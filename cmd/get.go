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

	val, err := GetParameter(&GetParameterOptions{
		Client:  client,
		Name:    path,
		Decrypt: decrypt,
	})
	if err != nil {
		panic(err)
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

// GetParameterOptions represents options for GetParameter
type GetParameterOptions struct {
	Client  *ssm.Client
	Name    string
	Decrypt bool
}

// GetParamter returns a parameter by name while removing any trailing new-line
func GetParameter(options *GetParameterOptions) (string, error) {
	out, err := options.Client.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           &options.Name,
		WithDecryption: options.Decrypt,
	})
	if err != nil {
		return "", err
	}

	// return val, nil
	return strings.TrimSuffix(*out.Parameter.Value, "\n"), nil
}

func init() {
	getCmd.Flags().BoolP("decrypt", "d", true, "decrypt \"SecureString\" value")
	getCmd.Flags().BoolP("copy", "c", false, "copy to clipboard")

	rootCmd.AddCommand(getCmd)
}
