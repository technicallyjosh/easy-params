package cmd

import (
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/jedib0t/go-pretty/table"
)

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

func getStringChunks(arr []*string, batchSize int) [][]*string {
	batches := make([][]*string, 0, (len(arr)+batchSize-1)/batchSize)

	for batchSize < len(arr) {
		arr, batches = arr[batchSize:], append(batches, arr[0:batchSize:batchSize])
	}
	batches = append(batches, arr)

	return batches
}

func sortParams(params []*ssm.Parameter) {
	sort.Slice(params, func(i, j int) bool {
		return strings.ToLower(*params[i].Name) < strings.ToLower(*params[j].Name)
	})
}
