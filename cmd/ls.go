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
	Run: runLsCmd,
}

func runLsCmd(cmd *cobra.Command, args []string) {
	path := args[0]

	recursive, _ := cmd.Flags().GetBool("recursive")
	decrypt, _ := cmd.Flags().GetBool("decrypt")
	displayValues, _ := cmd.Flags().GetBool("values")
	plain, _ := cmd.Flags().GetBool("plain")

	fmt.Println(text.FgBlue.Sprintf("Listing parameters for \"%s\"", path))

	options := &getParamsOptions{
		Client:    ssm.New(session),
		Path:      &path,
		Recursive: &recursive,
		Decrypt:   &decrypt,
	}

	params := getParams(options, []*ssm.Parameter{}, nil)

	sort.Slice(params, func(i, j int) bool {
		return strings.ToLower(*params[i].Name) < strings.ToLower(*params[j].Name)
	})

	tw := table.NewWriter()

	header := table.Row{"Name", "Type", "Last Modified"}

	if displayValues {
		header = insertColumn(header, 1, "Value")
	}

	tw.AppendHeader(header)

	for _, param := range params {
		name := *param.Name
		rest := strings.Replace(name, name[0:len(path)+1], "", -1)

		row := table.Row{
			rest,
			*param.Type,
			formatDate(param.LastModifiedDate),
		}

		if displayValues {
			row = insertColumn(row, 1, *param.Value)
		}

		if plain {
			if displayValues {
				fmt.Println(fmt.Sprintf("%s: %s", rest, *param.Value))
			} else {
				fmt.Println(rest)
			}
		}

		tw.AppendRow(row)
	}

	if !plain {
		fmt.Println(tw.Render())
	}
}

type getParamsOptions struct {
	Client    *ssm.SSM
	Path      *string
	Recursive *bool
	Decrypt   *bool
}

func getParams(options *getParamsOptions, params []*ssm.Parameter, nextToken *string) []*ssm.Parameter {
	cfg := &ssm.GetParametersByPathInput{
		Path:           options.Path,
		Recursive:      options.Recursive,
		WithDecryption: options.Decrypt,
	}

	if nextToken != nil {
		cfg.NextToken = nextToken
	}

	out, err := options.Client.GetParametersByPath(cfg)
	if err != nil {
		panic(err)
	}

	params = append(params, out.Parameters...)

	if out.NextToken != nil {
		return getParams(options, params, out.NextToken)
	}

	return params
}

func init() {
	lsCmd.Flags().BoolP("recursive", "r", true, "Recursively get values based on path")
	lsCmd.Flags().BoolP("decrypt", "d", true, "Decrypt SecureString values")
	lsCmd.Flags().BoolP("values", "v", false, "Display values")
	lsCmd.Flags().BoolP("plain", "p", false, "Plain text instead of tables")

	rootCmd.AddCommand(lsCmd)
}
