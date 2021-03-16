package cmd

import (
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/jedib0t/go-pretty/table"
)

// formatDate formats a time.Time instance.
func formatDate(dt *time.Time) string {
	if useLocalTime {
		return dt.Local().Format("2006-01-02 03:04:05 PM")
	}

	return dt.Format("2006-01-02 15:04:05")
}

// insertColumn inserts a given item into a Row at a specified index
func insertColumn(row table.Row, index int, item interface{}) table.Row {
	row = append(row, item)

	copy(row[index+1:], row[index:])

	row[index] = item

	return row
}

// getStringChunks returns chunks for a slice of strings.
func getStringChunks(items []string, chunkSize int) (chunks [][]string) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}

// sortParams sorts parameters by name.
func sortParams(params []types.Parameter) {
	sort.Slice(params, func(i, j int) bool {
		return strings.ToLower(*params[i].Name) < strings.ToLower(*params[j].Name)
	})
}

// sortDiffRows sorts diff rows by key.
func sortDiffRows(rows []*diffRow) {
	sort.Slice(rows, func(i, j int) bool {
		return strings.ToLower(rows[i].Key) < strings.ToLower(rows[j].Key)
	})
}

// stripSlash strips the first slash.
func stripSlash(str string) string {
	return strings.TrimRight(strings.TrimLeft(str, "/"), "/")
}
