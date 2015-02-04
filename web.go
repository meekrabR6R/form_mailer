package main

import (
	"code.google.com/p/gofpdf"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jordan-wright/email"
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
	SenderEmail string
	SenderPass  string
	ArtistTitle string
	ArtistBody  string
	ModelTitle  string
	ModelBody   string
}

func FormHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		fmt.Println(err)
	}
	form := req.PostForm

	artist := Artist{
		FirstName: form["firstName"][0],
		LastName:  form["lastName"][0],
		Email:     form["emailAddress"][0],
		Link:      form["downloadLink"][0],
	}
	//experimenting with returning the error
	//or just handling it internally in the
	//
	makeAPDF(artist)
	err2 := sendEmail(artist)

	if err2 != nil {
		panic(err2)
	}
}

func makeAPDF(artist Artist) {
	var config = getConf()
	body := fmt.Sprintf(config.ArtistBody, strings.ToUpper(artist.FullName()),
		strings.ToUpper(artist.FullName()))
	pdfBody := fmt.Sprintf("%s\n\n%s", config.ArtistTitle, body)
	//config.ArtistTitle,
	//fmt.Sprintf(config.ArtistBody, artist.FullName()))

	pdf := gofpdf.New("P", "mm", "A4", "../font")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 10)
	pdf.MultiCell(150, 5, pdfBody, "", "", false)

	err := pdf.OutputFileAndClose(
		fmt.Sprintf("%s_release.pdf", artist.LastName))

	if err != nil {
		panic(err)
	}
}

func sendEmail(artist Artist) error {
	var config = getConf()
	e := &email.Email{
		To:   []string{artist.Email},
		From: fmt.Sprintf("Perjus <%s>", config.SenderEmail),
		Subject: fmt.Sprintf("Artist Release Form for %s",
			artist.FullName()),
		Text:    []byte("Text Body is, of course, supported!"),
		HTML:    []byte("<h1>Fancy HTML is supported, too!</h1>"),
		Headers: textproto.MIMEHeader{},
	}
	e.AttachFile("temp1.pdf") //temporary!!

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
