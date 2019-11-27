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
	"log"
	"os"

	"github.com/spf13/cobra"
	jlog "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var cfgFile string
var Info bool
var Debug bool
var Trace bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "DeepSea",
	Short: "Red Team phishing gear",
	Long:  ` ROOT: see //dsnezhkov.github.io/deepsea...`,
	Run: func(cmd *cobra.Command, args []string) {
		jlog.INFO.Println("If you need help with usage => `deepsea help`")
	},
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		jlog.FATAL.Fatalf("Execute(): Cobra.Command: %v\n", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "", "config file (required)")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "v", false, "DEBUG output")
	rootCmd.PersistentFlags().BoolVarP(&Trace, "trace", "t", false, "TRACE output")
	rootCmd.PersistentFlags().BoolVarP(&Info, "info", "i", false, "INFO output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	jlog.SetStdoutThreshold(jlog.LevelWarn)

	if Info {
		jlog.SetStdoutThreshold(jlog.LevelInfo)
	}
	if Debug {
		jlog.SetStdoutThreshold(jlog.LevelDebug)
	}
	if Trace {
		jlog.SetStdoutThreshold(jlog.LevelTrace)
		jlog.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml")
	} else {
		jlog.ERROR.Println("Config file not provided")
		if err := rootCmd.Help(); err != nil {
			jlog.ERROR.Println("Error executing help()")
			os.Exit(2)
		}
		os.Exit(2)
	}
	//viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		jlog.DEBUG.Println("Using config file: ", viper.ConfigFileUsed())
	} else {
		jlog.ERROR.Fatalf("Config File Use Error: %v\n", err)
	}
}
