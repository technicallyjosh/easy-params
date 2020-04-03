package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	homedir "github.com/mitchellh/go-homedir"
)

var cfgFile string
var session *awsSession.Session
var useLocalTime bool
var region string

var rootCmd = &cobra.Command{
	Use:   "ezp",
	Short: "An easy AWS Parameter Store CLI",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func er(msg interface{}) {
	fmt.Println(text.FgRed.Sprintf("Error: %s", msg))
	os.Exit(1)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ez-params.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&useLocalTime, "useLocalTime", "l", true, "convert UTC to local time")
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "AWS Region to use.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".ez-params" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ez-params")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	session = awsSession.Must(awsSession.NewSession(&aws.Config{
		Region: &region,
	}))
}
