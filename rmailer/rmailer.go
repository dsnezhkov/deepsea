package rmailer

import (
	"crypto/tls"
	"deepsea/global"
	"errors"
	gomail "github.com/gophish/gomail"
	jlog "github.com/spf13/jwalterweatherman"
	thtml "html/template"
	ttext "html/template"
	"io"
	"path/filepath"
	"strings"
)

func GenMail(
	server string,
	from string,
	subject string,
	bodyTextTemplate string,
	bodyHtmlTemplate string,
	attachments []string,
	embeds []string,
	headers map[string]string,
	tdata *TemplateData) (*gomail.Message, error) {

	var err error

	jlog.DEBUG.Printf("Identifier: %s | Email: %s | First Name: %s | Last Name: %s\n",
		tdata.Mark.Identifier, tdata.Mark.Email,
		tdata.Mark.Firstname, tdata.Mark.Lastname)

	jlog.DEBUG.Println("Creating mail message")
	m := gomail.NewMessage()

	jlog.DEBUG.Println("Setting Headers")
	m.SetHeader("Subject", subject)
	m.SetHeader("From", from)
	m.SetHeader("To", tdata.Mark.Email)

	// Set Headers
	for key, value := range headers {
		jlog.TRACE.Printf("%s => %s\n", key, value)
		m.SetHeader(key, value)
	}

	// Create a Message-Id:
	jlog.TRACE.Println("Creating Message-ID")
	msgId := strings.Join([]string{global.RandString(16), server}, "@")
	m.SetHeader("Message-ID", "<"+msgId+">")

	// Templates HTML/Text
	jlog.DEBUG.Println("Parsing HTML Template")
	th, err := thtml.ParseFiles(bodyHtmlTemplate)
	if err != nil {
		jlog.ERROR.Println("Parsing HTML Template: ERROR")
		return new(gomail.Message), err
	}

	tt, err := ttext.ParseFiles(bodyTextTemplate)
	if err != nil {
		jlog.ERROR.Println("Parsing TEXT Template: ERROR")
		return new(gomail.Message), err
	}

	// Embedded images
	jlog.DEBUG.Println("Embedding Images")
	l := len(embeds)
	if l != 0 {
		tdata.EmbedImage = make([]string, l)
		for ix, file := range embeds {
			jlog.INFO.Println("Embedding: ", file)

			if global.FileExists(file) {
				jlog.DEBUG.Println("File Exists")
				m.Embed(file)
				tdata.EmbedImage[ix] = filepath.Base(file)
			} else {
				jlog.DEBUG.Println("Embedding failed")
				return new(gomail.Message),
					errors.New("Embedded File Not Found: " + file)
			}
		}
	}

	// URLs
	// Compile and write Templates
	jlog.INFO.Println("Merging Templates")
	jlog.DEBUG.Println("Merging Templates: TEXT")
	m.AddAlternativeWriter("text/plain", func(w io.Writer) error {
		return tt.Execute(w, tdata)
	})

	jlog.DEBUG.Println("Merging Templates: TEXT")
	m.AddAlternativeWriter("text/html", func(w io.Writer) error {
		return th.Execute(w, tdata)
	})

	// Attachments
	jlog.INFO.Println("Attaching Files (if any)")
	for _, file := range attachments {
		jlog.INFO.Println("Attaching asset: ", file)
		if global.FileExists(file) {
			m.Attach(file)
		} else {
			jlog.DEBUG.Println("Attaching failed")
			return new(gomail.Message),
				errors.New("Attachment File Not Found: " + file)
		}
	}
	return m, nil
}

func DialSend(
	m *gomail.Message,
	server string, port int, username string, password string, usetls string) {

	jlog.DEBUG.Println("Creating Dialer")
	d := gomail.NewDialer(server, port, username, password)
	if strings.ToLower(usetls) == "yes" {
		jlog.DEBUG.Println("Setting TLS")
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	jlog.DEBUG.Println("Dialing Now...")
	if err := d.DialAndSend(m); err != nil {
		jlog.ERROR.Fatalf("Could not dial and send: %v\n", err)
	}
}
