// Copyright © 2019 D.Snezhkov <dsnezhkov@gmail.com>
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
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
	"upper.io/db.v3/ql"

	"deepsea/global"
	"deepsea/rmailer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var SMTPServer string
var SMTPUser string
var SMTPPort int
var SMTPPass string
var TLS string // YAML does not handle bools

var From string
var To string
var Subject string

var Headers map[string]string
var BodyHTMLTemplate string
var BodyTextTemplate string

var Attachments []string
var Embeds []string
var err error
var tdata rmailer.TemplateData

var mailclientCmd = &cobra.Command{
	Use:   "mailclient",
	Short: "Email a phish",
	Long:  `Email a phish with features`,
	Run: func(cmd *cobra.Command, args []string) {
		mailDriver(cmd, args)
	},
}

func init() {

	// Bind command flags to configuration file directives.
	// So we can override them ont he command line

	// Example:
	// Cobra Flag name: SMTPPort
	// Viper key:       mailclient.SMTPPort
	// Variable :       SMTPPort
	// Bind Viper's key to Cobra's Flag name

	// mailclient.connection
	mailclientCmd.Flags().StringVarP(&SMTPServer, "SMTPServer", "s",
		"127.0.0.1", "SMTP server")
	mailclientCmd.Flags().IntVarP(&SMTPPort, "SMTPPort", "p",
		25, "SMTP server port")
	mailclientCmd.Flags().StringVarP(&SMTPUser, "SMTPUser", "U",
		"testuser", "SMTP user")
	mailclientCmd.Flags().StringVarP(&TLS, "TLS", "t",
		"yes", "Use TLS handshake (STARTTLS)")

	if err = viper.BindPFlag(
		"mailcient.connection.SMTPServer", mailclientCmd.Flags().Lookup("SMTPServer")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"mailclient.connection.SMTPPort", mailclientCmd.Flags().Lookup("SMTPPort")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"mailclient.connection.SMTPUser", mailclientCmd.Flags().Lookup("SMTPUser")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"mailclient.connection.TLS", mailclientCmd.Flags().Lookup("TLS")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}

	// mailclient.message
	mailclientCmd.Flags().StringVarP(&From, "From", "F",
		"", "Message From: header")
	mailclientCmd.Flags().StringVarP(&To, "To", "T",
		"", "Message To: header")
	mailclientCmd.Flags().StringVarP(&Subject, "Subject", "S",
		"", "Message Subject: header")

	if err = viper.BindPFlag(
		"mailclient.message.From", mailclientCmd.Flags().Lookup("From")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"mailclient.message.To", mailclientCmd.Flags().Lookup("To")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"mailclient.message.Subject", mailclientCmd.Flags().Lookup("Subject")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}

	// mailclient.message.body
	mailclientCmd.Flags().StringVarP(&BodyHTMLTemplate, "HTMLTemplate", "H",
		"", "HTML Template file (.htpl)")
	mailclientCmd.Flags().StringVarP(&BodyTextTemplate, "TextTemplate", "P",
		"", "Text Template file (.ttpl)")

	if err = viper.BindPFlag(
		"mailclient.message.body.html", mailclientCmd.Flags().Lookup("HTMLTemplate")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"mailclient.message.body.text", mailclientCmd.Flags().Lookup("TextTemplate")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}

	rootCmd.AddCommand(mailclientCmd)
}

// Processing
func getUserCredentials(server *string) (string, error) {

	fmt.Printf("-= SMTP Authentication Credentials for %s =- \n", *server)
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	password := string(bytePassword)

	return strings.TrimSpace(password), nil
}

func mailDriver(cmd *cobra.Command, args []string) {

	SMTPServer = viper.GetString("mailclient.connection.SMTPServer")
	SMTPPort = viper.GetInt("mailclient.connection.SMTPPort")
	SMTPUser = viper.GetString("mailclient.connection.SMTPUser")
	SMTPUser = viper.GetString("mailclient.connection.SMTPUser")
	TLS = viper.GetString("mailclient.connection.TLS")

	// Message headers
	From = viper.GetString("mailclient.message.From")
	To = viper.GetString("mailclient.message.To")
	Subject = viper.GetString("mailclient.message.Subject")

	// SMTP headers
	Headers = viper.GetStringMapString("mailclient.message.headers")

	// Message Templates
	BodyTextTemplate = viper.GetString("mailclient.message.body.text")
	BodyHTMLTemplate = viper.GetString("mailclient.message.body.html")

	// Attachments and Embeds
	Attachments = viper.GetStringSlice("mailclient.message.attach")
	Embeds = viper.GetStringSlice("mailclient.message.embed")

	// Additional Exposed Template Metadata
	// tdata.URLTop = viper.GetString("mailclient.message.template-data.URLTop")

	log.Println("Setting up template data")
	staticTmplData := viper.GetStringMapString(
		"mailclient.message.template-data")
	log.Printf("staticTmplData %#v Len:%d\n", staticTmplData, len(staticTmplData))
	if len(staticTmplData) != 0 {
		if _, found := staticTmplData["dictionary"]; found {
			log.Printf("Dictionary key found\n")
			kvDict := viper.GetStringMapString(
				"mailclient.message.template-data.Dictionary")
			log.Printf("Dict: %#v\n", kvDict)

			// fmt.Println(kvDict)
			tdata.Dictionary = map[string]string{}
			for k, v := range kvDict {
				tdata.Dictionary[k] = v
				log.Printf("Dict KEY: %s : %s", k, v)
			}
		}
	}

	// Debug
	fmt.Printf("-= Connection Parameters =-")
	fmt.Printf("SMTP Server : %s\n", SMTPServer)
	fmt.Printf("SMTP Port   : %d\n", SMTPPort)
	fmt.Printf("SMTP User : %s\n", SMTPUser)
	fmt.Printf("SMTP TLS : %s\n", TLS)

	fmt.Printf("From: %s\n", From)
	fmt.Printf("To: %s\n", To)
	fmt.Printf("Subject: %s\n", Subject)

	fmt.Printf("Text Template: %s\n", BodyTextTemplate)
	fmt.Printf("HTML Template: %s\n", BodyHTMLTemplate)

	// Get SMTP credentials once
	SMTPPass, err = getUserCredentials(&SMTPServer)

	if err != nil {
		fmt.Printf("ERROR: Unable to record credentials: %v\n", err)
		os.Exit(2)
	}

	// Direct email, compose Mark and send
	if global.EmailRe.MatchString(To) {
		log.Printf("Mark is diredctly in the config. \n")
		var mark global.Mark
		mark.Firstname = viper.GetString("mailclient.message.mark.firstname")
		mark.Lastname = viper.GetString("mailclient.message.mark.lastname")
		mark.Identifier = viper.GetString("mailclient.message.mark.identifier")
		mark.Email = To

		tdata.Mark = &mark
		invokeRmail(&tdata)
	}

	fmt.Println(viper.GetString("mailclient.message.To"))
	// Marks in CSV file
	log.Printf("Marks in the database. \n")
	if global.DBFileRe.MatchString(viper.GetString("mailclient.message.To")) {
		var marks []global.Mark

		var settings = ql.ConnectionURL{
			Database: viper.GetString("mailclient.message.To"), // Path to database file.
		}
		// Attemping to open the "example.db" database file.
		sess, err := ql.Open(settings)
		if err != nil {
			log.Fatalf("db.Open(): %q\n", err)
		}
		defer sess.Close() // Remember to close the database session.

		log.Printf("Pointing to mark table \n")
		markCollection := sess.Collection("mark")

		// Let's query for the results we've just inserted.
		log.Printf("Querying for result : find()\n")
		res := markCollection.Find()

		log.Printf("Getting all results\n")
		err = res.All(&marks)
		if err != nil {
			log.Fatalf("res.All(): %q\n", err)
		}

		// Printing to stdout.
		log.Printf("-= Marks =-\n")
		for _, mark := range marks {
			fmt.Printf("Emailing: %s [id:%s] - %s %s\n",
				mark.Email,
				mark.Identifier,
				mark.Firstname,
				mark.Lastname,
			)
			tdata.Mark = &mark
			invokeRmail(&tdata)
			time.Sleep(5 * time.Second)
		}

	}

}

func invokeRmail(tdata *rmailer.TemplateData) {
	rmailer.GenMail(
		SMTPUser,
		SMTPPass,
		SMTPServer,
		SMTPPort,
		TLS,
		From,
		Subject,
		BodyTextTemplate, BodyHTMLTemplate,
		Attachments,
		Embeds,
		Headers,
		tdata)
}
