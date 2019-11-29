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
	"fmt"
	"io/ioutil"

	"github.com/jaytaylor/html2text"
	"github.com/spf13/cobra"
	jlog "github.com/spf13/jwalterweatherman"

	"deepsea/global"
)


var sourceHFile string
var targetTFile string

var html2txtCmd = &cobra.Command{
	Use:   "html2txt",
	Short: "Generate Text content from HTML",
	Long:  `HTML2TXT: Generate Text content from HTML`,
	Run: func(cmd *cobra.Command, args []string) {
		jlog.DEBUG.Print("Cobra.Command: html2txtDriver()")
		html2txtDriver(cmd, args)
	},
}

func init() {
	html2txtCmd.Flags().StringVarP(
		&sourceHFile,
		"srcHtml",
		"K",
		"",
		"Path to Source of HTML content (required)")

	html2txtCmd.Flags().StringVarP(
		&targetTFile,
		"tgtTxt",
		"L",
		"",
		"Path to Destination TXT file")

	contentCmd.AddCommand(html2txtCmd)
}

func html2txtDriver(cmd *cobra.Command, args []string) {

	jlog.INFO.Println("Processing HTML")

	if len(sourceHFile) == 0 {
		jlog.ERROR.Fatalln("Check presence of input/output file options")
	}

	if !global.FileExists(sourceHFile) {
		jlog.ERROR.Fatalln("Check presence of input/output files.")
	}

	htmlContent := global.GetContentFromFileStr(sourceHFile)
	textContent, err := convertHTML2TXT(htmlContent)
	if err != nil {
		jlog.ERROR.Fatalf("Unable to convert HTML: %v\n", err )
	}

	if len(targetTFile) != 0 {
		err = ioutil.WriteFile(targetTFile, []byte(textContent), 0644)
		if err != nil {
			jlog.FATAL.Printf("Cannot save TXT file: %s\n", err)
		}
	}else{
		fmt.Println(textContent)
	}
}

func convertHTML2TXT(html string) (string, error) {
	return html2text.FromString(html)
}
