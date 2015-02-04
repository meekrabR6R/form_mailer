package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Form interface {
	FullName() string
	SetSignature(string) error
}

type ArtistForm struct {
	FirstName string
	LastName  string
	Email     string
	Link      string
	Sig       []map[string]int
	EmailSent bool
	Works     []Work
}

func (a *ArtistForm) FullName() string {
	return fmt.Sprintf("%s %s", a.FirstName, a.LastName)
}

func (a *ArtistForm) SetSignature(sigString string) error {
	rawSig := []byte(sigString)
	var sig []map[string]int
	err := json.Unmarshal(rawSig, &sig)
	a.Sig = sig
	return err
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

func getItemCount(filter string, form map[string][]string) int {
	count := 0
	for k, _ := range form {
		if strings.Contains(k, filter) {
			count++
		}
	}
	return count
}

/**
 * Returns an array of indices for works. This is needed b/c
 * it is possible for user to remove works.
 */
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
	return indices
}

func extractWorks(form map[string][]string) []Work {
	numItems := getItemCount("descOfWork", form)
	workIndices := getIndices("descOfWork", form)
	works := make([]Work, numItems)

	for i, e := range workIndices {
		works[i] = Work{
			Name:        form[fmt.Sprintf("nameOfWork%d", e)][0],
			Description: form[fmt.Sprintf("descOfWork%d", e)][0],
		}
	}
	return works
}
