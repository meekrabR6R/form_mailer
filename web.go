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
	MongoUrl    string
	DbName      string
	SenderEmail string
	SenderPass  string
	ArtistEmail string
	ArtistTitle string
	ArtistBody  string
	ModelTitle  string
	ModelBody   string
}

func FormHandler(w http.ResponseWriter, req *http.Request) {
	var config = getConf()
	err0 := req.ParseForm()

	if err0 != nil {
		panic(err0)
	}

	form := req.PostForm
	rawSig := []byte(form["output"][0])
	var sig []map[string]int

	err1 := json.Unmarshal(rawSig, &sig)

	if err1 != nil {
		panic(err1)
	}

	artist := Artist{
		FirstName: form["firstName"][0],
		LastName:  form["lastName"][0],
		Email:     form["emailAddress"][0],
		Link:      form["downloadLink"][0],
		Sig:       sig,
	}

	session, err2 := mgo.Dial(config.MongoUrl)

	if err2 != nil {
		panic(err2)
	}

	//experimenting with returning the error
	//or just handling it internally in the
	//function.
	makeAPDF(artist)
	err3 := sendEmail(artist)

	sent := true
	if err3 != nil {
		panic(err3)
		sent = false
	}

	artist.EmailSent = sent
	artistsCollection := session.DB(config.DbName).C("artists")
	artistsCollection.Insert(artist)
}

func makeAPDF(artist Artist) {
	var config = getConf()
	body := fmt.Sprintf(config.ArtistBody, strings.ToUpper(artist.FullName()),
		strings.ToUpper(artist.FullName()))

	pdfBody := fmt.Sprintf("%s\n\n%s", config.ArtistTitle, body)

	pdf := gofpdf.New("P", "mm", "A4", "../font")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 10)
	pdf.MultiCell(185, 5, pdfBody, "", "", false)

	err := pdf.OutputFileAndClose(
		fmt.Sprintf("%s_release.pdf", strings.ToLower(artist.LastName)))

	if err != nil {
		panic(err)
	}
}

func sendEmail(artist Artist) error {
	var config = getConf()
	e := &email.Email{
		To:      []string{artist.Email},
		From:    fmt.Sprintf("Perjus <%s>", config.SenderEmail),
		Subject: "PERJUS Magazine release forms",
		Text:    []byte(config.ArtistEmail),
		HTML:    []byte(fmt.Sprintf("<h1>%s</h1>", config.ArtistEmail)),
		Headers: textproto.MIMEHeader{},
	}

	e.AttachFile(fmt.Sprintf("%s_release.pdf",
		strings.ToLower(artist.LastName)))

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
