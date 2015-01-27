package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Artist struct {
	FirstName string
	LastName  string
	Email     string
	Link      string
	Works     []Work
}

type Work struct {
	Name        string
	Description string
	Photos      []Photo
}

type Photo struct {
	Name   string
	Models []Model
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

	fmt.Println(extractWorks(form))
	fmt.Println(artist.FirstName)
}

func getItemCount(filter string, form map[string][]string) int {
	count := 0
	for k, _ := range form {
		if strings.Contains(k, filter) {
			count++
		}
	}
	return count
}

func getIndices(filter string, form map[string][]string) []int {
	numItems := getItemCount(filter, form)
	indices := make([]int, numItems)
	i := 0
	for k, _ := range form {
		if strings.Contains(k, filter) {
			j, err := strconv.Atoi(k[len(k)-1:])
			if err != nil {
				panic(err)
			}

			indices[i] = j
			i++
		}
	}
	fmt.Println(indices)
	return indices
}

func extractWorks(form map[string][]string) []Work {
	numItems := getItemCount("descOfWork", form)
	workIndices := getIndices("descOfWork", form)
	works := make([]Work, numItems)
	j := 0
	for i := range workIndices {
		works[j] = Work{
			Name:        form[fmt.Sprintf("nameOfWork%d", i)][0],
			Description: form[fmt.Sprintf("descOfWork%d", i)][0],
		}
		j++
	}
	return works
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
