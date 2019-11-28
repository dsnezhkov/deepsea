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
	jlog "github.com/spf13/jwalterweatherman"
)

// contentCmd represents the content command
var contentCmd = &cobra.Command{
	Use:   "content",
	Short: `content generation and formatting of email templates`,
	Long: `content generation and formatting of email templates`,
	Run: func(cmd *cobra.Command, args []string) {
		contentDriver(cmd, args)
		jlog.ERROR.Println("`content` is a meta command. Try `content -h`")
	},
}

func init() {
	rootCmd.AddCommand(contentCmd)
}
func contentDriver(cmd *cobra.Command, args []string) {}
