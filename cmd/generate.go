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
	"github.com/spf13/viper"
	"log"
	"os"

	"github.com/matcornic/hermes/v2"
	"io/ioutil"
)

var SourceMDFile string
var SourceTemplateHTMLFile string
var TargetHTMLFile string

// Theme
type DTheme struct{}

// Name returns the name of the default theme
func (dt *DTheme) Name() string {
	return "dtheme"
}

// HTMLTemplate returns a Golang template that will generate an HTML email.
func (dt *DTheme) HTMLTemplate() string {
	log.Printf("HTML template file: %s\n", viper.GetString("content.generate.SourceTemplateHTMLFile"))
	return string(global.GetContentFromFile(viper.GetString("content.generate.SourceTemplateHTMLFile")))
}

func (dt *DTheme) PlainTextTemplate() string {
	// return string(getContentFromFile( viper.GetString("content.generate.SourceTemplateTXTFile")))
	return ""
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate HTML content from HTML template",
	Long:  `GENERATE: Help here`,
	Run: func(cmd *cobra.Command, args []string) {
		generateDriver(cmd, args)
	},
}

func init() {

	generateCmd.Flags().StringVarP(&SourceMDFile, "SourceMDFile", "M",
		"", "Path to Source of MarkDown content")
	generateCmd.Flags().StringVarP(&SourceTemplateHTMLFile, "SourceTemplateHTMLFile",
		"T", "", "Path to Source of HTML template for Markdown")
	generateCmd.Flags().StringVarP(&TargetHTMLFile, "TargetHTMLFile", "H", "", "Path to Destination of HTML file")

	if err = viper.BindPFlag(
		"content.generate.SourceMDFile", generateCmd.Flags().Lookup("SourceMDFile")); err != nil {
		_ = generateCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"content.generate.SourceTemplateHTMLFile", generateCmd.Flags().Lookup("SourceTemplateHTMLFile")); err != nil {
		_ = generateCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"content.generate.TargetHTMLFile", generateCmd.Flags().Lookup("TargetHTMLFile")); err != nil {
		_ = generateCmd.Help()
		os.Exit(2)
	}

	contentCmd.AddCommand(generateCmd)
}

func generateDriver(cmd *cobra.Command, args []string) {

	log.Println("Processing Markdown")
	var messageFPath = viper.GetString("content.generate.SourceMDFile")
	log.Println(messageFPath)
	var mdMessage = hermes.Markdown(global.GetContentFromFileStr(messageFPath))

	// Defaults
	h := hermes.Hermes{

		Theme: new(DTheme),
		Product: hermes.Product{
			// TODO: reconcile to option w/ mailclient
			// Logo: "http://url",
			// Logo: getLogoFromFile(logoPath)
			Name:        "",
			Link:        "",
			Copyright:   "",
			TroubleText: "",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{

			Greeting:     "",
			Name:         "",
			Intros:       []string{},
			FreeMarkdown: mdMessage,
			Outros:       []string{},
			Signature:    "",
			Dictionary:   []hermes.Entry{},
			Actions:      []hermes.Action{},
			Title:        "",
		},
	}

	// Generate an HTML email with the provided contents
	emailHtml, err := h.GenerateHTML(email)
	if err != nil {
		log.Fatalf("Cannot Generate HTML from email: %s\n", err)
	}
	// Generated HTML to a local file
	err = ioutil.WriteFile(viper.GetString(
		"content.generate.TargetHTMLFile"), []byte(emailHtml), 0644)
	if err != nil {
		log.Fatalf("Cannot save HTML file: %s\n", err)
	}
	/*
		// Generate an TXT email with the provided contents
		emailTxt, err := h.GeneratePlainText(email)
		if err != nil {
			log.Fatalf("Cannot Generate TXTfrom email: %s\n", err)
		}
		// Generated TXT to a local file
		err = ioutil.WriteFile(viper.GetString("content.generate.TargetTXTFile"), []byte(emailTxt), 0644)
		if err != nil {
			log.Fatalf("Cannot save HTML file: %s\n", err)
		}*/
}
