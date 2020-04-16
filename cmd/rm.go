package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	ssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <path>",
	Short: "Remove parameter(s) by path",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("requires a path")
		}

		return nil
	},
	Run: runRmCmd,
}

func runRmCmd(cmd *cobra.Command, args []string) {
	path := args[0]
	recursive, _ := cmd.Flags().GetBool("recursive")

	client := ssm.New(session)

	var names []*string

	if recursive {
		opts := &getParamsOptions{
			Client:    client,
			Path:      &path,
			Recursive: aws.Bool(true),
			Decrypt:   aws.Bool(false),
		}

		params := getParams(opts, []*ssm.Parameter{}, nil)

		if len(params) == 0 {
			fmt.Println("No parameters to delete at the specified path.")
			os.Exit(0)
		}

		sortParams(params)

		fmt.Println(text.FgYellow.Sprint("The following parameters will be removed..."))

		for _, param := range params {
			names = append(names, param.Name)

			fmt.Println(*param.Name)
		}

		fmt.Println(text.FgYellow.Sprint("Enter \"yes\" to delete the parameters above. Any other input will cancel the removal."))

		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		if input.Text() != "yes" {
			fmt.Println("Canceled the recursive delete!")
			os.Exit(0)
		}
	} else {
		names = append(names, &path)
	}

	chunks := getStringChunks(names, 10)

	for _, chunk := range chunks {
		_, err := client.DeleteParameters(&ssm.DeleteParametersInput{
			Names: chunk,
		})

		if err != nil {
			panic(err)
		}
	}

	fmt.Println(text.FgGreen.Sprintf("Deleted %d parameters successfully", len(names)))
}

func init() {
	rmCmd.Flags().Bool("recursive", false, "remove all children on path recursively")

	rootCmd.AddCommand(rmCmd)
}
