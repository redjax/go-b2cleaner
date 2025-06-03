package b2Ops

import (
	"context"
	"fmt"
	"strings"

	"github.com/kurin/blazer/b2"
	"github.com/redjax/go-b2cleaner/internal/config"
	"github.com/redjax/go-b2cleaner/internal/util"
)

type Client struct {
	b2Client *b2.Client
	cfg      config.Config
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

// ListObjects lists objects in a B2 bucket at a given path, optionally recursively
func (c *Client) ListObjects(bucketName, prefix string, recurse bool) error {
	ctx := context.Background()
	bucket, err := c.b2Client.Bucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket: %w", err)
	}

	it := bucket.List(ctx, b2.ListPrefix(prefix))
	seenDirs := make(map[string]struct{})

	for it.Next() {
		obj := it.Object()
		name := obj.Name()

		attrs, err := obj.Attrs(ctx)
		if err != nil {
			// Handle error, e.g., skip this object or report
			fmt.Printf("error getting attributes for %s: %v\n", obj.Name(), err)
			continue
		}
		size := attrs.Size

		// Remove prefix from name for easier processing
		relName := strings.TrimPrefix(name, prefix)
		if relName == "" {
			continue
		}

		// If not recursive, only show immediate children
		if !recurse {
			parts := strings.SplitN(relName, "/", 2)
			if len(parts) == 2 {
				dir := parts[0]
				if _, seen := seenDirs[dir]; !seen {
					fmt.Printf("[DIR] %s%s\n", prefix, dir)
					seenDirs[dir] = struct{}{}
				}

				continue
			}

			fmt.Printf("[FILE] (%10s) %s\n", util.HumanDiskSize(size), name)
		} else {
			if strings.HasSuffix(name, "/") {
				fmt.Printf("[DIR] %s\n", name)
			} else {
				fmt.Printf("[FILE] (%10s) %s\n", util.HumanDiskSize(size), name)
			}
		}
	}

	if err := it.Err(); err != nil {
		return fmt.Errorf("error listing objects: %w", err)
	}

	return nil
}
