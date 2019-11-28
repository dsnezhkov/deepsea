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
	"github.com/aymerick/douceur/inliner"
	"github.com/spf13/cobra"
	jlog "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
)

var SourceContentHTMLFile string
var TargetMailTemplateHTMLFile string

var inlineCmd = &cobra.Command{
	Use:   "inline",
	Short: "Inline HTML content for emailing",
	Long:  `INLINE: Inline HTML content for emailing`,
	Run: func(cmd *cobra.Command, args []string) {
		jlog.TRACE.Println("inlineDriver()")
		inlineDriver(cmd, args)
	},
}

func init() {

	inlineCmd.Flags().StringVarP(&SourceContentHTMLFile, "SourceContentHTMLFile", "S",
		"", "Path to Source of HTML content")
	inlineCmd.Flags().StringVarP(&TargetMailTemplateHTMLFile, "TargetMailTemplateHTMLFile", "E",
		"", "Path to Source of HTML template for Markdown")

	if err = viper.BindPFlag(
		"content.inline.SourceContentHTMLFile", inlineCmd.Flags().Lookup("SourceContentHTMLFile")); err != nil {
		_ = inlineCmd.Help()
		jlog.ERROR.Println("Error processing flag: `content.inline.SourceContentHTMLFile`")
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"content.inline.TargetMailTemplateHTMLFile", inlineCmd.Flags().Lookup("TargetMailTemplateHTMLFile")); err != nil {
		_ = inlineCmd.Help()
		jlog.ERROR.Println("Error processing flag: `content.inline.TargetContentHTMLFile`")
		os.Exit(2)
	}

	contentCmd.AddCommand(inlineCmd)
}

func inlineDriver(cmd *cobra.Command, args []string) {

	var source = viper.GetString("content.inline.SourceContentHTMLFile")
	var dest = viper.GetString("content.inline.TargetMailTemplateHTMLFile")
	log.Printf("content.inline.SourceContentHTMLFile: %s\n", source)
	log.Printf("content.inline.TargetMailTemplateHTMLFile: %s\n", dest)
	var htSrc = global.GetContentFromFileStr(source)

	html, err := inliner.Inline(htSrc)
	if err != nil {
		log.Fatalf("Unable to inline: %s : %s\n", html, err)
	}

	err = ioutil.WriteFile(dest, []byte(html), 0644)
	if err != nil {
		log.Fatalf("Unable to save: %s : %s\n", dest, err)
	}
}
