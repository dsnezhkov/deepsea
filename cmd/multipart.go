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
	"bytes"
	"deepsea/global"
	"github.com/jaytaylor/html2text"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"

	jlog "github.com/spf13/jwalterweatherman"
)

var SourceMailTemplateHTMLFile string
var TargetMailTemplateTXTFile string

var multipartCmd = &cobra.Command{
	Use:   "multipart",
	Short: "Multipart HTML->TXT content for emailing",
	Long:  `MULTIPART: Help here`,
	Run: func(cmd *cobra.Command, args []string) {
		jlog.DEBUG.Println("multipartDriver()")
		multipartDriver(cmd, args)
	},
}

func init() {

	multipartCmd.Flags().StringVarP(&SourceMailTemplateHTMLFile, "SourceMailTemplateHTMLFile", "S",
		"", "Path to Source of HTML template")
	multipartCmd.Flags().StringVarP(&TargetMailTemplateTXTFile, "TargetMailTemplateTXTFile", "E",
		"", "Path to Destination of TXT template")

	if err = viper.BindPFlag(
		"content.multipart.SourceMailTemplateHTMLFile", multipartCmd.Flags().Lookup("SourceMailTemplateHTMLFile")); err != nil {
		_ = multipartCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"content.multipart.TargetMailTemplateTXTFile", multipartCmd.Flags().Lookup("TargetMailTemplateTXTFile")); err != nil {
		_ = multipartCmd.Help()
		os.Exit(2)
	}

	contentCmd.AddCommand(multipartCmd)
}

func multipartDriver(cmd *cobra.Command, args []string) {

	var source = viper.GetString("content.multipart.SourceMailTemplateHTMLFile")
	var dest = viper.GetString("content.multipart.TargetMailTemplateTXTFile")

	jlog.DEBUG.Printf("content.multipart.SourceMailTemplateHTMLFile: %s\n", source)
	jlog.DEBUG.Printf("content.multipart.TargetMailTemplateTXTFile: %s\n", dest)

	var htSrc = global.GetContentFromFile(source)

	text, err := html2text.FromReader(bytes.NewReader(htSrc))
	if err != nil {
		jlog.ERROR.Fatalf("Unable to multipart: %s : %s\n", source, err)
	}

	err = ioutil.WriteFile(dest, []byte(text), 0644)
	if err != nil {
		jlog.ERROR.Fatalf("Unable to save text file: %s : %s\n", dest, err)
	}

}
