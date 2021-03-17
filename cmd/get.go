package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

	fmt.Println(val)
}

func init() {
	getCmd.Flags().BoolP("decrypt", "d", true, "decrypt \"SecureString\" value")

	rootCmd.AddCommand(getCmd)
}
