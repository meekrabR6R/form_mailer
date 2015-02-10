package main

import (
	"code.google.com/p/gofpdf"
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"
)

/**
 * Global config
 */
type Config struct {
	Url               string
	MongoUrl          string
	DbName            string
	SenderEmail       string
	SenderPass        string
	ArtistEmailBody   string
	ArtistTitle       string
	ArtistBody        string
	ModelEmailBodyOne string
	ModelTitle        string
	ModelBody         string
}

type Content struct {
	ArtistId string
	ModelId  string
	Model    ModelForm
	Conf     *Config
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

func makeArtistForm(form map[string][]string) (error, *ArtistForm) {
	artistForm := &ArtistForm{
		Form: Form{
			FirstName: form["firstName"][0],
			LastName:  form["lastName"][0],
			Email:     form["emailAddress"][0],
			Link:      form["downloadLink"][0],
		},
	}

	artistForm.SetWorks(form)
	writeNewMetaData(artistForm)
	err := artistForm.SetSignature(form["output"][0])

	return err, artistForm
}

func writeArtistFormToDb(url string, sent bool, artistForm *ArtistForm) error {
	artistForm.EmailSent = sent
	err, artistFormsCollection := makeOrGetCollection("artistForms")
	artistFormsCollection.Insert(artistForm)

	return err
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

func sendErrorEmail(err error) {
	sendEmail("Error Report", err.Error(), false, Form{Email: "nmiano84@gmail.com"})
}

func makeOrGetCollection(coll string) (error, *mgo.Collection) {
	config := getConf()
	session, err := mgo.Dial(config.MongoUrl)
	return err, session.DB(config.DbName).C("artistForms")
}

func getArtistFromCollection(id bson.ObjectId) (error, ArtistForm) {
	err, artistForms := makeOrGetCollection("artistForms")

	if err != nil {
		panic(err)
	}

	artistForm := ArtistForm{}
	query := bson.M{"works.photos.models._id": id}

	err1 := artistForms.Find(query).One(&artistForm)

	return err1, artistForm
}

func makeContent(id bson.ObjectId, artistForm ArtistForm) Content {
	model := artistForm.ModelById(id)

	return Content{
		ArtistId: artistForm.Id.Hex(),
		ModelId:  id.Hex(),
		Model:    model,
		Conf:     getConf(),
	}
}
