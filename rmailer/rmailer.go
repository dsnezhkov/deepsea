package rmailer

import (
	"crypto/tls"
	"fmt"
	"deepsea/global"
	gomail "github.com/gophish/gomail"
	thtml "html/template"
	ttext "html/template"
	"io"
	"path/filepath"
	"strings"
)

func GenMail(username string, password string, server string, port int,
	usetls string, from string, subject string,
	bodyTextTemplate string, bodyHtmlTemplate string, attachments []string,
	embeds []string, headers map[string]string, tdata *TemplateData) {

	var err error

	fmt.Printf("Identifier: %s | Email: %s | First Name: %s | Last Name: %s\n",
		tdata.Mark.Identifier, tdata.Mark.Email,
		tdata.Mark.Firstname, tdata.Mark.Lastname)

	m := gomail.NewMessage()
	m.SetHeader("Subject", subject)
	m.SetHeader("From", from)
	m.SetHeader("To", tdata.Mark.Email)

	// Set Headers
	for key, value := range headers {
		m.SetHeader(key, value)
	}

	// Create a Message-Id:
	msg_id := strings.Join([]string{ global.RandString(16), server}, "@")
	m.SetHeader("Message-ID", "<"+msg_id+">")

	// Templates HTML/Text
	th, err := thtml.ParseFiles(bodyHtmlTemplate)
	if err != nil {
		fmt.Printf("ERROR: HTML Template not parsed %v\n", err)
		return
	}
	tt, err := ttext.ParseFiles(bodyTextTemplate)
	if err != nil {
		fmt.Printf("ERROR: Text Template not parsed %v\n", err)
		return
	}

	// Embedded images
	l := len(embeds)
	if l != 0 {
		tdata.EmbedImage = make([]string, l)
		for ix, file := range embeds {
			fmt.Println("Embedding: ", file)

			if global.FileExists(file) {
				m.Embed(file)
				tdata.EmbedImage[ix] = filepath.Base(file)
			} else {
				fmt.Println("File Not Found: ", file)
			}
		}
	}

	// URLs
	// Compile and write Templates
	m.AddAlternativeWriter("text/plain", func(w io.Writer) error {
		return tt.Execute(w, tdata)
	})
	m.AddAlternativeWriter("text/html", func(w io.Writer) error {
		return th.Execute(w, tdata)
	})

	// Attachments
	for _, file := range attachments {
		fmt.Println("Attaching: ", file)
		if global.FileExists(file) {
			m.Attach(file)
		} else {
			fmt.Println("File Not Found: ", file)
		}
	}
    fmt.Println(m)
	dialSend(m, server, port, username, password, usetls)
}

func dialSend(m *gomail.Message, server string, port int, username string, password string, usetls string) {

	fmt.Println(m)
	d := gomail.NewDialer(server, port, username, password)
	if strings.ToLower(usetls) == "yes" {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // no idea who we are connecting. relax
	}

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("ERROR: Could not dial and send: %v", err)
	}
}
