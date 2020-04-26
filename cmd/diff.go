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
	Key    string
	Value1 string
	Value2 string
	Path   int
}

func runDiffCmd(cmd *cobra.Command, args []string) {
	path1 := args[0]
	path2 := args[1]

	showValues, _ := cmd.Flags().GetBool("values")
	decrypt, _ := cmd.Flags().GetBool("decrypt")
	widthLimit, _ := cmd.Flags().GetInt("width-limit")

	fmt.Println(text.FgBlue.Sprintf("Getting diff between \"%s\" and \"%s\"...", path1, path2))

	options := &getParamsOptions{
		Client:    ssm.New(session),
		Path:      &path1,
		Recursive: aws.Bool(true),
		Decrypt:   &decrypt,
	}

	params1 := getParams(options, []*ssm.Parameter{}, nil)

	options.Path = &path2
	params2 := getParams(options, []*ssm.Parameter{}, nil)

	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatLower

	headerRow := table.Row{path1, path2}

	if showValues {
		headerRow = insertColumn(headerRow, 1, "Value")
		headerRow = insertColumn(headerRow, 3, "Value")

		if widthLimit > 0 {
			tw.SetColumnConfigs([]table.ColumnConfig{
				{
					Number:   2,
					WidthMax: widthLimit,
				},
				{
					Number:   4,
					WidthMax: widthLimit,
				},
			})
		}
	}

	tw.AppendHeader(headerRow)

	var rows []*diffRow

	for i := range params1 {
		name := *params1[i].Name
		val := *params1[i].Value
		key := strings.Replace(name, name[0:len(path1)+1], "", -1)

		rows = append(rows, &diffRow{
			Key:    key,
			Value1: val,
			Path:   1,
		})
	}

	for i := range params2 {
		name := *params2[i].Name
		val := *params2[i].Value
		key := strings.Replace(name, name[0:len(path2)+1], "", -1)

		updated := false

		// look for existing
		for j := range rows {
			row := rows[j]

			if row.Key == key {
				row.Value2 = val
				row.Path = 0
				updated = true
			}
		}

		if !updated {
			rows = append(rows, &diffRow{
				Key:    key,
				Value2: val,
				Path:   2,
			})
		}
	}

	sortDiffRows(rows)

	for i := range rows {
		var key1 string
		var key2 string

		key := rows[i].Key
		greenKey := text.FgGreen.Sprint(key)
		redKey := text.FgRed.Sprint(key)

		switch path := rows[i].Path; path {
		case 0:
			key1 = key
			key2 = key
		case 1:
			key1 = greenKey
			key2 = redKey
		case 2:
			key1 = redKey
			key2 = greenKey
		}

		row := table.Row{key1, key2}

		if showValues {
			value1 := rows[i].Value1
			value2 := rows[i].Value2

			// if both exist and there is a difference
			if value1 != "" && value2 != "" && value1 != value2 {
				value1 = text.FgYellow.Sprint(value1)
				value2 = text.FgYellow.Sprint(value2)
			}

			row = insertColumn(row, 1, value1)
			row = insertColumn(row, 3, value2)
		}

		tw.AppendRow(row)
	}

	fmt.Println(tw.Render())
}

func init() {
	diffCmd.Flags().BoolP("values", "v", false, "show value diffs")
	diffCmd.Flags().BoolP("decrypt", "d", true, "decrypt \"SecureString\" values")
	diffCmd.Flags().IntP("width-limit", "w", 0, "width limit of value output")

	rootCmd.AddCommand(diffCmd)
}
