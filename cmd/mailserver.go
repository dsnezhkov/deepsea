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
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mhale/smtpd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/marcusolsson/tui-go"
)

var SMTPLocalServer string
var SMTPLocalPort int

type mailMsg struct {
	from    string
	subject string
	date    string
	body    string
}

var mailMsgs = make([]mailMsg, 0)

var mailserverCmd = &cobra.Command{
	Use:   "mailserver",
	Short: "Email server",
	Long:  `Email server for catching email threads`,
	Run: func(cmd *cobra.Command, args []string) {
		mailSDriver(cmd, args)
	},
}

func init() {

	// Bind command flags to configuration file directives.
	// So we can override them ont he command line

	// Example:
	// Cobra Flag name: SMTPPort
	// Viper key:       mailserver.SMTPPort
	// Variable :       SMTPPort
	// Bind Viper's key to Cobra's Flag name

	// mailserver.connection
	mailserverCmd.Flags().StringVarP(&SMTPLocalServer, "SMTPServer", "J",
		"127.0.0.1", "SMTP server")
	mailserverCmd.Flags().IntVarP(&SMTPLocalPort, "SMTPPort", "X",
		25, "SMTP server port")

	if err = viper.BindPFlag(
		"mailserver.connection.SMTPServer", mailserverCmd.Flags().Lookup("SMTPServer")); err != nil {
		_ = mailserverCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"mailserver.connection.SMTPPort", mailserverCmd.Flags().Lookup("SMTPPort")); err != nil {
		_ = mailserverCmd.Help()
		os.Exit(2)
	}


	rootCmd.AddCommand(mailserverCmd)
}

func mailSDriver(cmd *cobra.Command, args []string) {

	SMTPServer = viper.GetString("mailserver.connection.SMTPServer")
	SMTPPort = viper.GetInt("mailserver.connection.SMTPPort")

	// Debug
	fmt.Printf("SMTP Server : %s\n", SMTPServer)
	fmt.Printf("SMTP Port   : %d\n", SMTPPort)


	go func() {
		err = smtpd.ListenAndServe(strings.Join([]string{SMTPServer, strconv.Itoa(SMTPPort)}, ":"), mailHandler, "SMTP", "")
		if err != nil {
			log.Fatalf("Unable to start local SMTP server: %+v", err)
		}
	} ()



	go func () {

		for {

			reader := bufio.NewReader(os.Stdin)

			fmt.Print("-> ")
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)

			if strings.Compare("mail", text) == 0 {


					inbox := tui.NewTable(0, 0)
					inbox.SetColumnStretch(0, 3)
					inbox.SetColumnStretch(1, 2)
					inbox.SetColumnStretch(2, 1)
					inbox.SetFocused(true)

					for _, m := range mailMsgs {
						inbox.AppendRow(
							tui.NewLabel(m.subject),
							tui.NewLabel(m.from),
							tui.NewLabel(m.date),
						)
					}

					var (
						from    = tui.NewLabel("")
						subject = tui.NewLabel("")
						date    = tui.NewLabel("")
					)

					info := tui.NewGrid(0, 0)
					info.AppendRow(tui.NewLabel("From:"), from)
					info.AppendRow(tui.NewLabel("Subject:"), subject)
					info.AppendRow(tui.NewLabel("Date:"), date)

					body := tui.NewLabel("")
					body.SetSizePolicy(tui.Preferred, tui.Expanding)

					mailT := tui.NewVBox(info, body)
					mailT.SetSizePolicy(tui.Preferred, tui.Expanding)

					inbox.OnSelectionChanged(func(t *tui.Table) {
						if len(mailMsgs) != 0 {

							m := mailMsgs[t.Selected()]
							from.SetText(m.from)
							subject.SetText(m.subject)
							date.SetText(m.date)
							body.SetText(m.body)
						}

					})

					// Select first mail on startup.
					inbox.Select(0)

					root := tui.NewVBox(inbox, tui.NewLabel(""), mailT)

					ui, err := tui.New(root)
					if err != nil {
						log.Fatal(err)
					}

					ui.SetKeybinding("ESC", func() { ui.Quit() })

					if err := ui.Run(); err != nil {
						log.Fatal(err)
					}
			}
		}

	}()
	select {}
}

func startServer(IP string, Port int){

}

func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Printf("Error: %+v", err)
	}else{
		subject := msg.Header.Get("Subject")
		log.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)
		buf, err := ioutil.ReadAll(msg.Body)
		if err != nil {
			log.Fatal(err)
		}
		body := string(buf)

		message := mailMsg{
			from:    from,
			subject: subject,
			body:    body,
			date:  time.Now().String(),
		}

		mailMsgs = append(mailMsgs, message)
	}
}

