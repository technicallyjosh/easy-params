package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <path(s)>",
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

	client := ssm.NewFromConfig(awsConfig)

	var names []string

	if recursive {
		opts := &GetParametersOptions{
			Client:    client,
			Path:      &path,
			Recursive: true,
			Decrypt:   false,
		}

		params := GetParameters(opts, []types.Parameter{}, nil)

		if len(params) == 0 {
			fmt.Println("No parameters to delete at the specified path.")
			os.Exit(0)
		}

		fmt.Println(text.FgYellow.Sprint("The following parameters will be removed..."))

		for _, param := range params {
			names = append(names, *param.Name)

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
		for _, arg := range args {
			// declare in loop for pointer to not be repeated
			str := arg
			names = append(names, str)
		}
	}

	chunks := getStringChunks(names, 10)

	num := 0

	for _, chunk := range chunks {
		res, err := client.DeleteParameters(context.TODO(), &ssm.DeleteParametersInput{
			Names: chunk,
		})

		if err != nil {
			panic(err)
		}

		num += len(res.DeletedParameters)
	}

	fmt.Println(text.FgGreen.Sprintf("Deleted %d parameters successfully", num))
}

func init() {
	rmCmd.Flags().Bool("recursive", false, "remove all children on path recursively")

	rootCmd.AddCommand(rmCmd)
}
