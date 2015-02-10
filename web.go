package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
	"strings"
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

	go func() {
		makeAPDF(artistForm)
		err3 := sendEmail("PERJUS Magazine release form",
			config.ArtistEmailBody,
			true,
			artistForm.Form)

		//THERE MUST BE A BETTER WAY!!!!!!
		for i := 0; i < len(artistForm.Works); i++ {
			for j := 0; j < len(artistForm.Works[i].Photos); j++ {
				for k := 0; k < len(artistForm.Works[i].Photos[j].Models); k++ {
					//At least this will mitigate slowdown due to O(n^3) complexity somewhat! :-P
					go func(iIdx int, jIdx int, kIdx int) {
						modelErr := sendEmail("PERJUS Magazine model release form",
							fmt.Sprintf(config.ModelEmailBodyOne,
								strings.ToUpper(artistForm.FullName()),
								"http://www.google.com"),
							false,
							artistForm.Works[iIdx].Photos[jIdx].Models[kIdx].Form)

						snt := true
						if modelErr != nil {
							panic(modelErr)
							snt = false
						}

						artistForm.Works[iIdx].Photos[jIdx].Models[kIdx].EmailSent = snt
					}(i, j, k)
				}
			}
		}

		sent := true
		if err3 != nil {
			panic(err3)
			sent = false
		}
		writeArtistFormToDb(config.MongoUrl, sent, artistForm)
	}()

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
