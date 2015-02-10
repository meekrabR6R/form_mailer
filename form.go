package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Record interface {
	SetId()
	SetCreatedAt()
	SetUpdatedAt()
}

type BaseForm interface {
	FullName() string
	SetSignature(string) error
	IsArtist() bool
	IsModel() bool
}

type Form struct {
	Id        bson.ObjectId `bson:"_id"`
	FirstName string
	LastName  string
	Email     string
	Link      string
	Sig       []map[string]int
	EmailSent bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (f *Form) SetId() {
	f.Id = bson.NewObjectId()
}

func (f *Form) SetCreatedAt() {
	f.CreatedAt = time.Now()
}

func (f *Form) SetUpdatedAt() {
	f.UpdatedAt = time.Now()
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
	Works []Work `bson:"works"`
}

func (a *ArtistForm) IsArtist() bool {
	return true
}

func (a *ArtistForm) IsModel() bool {
	return false
}

func (a *ArtistForm) ModelById(id bson.ObjectId) ModelForm {
	var m ModelForm

	for i := 0; i < len(a.Works); i++ {
		for j := 0; j < len(a.Works[i].Photos); j++ {
			for k := 0; k < len(a.Works[i].Photos[j].Models); k++ {
				if a.Works[i].Photos[j].Models[k].Id == id {
					m = a.Works[i].Photos[j].Models[k]
					break
				}
			}
		}
	}

	return m
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
		writeNewMetaData(&a.Works[i])
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
	Id          bson.ObjectId `bson:"_id"`
	Name        string
	Description string
	Photos      []Photo `bson:"photos"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (w *Work) SetId() {
	w.Id = bson.NewObjectId()
}

func (w *Work) SetCreatedAt() {
	w.CreatedAt = time.Now()
}

func (w *Work) SetUpdatedAt() {
	w.UpdatedAt = time.Now()
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
		writeNewMetaData(&w.Photos[i])
	}
}

type Photo struct {
	Id        bson.ObjectId `bson:"_id"`
	Name      string
	Models    []ModelForm `bson:"models"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Photo) SetId() {
	p.Id = bson.NewObjectId()
}

func (p *Photo) SetCreatedAt() {
	p.CreatedAt = time.Now()
}

func (p *Photo) SetUpdatedAt() {
	p.UpdatedAt = time.Now()
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

		writeNewMetaData(&p.Models[i].Form)
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

func writeNewMetaData(record Record) {
	record.SetId()
	record.SetCreatedAt()
	record.SetUpdatedAt()
}