package rmailer

import "deepsea/global"

type TemplateData struct {
	EmbedImage []string
	Mark       *global.Mark
	Dictionary map[string]string
}
