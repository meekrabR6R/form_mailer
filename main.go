package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	//router.Handle("/", http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", router)

	fmt.Println("Listening on 8000")
	http.ListenAndServe(":8000", nil)
}
