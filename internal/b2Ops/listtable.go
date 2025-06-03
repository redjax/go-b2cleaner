package b2Ops

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/redjax/go-b2cleaner/internal/util"
)

// RenderFileEntriesTable prints a pretty table of file entries.
func RenderFileEntriesTable(entries []FileEntry, maxNameLen int, termWidth int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Format.Header = text.FormatTitle

	t.AppendHeader(table.Row{"Type", "Name", "Size", "Created"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignCenter}, // Type
		{
			Number: 2, Align: text.AlignLeft,
			WidthMax:         maxNameLen,
			WidthMaxEnforcer: text.WrapSoft,
			Transformer: func(val interface{}) string {
				s, _ := val.(string)
				if len(s) > maxNameLen*2 {
					return s[:maxNameLen*2-3] + "..."
				}
				return s
			},
		},
		{Number: 3, Align: text.AlignRight}, // Size
		{Number: 4, Align: text.AlignLeft},  // Created
	})

	fileCount := 0
	totalSize := int64(0)
	for _, e := range entries {
		var size string
		if !e.IsDir {
			size = util.HumanDiskSize(e.Size)
			fileCount++
			totalSize += e.Size
		}
		t.AppendRow(table.Row{
			map[bool]string{true: "DIR", false: "FILE"}[e.IsDir],
			e.Name,
			size,
			e.UploadTimestamp.Format("2006-01-02 15:04"),
		})
	}

	t.AppendFooter(table.Row{"", "Total File Count", fileCount, ""})
	t.AppendFooter(table.Row{"", "Total Size", util.HumanDiskSize(totalSize), ""})

	// Optionally print detected width to stderr for debugging
	fmt.Fprintf(os.Stderr, "Detected terminal width: %d\n", termWidth)

	t.Render()
}
