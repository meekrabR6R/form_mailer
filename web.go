package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func WorkFormHandler(w http.ResponseWriter, req *http.Request) {
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
						url := fmt.Sprintf("%s/models/%s/release", config.Url,
							artistForm.Works[iIdx].Photos[jIdx].Models[kIdx].Id.Hex())
						fmt.Println(url)
						modelErr := sendEmail("PERJUS Magazine model release form",
							fmt.Sprintf(config.ModelEmailBodyOne,
								strings.ToUpper(artistForm.FullName()),
								url),
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

func ModelFormHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := bson.ObjectIdHex(vars["id"])
	err, artistForms := makeOrGetCollection("artistForms")

	if err != nil {
		panic(err)
	}

	artistForm := ArtistForm{}
	query := bson.M{"works.photos.models._id": id}
	err1 := artistForms.Find(query).One(&artistForm)

	if err1 != nil {
		panic(err1)
	}

	model := artistForm.ModelById(id)

	fmt.Println(id)
	fmt.Println(model)
}

func ThanksHandler(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("static/release_landing_page.html")
	t.Execute(w, new(interface{}))
}

func main() {
	portPtr := flag.String("port", "foo", "port number")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/work", WorkFormHandler)
	router.HandleFunc("/thanks", ThanksHandler)
	router.HandleFunc("/models/{id}/release", ModelFormHandler)

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
