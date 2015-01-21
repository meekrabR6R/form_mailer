package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.Handle("/", http.FileServer(http.Dir("./static/")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
