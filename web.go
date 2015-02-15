package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"os"
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
		err2, sent := sendArtistEmail(artistForm)
		if err2 != nil {
			panic(err2)
		}
		sendAdminEmail(artistForm)
		sendAllModelEmails(artistForm)

		writeArtistFormToDb(config.MongoUrl, sent, artistForm)
	}()

	http.Redirect(w, req, "/thanks", 301)
}

func ModelFormHandler(w http.ResponseWriter, req *http.Request) {
	err0 := req.ParseForm()
	if err0 != nil {
		//panic(err0)
		errorHandler(w, req, err0)
	}

	form := req.PostForm

	artistId := bson.ObjectIdHex(form["artistId"][0])
	modelId := bson.ObjectIdHex(form["modelId"][0])

	err1, artistForms := makeOrGetCollection("artistForms")

	if err1 != nil {
		//panic(err1)
		errorHandler(w, req, err1)
	}

	go func() {
		var artist ArtistForm
		artistForms.FindId(artistId).One(&artist)
		artist.SetModelSigById(modelId, form["output"][0])
		artistForms.UpdateId(artistId, artist)
		model := artist.ModelById(modelId)

		//TODO Store Admin Sent status
		err2, adminSent := sendAdminEmail(model)
		if err2 != nil {
			fmt.Println(adminSent)
			//panic(err2)
			errorHandler(w, req, err2)
		}

		err3, sent := sendModelEmailWithForm(model)
		if err3 != nil {
			//panic(err3)
			errorHandler(w, req, err3)
		}
		artist.SetModelSentById(modelId, sent)
	}()

	http.Redirect(w, req, "/thanks", 301)
}

func ModelLandingHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := bson.ObjectIdHex(vars["id"])

	err, artistForm := getArtistFromCollection(id)
	if err != nil {
		errorHandler(w, req, err)
		//panic(err1)
	}

	t, _ := template.ParseFiles("static/model_release.html")
	t.Execute(w, makeContent(id, artistForm))
}

func ThanksHandler(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("static/release_landing_page.html")
	t.Execute(w, new(interface{}))
}

func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(500)
	errString := fmt.Sprint("Something broke.. :-/ womp womp:\n %s ")
	fmt.Fprint(w, errString)
}

func main() {
	portPtr := flag.String("port", "foo", "port number")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/work", WorkFormHandler)
	router.HandleFunc("/model", ModelFormHandler)
	router.HandleFunc("/thanks", ThanksHandler)
	router.HandleFunc("/models/{id}/release", ModelLandingHandler)

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
