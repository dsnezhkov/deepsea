// Copyright © 2019 Dimitry Snezhkov <dsnezhkov@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "DeepSea",
	Short: "Red Team phishing gear",
	Long:  ` ROOT: see //dsnezhkov.github.io/deepsea...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("... Deep Sea ...")

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("[Error] cobra.Command: %v\n", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "", "config file (required)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml")
	} else {
		log.Printf("[Error]: Config file not provided")
		if err := rootCmd.Help(); err != nil {
			fmt.Print("Error executing help")
			os.Exit(2)
		}
		os.Exit(1)
	}
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("[Info] Using config file:", viper.ConfigFileUsed())
	} else {
		log.Printf("[Error] Config File Use Error: %v\n", err)
		os.Exit(2)
	}
}
