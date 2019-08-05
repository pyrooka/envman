/*
Copyright © 2019 pyrooka

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "envman",
	Short: "An environment manager tool.",
	Long: `Envman is an environment manager tool for your CLI.
You can simply:
  - list your environments or variables
  - load a previously save environment to your current session
  - save currently set environment variables to a storage
  - remove existing environments or just variables
  - cleanup your workspace by removing any created data from the storages and the local computer`,
	PersistentPostRun: postRun,
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

	rootCmd.PersistentFlags().String("storage", "", "set the storage for use")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	_, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".envman" (without extension).
	// viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	viper.SetConfigName("envman")
	viper.SetConfigType("yaml")
	viper.SetDefault("storage", "local")
	viper.BindPFlag("storage", rootCmd.Flags().Lookup("storage"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok == true {
			// Nothing to do here, the config doesn't exist yet. Don't worry we will create it!
		} else {
			fmt.Println("Error occured during the config initialization:", err)
			os.Exit(1)
		}
	}
}

func postRun(com *cobra.Command, args []string) {
	err := viper.WriteConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok == true {
			err = nil
			err = viper.WriteConfigAs(".envman.yaml")
		}
		if err != nil {
			fmt.Println("Error occured while trying to write the config:", err)
		}
	}
}
