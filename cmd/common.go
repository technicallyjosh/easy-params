package cmd

import (
	"github.com/jedib0t/go-pretty/table"
	"time"
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
