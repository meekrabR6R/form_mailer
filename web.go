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
	ModelEmailBody  string
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

	artistForm.SetWorks(form)

	err1 := artistForm.SetSignature(form["output"][0])
	if err1 != nil {
		panic(err1)
	}

	session, err2 := mgo.Dial(config.MongoUrl)
	if err2 != nil {
		panic(err2)
	}

	makeAPDF(artistForm)

	err3 := sendEmail("PERJUS Magazine release form",
		config.ArtistEmailBody,
		true,
		artistForm.Form)

	//THERE MUST BE A BETTER WAY!!!!!!
	for i := 0; i < len(artistForm.Works); i++ {
		for j := 0; j < len(artistForm.Works[i].Photos); j++ {
			for k := 0; k < len(artistForm.Works[i].Photos[j].Models); k++ {
				modelErr := sendEmail("PERJUS Magazine model release form",
					config.ModelEmailBody,
					false,
					artistForm.Works[i].Photos[j].Models[k].Form)

				snt := true
				if modelErr != nil {
					panic(modelErr)
					snt = false
				}

				artistForm.Works[i].Photos[j].Models[k].EmailSent = snt
			}
		}
	}

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
