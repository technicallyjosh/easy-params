package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	ssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
	Use:   "put <path> <value>",
	Short: "Put parameter by path",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("requires a path")
		}

		if len(args) == 1 {
			return errors.New("requires a value")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		value := args[1]
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		valueType, _ := cmd.Flags().GetString("type")

		client := ssm.New(session)

		fmt.Println(text.FgBlue.Sprintf("Putting Parameter \"%s\"", path))

		_, err := client.PutParameter(&ssm.PutParameterInput{
			Name:      &path,
			Value:     &value,
			Type:      &valueType,
			Overwrite: &overwrite,
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
	},
}

func init() {
	putCmd.Flags().BoolP("overwrite", "o", false, "overwrite param if exists.")
	putCmd.Flags().StringP("type", "t", "SecureString", "Type of parameter.")

	rootCmd.AddCommand(putCmd)
}
