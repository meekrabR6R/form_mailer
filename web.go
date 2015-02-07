package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
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
