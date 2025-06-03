package b2Ops

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/kurin/blazer/b2"
	"github.com/redjax/go-b2cleaner/internal/config"
	"github.com/redjax/go-b2cleaner/internal/util"
)

type Client struct {
	b2Client *b2.Client
	cfg      config.Config
}

type FileEntry struct {
	Name            string
	Size            int64
	UploadTimestamp time.Time
	IsDir           bool
}

func NewClient(cfg config.Config) *Client {
	ctx := context.Background()

	b2Client, err := b2.NewClient(ctx, cfg.KeyId, cfg.AppKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to create B2 client: %v", err))
	}

	return &Client{
		b2Client: b2Client,
		cfg:      cfg,
	}
}

// ListObjects lists objects in a B2 bucket at a given path, optionally recursively,
// and sorts them according to cfg.SortBy.
func (c *Client) ListObjects(bucketName, prefix string, recurse bool) error {
	ctx := context.Background()
	bucket, err := c.b2Client.Bucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket: %w", err)
	}

	it := bucket.List(ctx, b2.ListPrefix(prefix))
	seenDirs := make(map[string]struct{})
	var entries []FileEntry

	for it.Next() {
		obj := it.Object()
		name := obj.Name()

		attrs, err := obj.Attrs(ctx)
		if err != nil {
			fmt.Printf("error getting attributes for %s: %v\n", name, err)
			continue
		}

		// Remove prefix from name for easier processing
		relName := strings.TrimPrefix(name, prefix)
		if relName == "" {
			continue
		}

		isDir := strings.HasSuffix(name, "/")

		if !recurse {
			// Non-recursive: only immediate children
			parts := strings.SplitN(relName, "/", 2)
			if len(parts) == 2 {
				dir := parts[0]
				if _, seen := seenDirs[dir]; !seen {
					entries = append(entries, FileEntry{
						Name:  prefix + dir + "/",
						IsDir: true,
					})
					seenDirs[dir] = struct{}{}
				}
				continue
			}
		}

		entries = append(entries, FileEntry{
			Name:            name,
			Size:            attrs.Size,
			UploadTimestamp: attrs.UploadTimestamp,
			IsDir:           isDir,
		})
	}

	if err := it.Err(); err != nil {
		return fmt.Errorf("error listing objects: %w", err)
	}

	// Sort entries based on cfg.SortBy (default "name")
	sortBy := strings.ToLower(c.cfg.SortBy)
	switch sortBy {
	case "size":
		sort.Slice(entries, func(i, j int) bool {
			// Directories first, then by size ascending
			if entries[i].IsDir != entries[j].IsDir {
				return entries[i].IsDir
			}
			return entries[i].Size < entries[j].Size
		})
	case "created":
		sort.Slice(entries, func(i, j int) bool {
			// Directories first, then by upload timestamp ascending
			if entries[i].IsDir != entries[j].IsDir {
				return entries[i].IsDir
			}
			return entries[i].UploadTimestamp.Before(entries[j].UploadTimestamp)
		})
	default: // "name" or any other value defaults to name sorting
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name < entries[j].Name
		})
	}

	// Reverse if order is "desc"
	if strings.ToLower(c.cfg.Order) == "desc" {
		for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
			entries[i], entries[j] = entries[j], entries[i]
		}
	}

	// Print entries
	for _, e := range entries {
		created := e.UploadTimestamp.Format("2006-01-02 15:04:05")
		if e.IsDir {
			fmt.Printf("[DIR]  %-40s %s\n", e.Name, created)
		} else {
			fmt.Printf("[FILE] |%s| (%10s) %-40s\n", created, util.HumanDiskSize(e.Size), e.Name)
		}
	}

	return nil
}
