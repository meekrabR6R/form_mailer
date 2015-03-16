package main

import (
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"
	"time"
)

/**
 * Global config
 */
type Config struct {
	Url                string
	MongoUrl           string
	MongoUser          string
	MongoPass          string
	DbName             string
	SenderEmail        string
	SenderPass         string
	AdminBodyForArtist string
	AdminBodyForModel  string
	ArtistEmailBody    string
	ArtistTitle        string
	ArtistBody         string
	ModelEmailBodyOne  string
	ModelTitle         string
	ModelBody          string
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

func makeModelPDF(form BaseForm) {
	makeAPDF(form, 100, 520)
}

func makeArtistPDF(form BaseForm) {
	makeAPDF(form, 100, 720)
}

func makeAPDF(form BaseForm, x float64, y float64) {
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

	var body string
	if len(form.GetDataAsString()) > 0 {
		body = fmt.Sprintf(content, strings.ToUpper(form.FullName()),
			form.GetDataAsString())
	} else {
		body = fmt.Sprintf(content, strings.ToUpper(form.FullName()))
	}

	//time formatting
	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	t := time.Now()
	//hacky formtting... but i'm so tired..
	pdfBody := fmt.Sprintf("%s\n\n%s\n\n\nDate: %s\n\nSignature:",
		title,
		body,
		t.Format(layout))

	pdf := gofpdf.New("P", "mm", "A4", "../font")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 10)
	pdf.MultiCell(185, 5, pdfBody, "", "", false)

	//write sig
	for i := 0; i < len(form.GetSignature()); i++ {
		dot := form.GetSignature()[i]
		pdf.Line((dot["lx"]+x)/4,
			(dot["ly"]+y)/4,
			(dot["mx"]+x)/4,
			(dot["my"]+y)/4)
	}

	err := pdf.OutputFileAndClose(
		fmt.Sprintf("%s_release.pdf", form.FullNameForFile()))

	if err != nil {
		//panic(err)
	}
}

func writeArtistFormToDb(url string, sent bool, artistForm *ArtistForm) error {
	artistForm.EmailSent = sent
	err, artistFormsCollection := makeOrGetCollection("artistForms")
	artistFormsCollection.Insert(artistForm)

	return err
}

func makeOrGetCollection(coll string) (error, *mgo.Collection) {
	config := getConf()

	mongoUrl := fmt.Sprintf("mongodb://%s:%s/",
		os.Getenv("OPENSHIFT_MONGODB_DB_HOST"), os.Getenv("OPENSHIFT_MONGODB_DB_PORT"))
	session, err := mgo.Dial(mongoUrl)
	session.DB(config.DbName).Login(config.MongoUser, config.MongoPass)
	return err, session.DB(config.DbName).C(coll)
}

func getArtistFromCollection(id bson.ObjectId) (error, ArtistForm) {
	err, artistForms := makeOrGetCollection("artistForms")

	if err != nil {
		//panic(err)
	}

	artistForm := ArtistForm{}
	query := bson.M{"works.photos.models._id": id}

	err1 := artistForms.Find(query).One(&artistForm)

	return err1, artistForm
}

func makeContent(id bson.ObjectId, artistForm ArtistForm) Content {
	model := *artistForm.ModelById(id)

	return Content{
		ArtistId: artistForm.Id.Hex(),
		ModelId:  id.Hex(),
		Model:    model,
		Conf:     getConf(),
	}
}

/**
 * Email needs to be its own parameter so this function
 * can be used to send admin email too
 */
func sendEmail(emailAddress string, sub string, bod string,
	attachPdf bool, form BaseForm) error {

	var config = getConf()
	e := &email.Email{
		To:      []string{emailAddress},
		From:    fmt.Sprintf("PERJUS <%s>", config.SenderEmail),
		Subject: sub,
		Text:    []byte(bod),
		HTML:    []byte(fmt.Sprintf("<h1>%s</h1>", bod)),
		Headers: textproto.MIMEHeader{},
	}
	if attachPdf {
		e.AttachFile(fmt.Sprintf("%s_release.pdf", form.FullNameForFile()))
	}

	return e.Send("smtp.gmail.com:587",
		smtp.PlainAuth("", config.SenderEmail, config.SenderPass,
			"smtp.gmail.com"))
}

/**
 * Wrapper function for generating pdf and sending email
 * to artist
 */
func sendArtistEmail(artistForm *ArtistForm) (error, bool) {
	config := getConf()
	makeArtistPDF(artistForm)
	err := sendEmail(artistForm.Form.Email,
		"PERJUS Magazine release form",
		config.ArtistEmailBody,
		true,
		artistForm)

	sent := true
	if err != nil {
		//panic(err3)
		sent = false
	}

	return err, sent
}

func sendAdminEmailForArtist(form *ArtistForm) (error, bool) {
	config := getConf()
	makeArtistPDF(form)
	return sendAdminEmail(form,
		fmt.Sprintf(config.AdminBodyForArtist,
			form.FullName(),
			form.Link))
}

func sendAdminEmailForModel(form *ModelForm) (error, bool) {
	config := getConf()
	makeModelPDF(form)
	return sendAdminEmail(form,
		fmt.Sprintf(config.AdminBodyForModel,
			form.FullName()))
}

func sendAdminEmail(form BaseForm, body string) (error, bool) {
	config := getConf()
	err := sendEmail(config.SenderEmail,
		fmt.Sprintf("%s - Signed Release Forms - Issue %d",
			strings.ToUpper(form.FullName()),
			2), //temp holder for version
		body,
		true,
		form)
	sent := true
	if err != nil {
		sent = false
	}
	return err, sent
}

func sendModelEmailWithLink(artistForm *ArtistForm, modelForm *ModelForm) (error, bool) {
	config := getConf()

	url := fmt.Sprintf("%s/models/%s/release", config.Url,
		modelForm.Id.Hex())

	modelErr := sendEmail(modelForm.Email,
		"PERJUS Magazine model release form",
		fmt.Sprintf(config.ModelEmailBodyOne,
			strings.ToUpper(artistForm.FullName()),
			url),
		false,
		modelForm)

	snt := true
	if modelErr != nil {
		//panic(modelErr)
		snt = false
	}

	return modelErr, snt
}

func sendModelEmailWithForm(modelForm *ModelForm) (error, bool) {
	config := getConf()
	makeModelPDF(modelForm)
	err := sendEmail(modelForm.Form.Email,
		"PERJUS Magazine release form",
		config.ArtistEmailBody,
		true,
		modelForm)

	sent := true
	if err != nil {
		//panic(err3)
		sent = false
	}

	return err, sent
}

/**
 * Ugly little helper function that sends emails
 * to each model that belongs to an ArtistForm.
 * It is insanely inefficient, and could use some
 * optimization :)
 * (Should ideally be run in its own goroutine)
 */
func sendAllModelEmails(artistForm *ArtistForm) {
	//THERE MUST BE A BETTER WAY!!!!!!
	for i := 0; i < len(artistForm.Works); i++ {
		for j := 0; j < len(artistForm.Works[i].Photos); j++ {
			for k := 0; k < len(artistForm.Works[i].Photos[j].Models); k++ {
				//At least this will mitigate slowdown due to O(n^3) complexity somewhat! :-P
				go func(iIdx int, jIdx int, kIdx int) {
					modelErr, snt := sendModelEmailWithLink(artistForm,
						&artistForm.Works[iIdx].Photos[jIdx].Models[kIdx])
					if modelErr != nil {
						//panic(modelErr)
					}

					artistForm.Works[iIdx].Photos[jIdx].Models[kIdx].EmailSent = snt
				}(i, j, k)
			}
		}
	}
}
