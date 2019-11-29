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
	"deepsea/global"
	"github.com/spf13/cobra"
	"fmt"
	"html/template"
	"io/ioutil"
	"github.com/russross/blackfriday/v2"
	jlog "github.com/spf13/jwalterweatherman"
)


var sourceFile string
var targetFile string

var md2htmlCmd = &cobra.Command{
	Use:   "md2html",
	Short: "Generate HTML content from Markdown",
	Long:  `MD2HTML: Generate HTML content from Markdown`,
	Run: func(cmd *cobra.Command, args []string) {
		jlog.DEBUG.Print("Cobra.Command: md2htmlDriver()")
		md2htmlDriver(cmd, args)
	},
}



func init() {
	md2htmlCmd.Flags().StringVarP(
		&sourceFile,
		"srcMd",
		"M",
		"",
		"Path to Source of MarkDown content (required)")

	md2htmlCmd.Flags().StringVarP(
		&targetFile,
		"tgtHtml",
		"H",
		"",
		"Path to Destination HTML file")

	contentCmd.AddCommand(md2htmlCmd)
}

func md2htmlDriver(cmd *cobra.Command, args []string) {

	jlog.INFO.Println("Processing Markdown")

	if len(sourceFile) == 0 {
		jlog.ERROR.Fatalln("Check presence of input/output file options")
	}

	if !global.FileExists(sourceFile) {
		jlog.ERROR.Fatalln("Check presence of input/output files.")
	}

	markdownContent := global.GetContentFromFileStr(sourceFile)
	htmlContent := convertMdToHTML(markdownContent)

	if len(targetFile) != 0 {
		err = ioutil.WriteFile(targetFile, []byte(htmlContent), 0644)
		if err != nil {
			jlog.FATAL.Printf("Cannot save HTML file: %s\n", err)
		}
	}else{
		fmt.Println(htmlContent)
	}
}

func convertMdToHTML(markdown string) template.HTML {
	return template.HTML(blackfriday.Run([]byte(string(markdown))))
}
