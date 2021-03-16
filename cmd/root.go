package cmd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	version      = ""
	cfgFile      string
	showVersion  bool
	awsConfig    aws.Config
	useLocalTime bool
	region       string
	loadConfig   bool
)

var rootCmd = &cobra.Command{
	Use:   "ezparams",
	Short: "An easy AWS Parameter Store CLI",
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			cmd.Println("version:", version)
			os.Exit(0)
		}

		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				panic(err)
			}

			os.Exit(0)
		}
	},
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ezparams.yaml)")
	rootCmd.PersistentFlags().BoolVar(&showVersion, "version", false, "show version")
	rootCmd.PersistentFlags().BoolVarP(&useLocalTime, "local-time", "l", true, "convert UTC to local time")
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "AWS region to use")
	rootCmd.PersistentFlags().BoolVar(&loadConfig, "load-config", true, "load aws config from ~/.aws/config")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}

		// Search config in home directory with name ".ezparams" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ezparams")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if loadConfig {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic(err)
	}

	awsConfig = cfg
}
