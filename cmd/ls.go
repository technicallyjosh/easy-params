package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
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
	path := fmt.Sprintf("/%s", stripSlash(args[0]))

	recursive, _ := cmd.Flags().GetBool("recursive")
	decrypt, _ := cmd.Flags().GetBool("decrypt")
	displayValues, _ := cmd.Flags().GetBool("values")
	plain, _ := cmd.Flags().GetBool("plain")
	toEnv, _ := cmd.Flags().GetBool("env")

	fmt.Println(text.FgBlue.Sprintf("Listing parameters for \"%s\"", path))

	options := &getParamsOptions{
		Client:    ssm.NewFromConfig(awsConfig),
		Path:      &path,
		Recursive: recursive,
		Decrypt:   decrypt,
	}

	params := getParams(options, []types.Parameter{}, nil)

	sortParams(params)

	tw := table.NewWriter()

	header := table.Row{"Name", "Type", "Last Modified"}

	if displayValues {
		header = insertColumn(header, 1, "Value")
	}

	tw.AppendHeader(header)

	for _, param := range params {
		name := *param.Name
		var rest string

		if path == "/" {
			rest = strings.TrimPrefix(name, path)
		} else {
			rest = strings.TrimPrefix(name, fmt.Sprintf("%s/", path))
		}

		row := table.Row{
			rest,
			param.Type,
			formatDate(param.LastModifiedDate),
		}

		val := *param.Value

		if displayValues {
			row = insertColumn(row, 1, val)
		}

		if plain {
			if displayValues {
				if toEnv {
					fmt.Printf("%s=%s\n", strings.ToUpper(strings.ReplaceAll(rest, "-", "_")), val)
				} else {
					fmt.Printf("%s: %s\n", rest, val)
				}
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
	Client    *ssm.Client
	Path      *string
	Recursive bool
	Decrypt   bool
}

func getParams(options *getParamsOptions, params []types.Parameter, nextToken *string) []types.Parameter {
	cfg := &ssm.GetParametersByPathInput{
		Path:           options.Path,
		Recursive:      options.Recursive,
		WithDecryption: options.Decrypt,
	}

	if nextToken != nil {
		cfg.NextToken = nextToken
	}

	out, err := options.Client.GetParametersByPath(context.TODO(), cfg)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(out.Parameters); i++ {
		val := *out.Parameters[i].Value

		// seems like the upgrade to v2 of the aws-sdk sends back a new line on values :shrug:
		if strings.HasSuffix(val, "\n") {
			trimmed := strings.TrimSuffix(val, "\n")
			out.Parameters[i].Value = &trimmed
		}
	}

	params = append(params, out.Parameters...)

	if out.NextToken != nil {
		return getParams(options, params, out.NextToken)
	}

	return params
}

func init() {
	lsCmd.Flags().BoolP("recursive", "r", true, "recursively get values based on path")
	lsCmd.Flags().BoolP("decrypt", "d", true, "decrypt \"SecureString\" values")
	lsCmd.Flags().BoolP("values", "v", false, "display values")
	lsCmd.Flags().BoolP("plain", "p", false, "plain text instead of table")
	lsCmd.Flags().BoolP("env", "e", false, "output plain .env format")

	rootCmd.AddCommand(lsCmd)
}
