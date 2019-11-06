// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var DBFile string

// storageCmd represents the storage command
var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Manage persistent record storage",
	Long:  `STORAGE: TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		storageDriver(cmd, args)
	},
}

func init() {

	storageCmd.Flags().StringVarP(&DBFile, "DBFile", "d",
		"", "Path to QL DB file")

	if err = viper.BindPFlag(
		"storage.DBFile", storageCmd.Flags().Lookup("DBFile")); err != nil {
		_ = storageCmd.Help()
		os.Exit(2)
	}

	rootCmd.AddCommand(storageCmd)
}
func storageDriver(cmd *cobra.Command, args []string) {}
