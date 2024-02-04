package helpers

import (
	"bytes"
	"html/template"
	"log"
	"os"
)

func MailHTMLTemplate(path string, vars interface{}) string {
	var buf bytes.Buffer

	baseDir, _ := os.Getwd()

	tmpl, err := template.ParseFiles(baseDir + "/Layout/Email/" + path)
	if err != nil {
		log.Panicf("Error parsing template: %v", err)
	}
	tmpl.Execute(&buf, vars)
	return buf.String()
}

func PDFHTMLTemplate(path string, vars interface{}) bytes.Buffer {
	var buf bytes.Buffer

	baseDir, _ := os.Getwd()

	tmpl, err := template.ParseFiles(baseDir + "/Layout/PDF/" + path)
	if err != nil {
		log.Panicf("Error parsing template: %v", err)
	}
	tmpl.Execute(&buf, vars)
	return buf
}
