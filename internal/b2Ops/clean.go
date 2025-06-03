package b2Ops

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kurin/blazer/b2"
	"github.com/redjax/go-b2cleaner/internal/util"
)

// CleanObjects finds and deletes (or dry-runs) files older than the given age string.
func (c *Client) CleanObjects(bucketName, prefix, ageStr string, dryRun bool, recurse bool) error {
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

	// Actually delete files
	for _, entry := range toDelete {
		obj := bucket.Object(entry.Name)
		if err := obj.Delete(ctx); err != nil {
			fmt.Printf("Failed to delete %s: %v\n", entry.Name, err)
		} else {
			fmt.Printf("Deleted %s\n", entry.Name)
		}
	}
	return nil
}
