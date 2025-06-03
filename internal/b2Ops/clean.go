package b2Ops

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kurin/blazer/b2"
	"github.com/redjax/go-b2cleaner/internal/util"
)

// CleanObjects finds and deletes (or dry-runs) files older than the given age string.
func (c *Client) CleanObjects(bucketName, prefix, ageStr string, dryRun bool, recurse bool, outputPath string, filetypes []string) error {
	// Parse the age string (e.g., "30d", "2m", "1y")
	age, err := util.ParseAgeString(ageStr)
	if err != nil {
		return fmt.Errorf("invalid age: %w", err)
	}
	cutoff := time.Now().Add(-age)

	ctx := context.Background()
	bucket, err := c.b2Client.Bucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket: %w", err)
	}

	it := bucket.List(ctx, b2.ListPrefix(prefix))

	var toDelete []FileEntry

	var normalizedTypes []string
	for _, ft := range filetypes {
		if ft == "" {
			continue
		}
		if !strings.HasPrefix(ft, ".") {
			normalizedTypes = append(normalizedTypes, "."+ft)
		} else {
			normalizedTypes = append(normalizedTypes, ft)
		}
	}

	if ageStr == "" && len(filetypes) == 0 {
		fmt.Printf("Cleaning bucket %s at path %s\n", bucketName, prefix)
	} else if ageStr != "" || len(filetypes) == 0 {
		fmt.Printf("Cleaning bucket %s at path %s older than %s\n", bucketName, prefix, ageStr)
	} else {
		fmt.Printf("Cleaning bucket %s at path %s with filetypes %v\n", bucketName, prefix, normalizedTypes)
	}

	// DELETE loop
	for it.Next() {
		obj := it.Object()
		attrs, err := obj.Attrs(ctx)
		if err != nil {
			continue
		}

		// If not recursing, skip files not immediate children of prefix
		if !recurse {
			relName := strings.TrimPrefix(obj.Name(), prefix)
			relName = strings.TrimPrefix(relName, "/")
			if relName == "" {
				continue
			}
			if strings.Contains(relName, "/") {
				continue // skip files in subdirectories
			}
		}

		// Only consider files (not directories) older than cutoff
		if attrs.UploadTimestamp.Before(cutoff) && !strings.HasSuffix(obj.Name(), "/") {
			// Filter by filetype
			if len(normalizedTypes) > 0 {
				matched := false
				for _, ext := range normalizedTypes {
					if strings.HasSuffix(obj.Name(), ext) {
						matched = true
						break
					}
				}
				if !matched {
					continue // skip files not matching any extension
				}
			}

			toDelete = append(toDelete, FileEntry{
				Name:            obj.Name(),
				Size:            attrs.Size,
				UploadTimestamp: attrs.UploadTimestamp,
				IsDir:           false,
			})
		}
	}

	if err := it.Err(); err != nil {
		return fmt.Errorf("error listing objects: %w", err)
	}

	if dryRun {
		fmt.Println("Dry run: The following files would be deleted:")
		termWidth := util.DetectTerminalWidth(160)
		maxNameLen := util.MaxNameLen(termWidth, 6, 10, 18, 9)
		RenderFileEntriesTable(toDelete, maxNameLen, termWidth)
		return nil
	}

	var deletedEntries []FileEntry

	// Actually delete files and print each deleted object
	for _, entry := range toDelete {
		obj := bucket.Object(entry.Name)
		if err := obj.Delete(ctx); err != nil {
			fmt.Printf("Failed to delete %s: %v\n", entry.Name, err)
		} else {
			fmt.Printf("Deleted %s\n", entry.Name)
			deletedEntries = append(deletedEntries, entry)
		}
	}

	// Write deleted entries to CSV if requested
	if outputPath != "" && len(deletedEntries) > 0 {
		outputPath = util.SanitizeFileName(outputPath, ".csv")

		file, err := os.Create(outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create output file: %v\n", err)
			return err
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write header
		if err := writer.Write([]string{"Name", "Size", "Created"}); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write CSV header: %v\n", err)
			return err
		}

		// Write each deleted entry
		for _, entry := range deletedEntries {
			record := []string{
				entry.Name,
				fmt.Sprintf("%d", entry.Size),
				entry.UploadTimestamp.Format("2006-01-02 15:04:05"),
			}
			if err := writer.Write(record); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write CSV record: %v\n", err)
				return err
			}
		}
		fmt.Printf("Deleted objects written to %s\n", outputPath)
	}

	return nil
}
