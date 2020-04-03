package cmd

import (
	"errors"
	"fmt"

	ssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <path>",
	Short: "Remove parameter or parameters by path",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("requires a path")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		client := ssm.New(session)

		_, err := client.DeleteParameters(&ssm.DeleteParametersInput{
			Names: []*string{&path},
		})

		if err != nil {
			panic(err)
		}

		fmt.Println(text.FgGreen.Sprintf("Deleted parameter \"%s\"", path))
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
