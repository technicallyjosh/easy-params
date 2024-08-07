package cmd

import (
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

var migrateCmd = &cobra.Command{
	Use:   "migrate <source path> [destination path]",
	Short: "Migrate parameters by path",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("requires a path")
		}

		return nil
	},
	Run: runMigrateCmd,
}

func runMigrateCmd(cmd *cobra.Command, args []string) {
	pathFrom := args[0]
	pathTo := pathFrom

	if len(args) > 1 && args[1] != "" {
		pathTo = args[1]
	}

	regionFrom, _ := cmd.Flags().GetString("region-from")
	regionTo, _ := cmd.Flags().GetString("region-to")
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	if regionTo == "" {
		if pathTo == pathFrom {
			cmd.PrintErr("destination path cannot match source path if region-from matches region-to")
			os.Exit(1)
		}

		regionTo = regionFrom
	}

	clientFrom := ssm.NewFromConfig(awsConfig, func(o *ssm.Options) {
		o.Region = regionFrom
	})

	clientTo := ssm.NewFromConfig(awsConfig, func(o *ssm.Options) {
		o.Region = regionTo
	})

	fmt.Println(text.FgBlue.Sprintf("Migrating %s \"%s\" ==> %s \"%s\"", regionFrom, pathFrom, regionTo, pathTo))

	options := &GetParametersOptions{
		Client:    clientFrom,
		Path:      &pathFrom,
		Recursive: true,
		Decrypt:   true,
	}

	params := GetParameters(options, []types.Parameter{}, nil)

	fmt.Println(text.FgBlue.Sprintf("Found %d parameters to migrate...", len(params)))

	for _, param := range params {
		name := *param.Name

		// if pathTo is defined, remove source path name and prepend pathTo
		if pathTo != pathFrom {
			name = fmt.Sprintf("%s/%s", pathTo, strings.Replace(name, name[0:len(pathFrom)+1], "", -1))
		}

		fmt.Println(text.FgBlue.Sprintf("Migrating %s \"%s\" ==> %s \"%s\"", regionFrom, *param.Name, regionTo, name))

		input := &ssm.PutParameterInput{
			Name:      &name,
			Type:      param.Type,
			Value:     param.Value,
			Overwrite: &overwrite,
		}

		if _, err := clientTo.PutParameter(context.TODO(), input); err != nil {
			if strings.HasPrefix(err.Error(), "ParameterAlreadyExists") {
				fmt.Println(text.FgYellow.Sprintf("%s already exists... To overwrite, add the --overwrite flag.", name))
				continue
			}

			cmd.PrintErr(err)
			os.Exit(1)
		}

		fmt.Println(text.FgGreen.Sprintf("Created parameter \"%s\" successfully.", name))
	}
}

func init() {
	migrateCmd.Flags().StringP("region-from", "f", "", "the region to migrate from")
	migrateCmd.Flags().StringP("region-to", "t", "", "the region to migrate to")
	migrateCmd.Flags().Bool("overwrite", false, "overwrite destination params")

	if err := migrateCmd.MarkFlagRequired("region-from"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(migrateCmd)
}
