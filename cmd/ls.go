package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	ssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls <path>",
	Short: "List parameters by path",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("requires a path")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		recursive, _ := cmd.Flags().GetBool("recursive")
		decrypt, _ := cmd.Flags().GetBool("decrypt")

		fmt.Println(text.FgBlue.Sprintf("Listing parameter(s) for \"%s\"", path))

		getParams(&path, &recursive, &decrypt, []*ssm.Parameter{}, nil)
	},
}

func getParams(path *string, recursive *bool, decrypt *bool, params []*ssm.Parameter, nextToken *string) {
	client := ssm.New(session)

	cfg := &ssm.GetParametersByPathInput{
		Path:           path,
		Recursive:      recursive,
		WithDecryption: decrypt,
	}

	if nextToken != nil {
		cfg.NextToken = nextToken
	}

	out, err := client.GetParametersByPath(cfg)
	if err != nil {
		panic(err)
	}

	params = append(params, out.Parameters...)

	if out.NextToken != nil {
		getParams(path, recursive, decrypt, params, out.NextToken)
		return
	}

	sort.Slice(params, func(i, j int) bool {
		return strings.ToLower(*params[i].Name) < strings.ToLower(*params[j].Name)
	})

	tw := table.NewWriter()

	tw.AppendHeader(table.Row{"Name", "Value", "Type", "Last Modified"})

	for _, param := range params {
		name := *param.Name
		rest := strings.Replace(name, name[0:len(*path)+1], "", -1)

		tw.AppendRow(table.Row{
			rest,
			*param.Value,
			*param.Type,
			formatDate(param.LastModifiedDate),
		})
	}

	fmt.Println(tw.Render())
}

func init() {
	lsCmd.Flags().BoolP("recursive", "r", true, "Recursively get values based on path.")
	lsCmd.Flags().BoolP("decrypt", "d", true, "Decrypt SecureString values.")

	rootCmd.AddCommand(lsCmd)
}
