package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
	Use:   "put <path> <value>",
	Short: "Put parameter by path",
	Args: func(cmd *cobra.Command, args []string) error {
		context, _ := cmd.Flags().GetString("context")

		if context == "" && len(args) < 2 {
			return errors.New("requires a path and a value")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		context, _ := cmd.Flags().GetString("context")

		if context != "" {
			// when running the command, we remove any leading or trailing /
			runPutCmdContext(cmd, args, stripSlash(context))
			return
		}

		runPutCmd(cmd, args)
	},
}

func runPutCmdContext(cmd *cobra.Command, args []string, ctx string) {
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	valueType, _ := cmd.Flags().GetString("type")

	client := ssm.NewFromConfig(awsConfig)

	ctxMessage := fmt.Sprintf("context = /%s", ctx)
	kvMessage := fmt.Sprintf("enter key/value pairs to put in the format of \"key value\".")

	fmt.Printf("%v\n", text.FgYellow.Sprint(ctxMessage))
	fmt.Printf("%v\n", text.FgYellow.Sprintf("overwrite = %v", overwrite))
	fmt.Printf("%v\n", text.FgYellow.Sprintf("type = %s", valueType))
	fmt.Printf("%s\n\n", kvMessage)

	for {
		fmt.Print("> ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		txt := strings.TrimSpace(input.Text())
		if txt == "" {
			break
		}

		pair := strings.Split(txt, " ")
		if len(pair) < 2 {
			fmt.Println(text.FgRed.Sprint("param and value must be defined"))
			continue
		}

		param := stripSlash(pair[0])
		value := strings.TrimSpace(strings.Join(pair[1:], " "))

		if param == "" || value == "" {
			fmt.Println(text.FgRed.Sprint("param and value cannot be empty"))
			continue
		}

		path := fmt.Sprintf("/%s/%s", ctx, param)

		_, err := client.PutParameter(context.TODO(), &ssm.PutParameterInput{
			Name:      &path,
			Value:     &value,
			Type:      types.ParameterType(valueType),
			Overwrite: overwrite,
		})

		if err != nil {
			if strings.HasPrefix(err.Error(), "ParameterAlreadyExists") {
				fmt.Println(text.FgRed.Sprintf("Parameter \"%s\" already exists. Use the --overwrite option to update.", path))
			} else {
				fmt.Println(text.FgRed.Sprint(err.Error()))
			}
		}
	}
}

func runPutCmd(cmd *cobra.Command, args []string) {
	path := args[0]
	value := args[1]
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	valueType, _ := cmd.Flags().GetString("type")

	client := ssm.NewFromConfig(awsConfig)

	_, err := client.PutParameter(context.TODO(), &ssm.PutParameterInput{
		Name:      &path,
		Value:     &value,
		Type:      types.ParameterType(valueType),
		Overwrite: overwrite,
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "ParameterAlreadyExists") {
			fmt.Println(text.FgRed.Sprintf("Parameter \"%s\" already exists. Use the --overwrite option to update.", path))
			os.Exit(1)
		} else {
			panic(err)
		}
	}

	fmt.Println(text.FgGreen.Sprintf("Put parameter \"%s\" successfully", path))
}

func init() {
	putCmd.Flags().BoolP("overwrite", "o", false, "overwrite param if exists.")
	putCmd.Flags().StringP("type", "t", "SecureString", "type of parameter.")
	putCmd.Flags().StringP("context", "c", "", "context mode for setting many values.")

	rootCmd.AddCommand(putCmd)
}
