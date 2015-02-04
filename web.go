package main

import (
	"code.google.com/p/gofpdf"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jordan-wright/email"
	"gopkg.in/mgo.v2"
	"net/http"
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
	ModelTitle      string
	ModelBody       string
}

func FormHandler(w http.ResponseWriter, req *http.Request) {
	var config = getConf()
	err0 := req.ParseForm()

	if err0 != nil {
		panic(err0)
	}

	form := req.PostForm

	artistForm := &ArtistForm{
		Form: Form{
			FirstName: form["firstName"][0],
			LastName:  form["lastName"][0],
			Email:     form["emailAddress"][0],
			Link:      form["downloadLink"][0],
		},
	}
	err1 := artistForm.SetSignature(form["output"][0])
	if err1 != nil {
		panic(err1)
	}

	session, err2 := mgo.Dial(config.MongoUrl)

	if err2 != nil {
		panic(err2)
	}

	//experimenting with returning the error
	//or just handling it internally in the
	//function.
	makeAPDF(artistForm)
	err3 := sendEmail(artistForm)

	sent := true
	if err3 != nil {
		panic(err3)
		sent = false
	}

	artistForm.EmailSent = sent
	artistFormsCollection := session.DB(config.DbName).C("artistForms")
	artistFormsCollection.Insert(artistForm)
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

func sendEmail(artistForm *ArtistForm) error {
	var config = getConf()
	e := &email.Email{
		To:      []string{artistForm.Email},
		From:    fmt.Sprintf("Perjus <%s>", config.SenderEmail),
		Subject: "PERJUS Magazine release forms",
		Text:    []byte(config.ArtistEmailBody),
		HTML:    []byte(fmt.Sprintf("<h1>%s</h1>", config.ArtistEmailBody)),
		Headers: textproto.MIMEHeader{},
	}

	e.AttachFile(fmt.Sprintf("%s_release.pdf",
		strings.ToLower(artistForm.FullName())))

	return e.Send("smtp.gmail.com:587",
		smtp.PlainAuth("", config.SenderEmail, config.SenderPass, "smtp.gmail.com"))
}

func getConf() (conf *Config) {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	conf = &Config{}
	decoder.Decode(&conf)
	return
}

func main() {
	portPtr := flag.String("port", "foo", "port number")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/work", FormHandler)
	//Must go after all routes..
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", router)

	bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

	if *portPtr != "foo" {
		fmt.Printf("listening on %s...\n", *portPtr)
		port := fmt.Sprintf(":%s", *portPtr)
		err := http.ListenAndServe(port, nil)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("listening on %s...", bind)
		err := http.ListenAndServe(bind, nil)
		if err != nil {
			panic(err)
		}
	}
}
