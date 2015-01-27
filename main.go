package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type Artist struct {
	FirstName string
	LastName  string
	Email     string
	Link      string
}

type Work struct {
	Name        string
	Description string
}

type Photo struct {
	Name string
}

type Model struct {
	Name  string
	Email string
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

	fmt.Println(artist.FirstName)
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
