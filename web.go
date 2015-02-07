package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"html/template"
	"net/http"
	"os"
)

func FormHandler(w http.ResponseWriter, req *http.Request) {
	var config = getConf()
	err0 := req.ParseForm()
	if err0 != nil {
		panic(err0)
	}

	form := req.PostForm

	err1, artistForm := makeArtistForm(form)
	if err1 != nil {
		panic(err1)
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

	session, err2 := mgo.Dial(config.MongoUrl)
	if err2 != nil {
		panic(err2)
	}

	artistForm.EmailSent = sent
	artistFormsCollection := session.DB(config.DbName).C("artistForms")
	artistFormsCollection.Insert(artistForm)

	http.Redirect(w, req, "/thanks", 301)
}

func ThanksHandler(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("static/release_landing_page.html")
	t.Execute(w, new(interface{}))
}

func main() {
	portPtr := flag.String("port", "foo", "port number")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/work", FormHandler)
	router.HandleFunc("/thanks", ThanksHandler)
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
