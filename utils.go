package main

import (
	"code.google.com/p/gofpdf"
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"
)

/**
 * Global config
 */
type Config struct {
	MongoUrl        string
	DbName          string
	SenderEmail     string
	SenderPass      string
	ArtistEmailBody string
	ArtistTitle     string
	ArtistBody      string
	ModelEmailBody  string
	ModelTitle      string
	ModelBody       string
}

func getConf() (conf *Config) {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	conf = &Config{}
	decoder.Decode(&conf)
	return
}

func makeAPDF(form BaseForm) {
	var config = getConf()
	var title string
	var content string

	if form.IsArtist() {
		title = config.ArtistTitle
		content = config.ArtistBody
	} else {
		title = config.ModelTitle
		content = config.ModelBody
	}

	body := fmt.Sprintf(content, strings.ToUpper(form.FullName()),
		strings.ToUpper(form.FullName()))

	pdfBody := fmt.Sprintf("%s\n\n%s", title, body)

	pdf := gofpdf.New("P", "mm", "A4", "../font")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 10)
	pdf.MultiCell(185, 5, pdfBody, "", "", false)

	err := pdf.OutputFileAndClose(
		fmt.Sprintf("%s_release.pdf", form.FullName()))

	if err != nil {
		panic(err)
	}
}

func sendEmail(sub string, bod string, attachPdf bool, form Form) error {
	var config = getConf()
	e := &email.Email{
		To:      []string{form.Email},
		From:    fmt.Sprintf("PERJUS <%s>", config.SenderEmail),
		Subject: sub,
		Text:    []byte(bod),
		HTML:    []byte(fmt.Sprintf("<h1>%s</h1>", bod)),
		Headers: textproto.MIMEHeader{},
	}
	if attachPdf {
		e.AttachFile(fmt.Sprintf("%s_release.pdf",
			strings.ToLower(form.FullName())))
	}

	return e.Send("smtp.gmail.com:587",
		smtp.PlainAuth("", config.SenderEmail, config.SenderPass,
			"smtp.gmail.com"))
}
