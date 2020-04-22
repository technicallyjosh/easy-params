package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff <path 1> <path 2>",
	Short: "Shows the difference recursively between 2 paths.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires 2 paths")
		}

		return nil
	},
	Run: runDiffCmd,
}

type diffRow struct {
	Key  string
	Path int
}

func runDiffCmd(cmd *cobra.Command, args []string) {
	path1 := args[0]
	path2 := args[1]

	fmt.Println(text.FgBlue.Sprintf("Getting diff between \"%s\" and \"%s\"...", path1, path2))

	options := &getParamsOptions{
		Client:    ssm.New(session),
		Path:      &path1,
		Recursive: aws.Bool(true),
	}

	params1 := getParams(options, []*ssm.Parameter{}, nil)

	options.Path = &path2
	params2 := getParams(options, []*ssm.Parameter{}, nil)

	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatLower

	tw.AppendHeader(table.Row{path1, path2})

	var rows []*diffRow

	for i := range params1 {
		name := *params1[i].Name
		key := strings.Replace(name, name[0:len(path1)+1], "", -1)

		rows = append(rows, &diffRow{
			Key:  key,
			Path: 1,
		})
	}

	for i := range params2 {
		name := *params2[i].Name
		key := strings.Replace(name, name[0:len(path2)+1], "", -1)

		updated := false

		// look for existing
		for j := range rows {
			row := rows[j]

			if row.Key == key {
				row.Path = 0
				updated = true
			}
		}

		if !updated {
			rows = append(rows, &diffRow{
				Key:  key,
				Path: 2,
			})
		}
	}

	sortDiffRows(rows)

	for i := range rows {
		var text1 string
		var text2 string

		key := rows[i].Key
		greenKey := text.FgGreen.Sprint(key)
		redKey := text.FgRed.Sprint(key)

		switch path := rows[i].Path; path {
		case 0:
			text1 = key
			text2 = key
		case 1:
			text1 = greenKey
			text2 = redKey
		case 2:
			text1 = redKey
			text2 = greenKey
		}

		tw.AppendRow(table.Row{text1, text2})
	}

	fmt.Println(tw.Render())
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
