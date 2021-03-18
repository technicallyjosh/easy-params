package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Get parameter value by path",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("requires a path")
		}

		return nil
	},
	Run: runGetCmd,
}

func runGetCmd(cmd *cobra.Command, args []string) {
	path := fmt.Sprintf("/%s", stripSlash(args[0]))
	decrypt, _ := cmd.Flags().GetBool("decrypt")
	copy, _ := cmd.Flags().GetBool("copy")

	client := ssm.NewFromConfig(awsConfig)

	out, err := client.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           &path,
		WithDecryption: decrypt,
	})
	if err != nil {
		panic(err)
	}

	val := *out.Parameter.Value

	if strings.HasSuffix(val, "\n") {
		val = strings.TrimSuffix(val, "\n")
	}

	if copy {
		if err := clipboard.WriteAll(val); err != nil {
			panic(err)
		}
		fmt.Println("Value copied to clipboard!")
	} else {
		fmt.Println(val)
	}
}

func init() {
	getCmd.Flags().BoolP("decrypt", "d", true, "decrypt \"SecureString\" value")
	getCmd.Flags().BoolP("copy", "c", false, "copy to clipboard")

	rootCmd.AddCommand(getCmd)
}
