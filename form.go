package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Indexable interface {
	SetId()
}

type BaseForm interface {
	FullName() string
	SetSignature(string) error
	IsArtist() bool
	IsModel() bool
}

type Form struct {
	Id        string `bson:"_id"`
	FirstName string
	LastName  string
	Email     string
	Link      string
	Sig       []map[string]int
	EmailSent bool
}

func (f *Form) SetId() {
	f.Id = randomHex()
}

func (f *Form) FullName() string {
	return fmt.Sprintf("%s %s", f.FirstName, f.LastName)
}

func (f *Form) SetSignature(sigString string) error {
	rawSig := []byte(sigString)
	var sig []map[string]int
	err := json.Unmarshal(rawSig, &sig)
	f.Sig = sig
	return err
}

type ArtistForm struct {
	Form  `bson:",inline"`
	Works []Work
}

func (a *ArtistForm) IsArtist() bool {
	return true
}

func (a *ArtistForm) IsModel() bool {
	return false
}

func (a *ArtistForm) SetWorks(form map[string][]string) {
	numItems := getItemCount("descOfWork", form)
	workIndices := getIndices("descOfWork", form)
	a.Works = make([]Work, numItems)

	for i, e := range workIndices {
		a.Works[i] = Work{
			Name:        form[fmt.Sprintf("nameOfWork%d", e)][0],
			Description: form[fmt.Sprintf("descOfWork%d", e)][0],
		}

		a.Works[i].SetPhotos(form, e)
		a.Works[i].SetId()
	}
}

type ModelForm struct {
	Form `bson:",inline"`
}

func (m *ModelForm) IsArtist() bool {
	return false
}

func (m *ModelForm) IsModel() bool {
	return true
}

type Work struct {
	Id          string `bson:"_id"`
	Name        string
	Description string
	Photos      []Photo
}

func (w *Work) SetId() {
	w.Id = randomHex()
}

func (w *Work) SetPhotos(form map[string][]string, workIndex int) {
	filter := fmt.Sprintf("nameOfPhoto%d", workIndex)
	numItems := getItemCount(filter, form)
	photoIndices := getIndices(filter, form)
	w.Photos = make([]Photo, numItems)

	for i, e := range photoIndices {
		w.Photos[i] = Photo{
			Name: form[fmt.Sprintf("nameOfPhoto%d%d", workIndex, e)][0],
		}

		w.Photos[i].SetModels(form, workIndex, e)
		w.Photos[i].SetId()
	}
}

type Photo struct {
	Id     string `bson:"_id"`
	Name   string
	Models []ModelForm
}

func (p *Photo) SetId() {
	p.Id = randomHex()
}

func (p *Photo) SetModels(form map[string][]string, workIndex int,
	photoIndex int) {

	filter := fmt.Sprintf("firstNameOfModel%d%d", workIndex, photoIndex)
	numItems := getItemCount(filter, form)
	modelIndices := getIndices(filter, form)
	p.Models = make([]ModelForm, numItems)

	for i, e := range modelIndices {
		p.Models[i] = ModelForm{
			Form: Form{
				FirstName: form[fmt.Sprintf("firstNameOfModel%d%d%d",
					workIndex, photoIndex, e)][0],
				LastName: form[fmt.Sprintf("lastNameOfModel%d%d%d",
					workIndex, photoIndex, e)][0],
				Email: form[fmt.Sprintf("emailOfModel%d%d%d",
					workIndex, photoIndex, e)][0],
			},
		}

		p.Models[i].SetId()
	}
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

func randomHex() string {
	var numbers = []rune("abcdef0123456789")
	b := make([]rune, 32)
	for i := range b {
		b[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(b)
}
