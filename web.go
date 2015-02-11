package main

import (
	"encoding/json"
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
		sendAllEmails(artistForm)
		writeArtistFormToDb(config.MongoUrl, sent, artistForm)
	}()

	http.Redirect(w, req, "/thanks", 301)
}

func ModelFormHandler(w http.ResponseWriter, req *http.Request) {
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

	selector := bson.M{"_id": bson.ObjectIdHex(form["artistId"][0])}
	update := bson.M{"$set": bson.M{
		"works.$.photos.$.models.firstname": form["firstName"][0],
		"works.$.photos.$.models.lastname":  form["lastName"][0],
		"works.$.photos.$.models.email":     form["emailAddress"][0],
		"works.$.photos.$.models.sig":       sig,
	}}

	err2, artistForms := makeOrGetCollection("artistForms")

	if err2 != nil {
		panic(err2)
	}

	err3 := artistForms.Update(selector, update)

	if err3 != nil {
		panic(err3)
	}
}

func ModelLandingHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := bson.ObjectIdHex(vars["id"])

	err1, artistForm := getArtistFromCollection(id)
	if err1 != nil {
		panic(err1)
	}

	t, _ := template.ParseFiles("static/model_release.html")
	t.Execute(w, makeContent(id, artistForm))
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
