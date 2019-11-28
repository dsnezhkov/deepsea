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
	"os"
	"strings"
	"syscall"
	"time"
	"upper.io/db.v3/ql"

	"deepsea/global"
	"deepsea/rmailer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	jlog "github.com/spf13/jwalterweatherman"
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
	Long:  `MAILCLIENT: Connect and Email a phish with features`,
	Run: func(cmd *cobra.Command, args []string) {
		jlog.DEBUG.Println("mailDriver()")
		mailDriver(cmd, args)
	},
}

func init() {

	// Connection
	mailclientCmd.Flags().StringVarP(
		&SMTPServer,
		"SMTPServer",
		"s",
		"127.0.0.1",
		"SMTP server")

	mailclientCmd.Flags().IntVarP(
		&SMTPPort,
		"SMTPPort",
		"p",
		25,
		"SMTP server port")

	mailclientCmd.Flags().StringVarP(
		&SMTPUser,
		"SMTPUser",
		"U",
		"testuser",
		"SMTP user")

	mailclientCmd.Flags().StringVarP(
		&TLS,
		"TLS",
		"t",
		"yes",
		"Use TLS handshake (STARTTLS)")

	if err = viper.BindPFlag(
		"mailcient.connection.SMTPServer",
		mailclientCmd.Flags().Lookup("SMTPServer")); err != nil {
		jlog.DEBUG.Println("Setting SMTPServer")
		_ = mailclientCmd.Help()
		os.Exit(2)
	}

	if err = viper.BindPFlag(
		"mailclient.connection.SMTPPort",
		mailclientCmd.Flags().Lookup("SMTPPort")); err != nil {
		jlog.DEBUG.Println("Setting SMTPPort")
		_ = mailclientCmd.Help()
		os.Exit(2)
	}

	if err = viper.BindPFlag(
		"mailclient.connection.SMTPUser",
		mailclientCmd.Flags().Lookup("SMTPUser")); err != nil {
		jlog.DEBUG.Println("Setting SMTPUser")
		_ = mailclientCmd.Help()
		os.Exit(2)
	}

	if err = viper.BindPFlag(
		"mailclient.connection.TLS",
		mailclientCmd.Flags().Lookup("TLS")); err != nil {
		jlog.DEBUG.Println("Setting TLS")
		_ = mailclientCmd.Help()
		os.Exit(2)
	}

	// Message
	mailclientCmd.Flags().StringVarP(
		&From,
		"From",
		"F",
		"",
		"Message From: header")

	mailclientCmd.Flags().StringVarP(
		&To,
		"To",
		"T",
		"",
		"Message To: header")

	mailclientCmd.Flags().StringVarP(
		&Subject,
		"Subject",
		"S",
		"",
		"Message Subject: header")

	// message body
	mailclientCmd.Flags().StringVarP(
		&BodyHTMLTemplate,
		"HTMLTemplate",
		"H",
		"", "HTML Template file (.htpl)")

	mailclientCmd.Flags().StringVarP(
		&BodyTextTemplate,
		"TextTemplate",
		"P",
		"", "Text Template file (.ttpl)")

	if err = viper.BindPFlag(
		"mailclient.message.body.html",
		mailclientCmd.Flags().Lookup("HTMLTemplate")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}

	if err = viper.BindPFlag(
		"mailclient.message.body.text",
		mailclientCmd.Flags().Lookup("TextTemplate")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}

	if err = viper.BindPFlag(
		"mailclient.message.From",
		mailclientCmd.Flags().Lookup("From")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}

	if err = viper.BindPFlag(
		"mailclient.message.To",
		mailclientCmd.Flags().Lookup("To")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}

	if err = viper.BindPFlag(
		"mailclient.message.Subject",
		mailclientCmd.Flags().Lookup("Subject")); err != nil {
		_ = mailclientCmd.Help()
		os.Exit(2)
	}
	rootCmd.AddCommand(mailclientCmd)
}

// Processing
func getUserCredentials(server *string) (string, error) {

	jlog.INFO.Printf("SMTP Auth Creds for %s: \n", *server)
	fmt.Println("\nEnter Password: ")

	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		jlog.DEBUG.Println("Password sourcing error")
		return "", err
	}
	password := string(bytePassword)
	jlog.DEBUG.Println("Returning trimmed password")
	jlog.TRACE.Printf("PW: %s\n", password)
	return strings.TrimSpace(password), nil
}

func mailDriver(cmd *cobra.Command, args []string) {

	// Connectivity
	SMTPServer = viper.GetString("mailclient.connection.SMTPServer")
	SMTPPort = viper.GetInt("mailclient.connection.SMTPPort")
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

	if len(TLS) == 0 {
		jlog.ERROR.Fatalln("TLS: cannot be empty.")
	}
	if len(Subject) == 0 {
		jlog.ERROR.Fatalln("Subject: cannot be empty.")
	}
	if len(From) == 0 {
		jlog.ERROR.Fatalln("From: cannot be empty")
	}
	if len(To) == 0 {
		jlog.ERROR.Fatalln("To: cannot be empty")
	}
	if len(BodyTextTemplate) == 0 {
		jlog.ERROR.Fatalln("BodyTextTemplate: cannot be empty")
	}
	if len(BodyHTMLTemplate) == 0 {
		jlog.ERROR.Fatalln("BodyHTMLTemplate: cannot be empty")
	}

	// Additional Exposed Template Metadata
	jlog.INFO.Println("Setting up template data")
	staticTmplData := viper.GetStringMapString(
		"mailclient.message.template-data")

	if len(staticTmplData) != 0 {
		jlog.DEBUG.Println("Template data is present")
		if _, found := staticTmplData["dictionary"]; found {
			jlog.DEBUG.Println("Template data [dictionary] is present")
			kvDict := viper.GetStringMapString(
				"mailclient.message.template-data.Dictionary")
			jlog.DEBUG.Printf("Template data: %#v\n", kvDict)

			tdata.Dictionary = map[string]string{}
			for k, v := range kvDict {
				tdata.Dictionary[k] = v
				jlog.DEBUG.Printf("Template dictionary layout: %s : %s", k, v)
			}
		}
	}

	jlog.DEBUG.Println("-= = = = Connection Parameters = = = =-")
	jlog.DEBUG.Printf("SMTP Server : %s\n", SMTPServer)
	jlog.DEBUG.Printf("SMTP Port   : %d\n", SMTPPort)
	jlog.DEBUG.Printf("SMTP TLS    : %s\n", TLS)
	jlog.DEBUG.Printf("SMTP User   : %s\n", SMTPUser)
	jlog.DEBUG.Printf("SMTP Pass(n): %d\n", len(SMTPPass))

	jlog.DEBUG.Println("-= = = = Message Envelope Parameters = = = =-")
	jlog.DEBUG.Printf("From   : %s\n", From)
	jlog.DEBUG.Printf("To     : %s\n", To)
	jlog.DEBUG.Printf("Subject: %s\n", Subject)

	jlog.DEBUG.Println("-= = = = Message BODY Parameters = = = =-")
	jlog.DEBUG.Printf("HTML Template (f): %s\n", BodyHTMLTemplate)
	jlog.DEBUG.Printf("Text Template (f): %s\n", BodyTextTemplate)

	if global.EmailRe.MatchString(To) {
		jlog.DEBUG.Println("Delivery to a single email.")

		var mark global.Mark
		mark.Firstname = viper.GetString("mailclient.message.mark.firstname")
		mark.Lastname = viper.GetString("mailclient.message.mark.lastname")
		mark.Identifier = viper.GetString("mailclient.message.mark.identifier")
		mark.Email = To

		tdata.Mark = &mark
		invokeRmail(&tdata)
	}

	// Marks in CSV file
	jlog.DEBUG.Printf("Delivery to a list of marks. \n")
	if global.DBFileRe.MatchString(
		viper.GetString("mailclient.message.To")) {

		var marks []global.Mark
		var settings = ql.ConnectionURL{
			Database: viper.GetString("mailclient.message.To"),
		}

		sess, err := ql.Open(settings)
		if err != nil {
			jlog.ERROR.Fatalf("db.Open(): %q\n", err)
		}
		defer sess.Close()

		jlog.DEBUG.Println("Pointing to mark table")
		markCollection := sess.Collection("mark")

		// Let's query for the results we've just inserted.
		jlog.DEBUG.Println("Querying for result : find()")
		res := markCollection.Find()

		jlog.DEBUG.Println("Getting all results")
		err = res.All(&marks)
		if err != nil {
			jlog.ERROR.Fatalf("res.All(): %q\n", err)
		}

		jlog.DEBUG.Println("-= = = = Marks = = = =-")
		for _, mark := range marks {
			jlog.INFO.Printf(" M => Dialer: %s [id:%s] - %s %s\n",
				mark.Email,
				mark.Identifier,
				mark.Firstname,
				mark.Lastname,
			)
			tdata.Mark = &mark
			invokeRmail(&tdata)
			jlog.DEBUG.Println("Sleeping")
			time.Sleep(5 * time.Second)
		}
	}
}

func invokeRmail(tdata *rmailer.TemplateData) {
	m, err := rmailer.GenMail(
		SMTPServer,
		From,
		Subject,
		BodyTextTemplate, BodyHTMLTemplate,
		Attachments,
		Embeds,
		Headers,
		tdata)

	if err != nil {
		jlog.ERROR.Fatalln("Mail Generation Error: ", err)
	}

	// Get SMTP credentials once
	// TODO: Implement Dry-run, with only final message generation but no send
	if len(SMTPPass) == 0 {
		// Ask and Cache password
		SMTPPass, err = getUserCredentials(&SMTPServer)
		if err != nil {
			jlog.ERROR.Fatalf("Unable to record credentials: %v\n", err)
		}
	}
	jlog.INFO.Print("Sending Email ... ")
	rmailer.DialSend(m, SMTPServer, SMTPPort, SMTPUser, SMTPPass, TLS)
	jlog.INFO.Print("OK")
}
