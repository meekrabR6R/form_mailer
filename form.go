package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"sort"
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
	FullNameForFile() string
	FullAddress() string
	SetSignature(string) error
	GetSignature() []map[string]float64
	GetDataAsString() string
	IsArtist() bool
	IsModel() bool
}

type Form struct {
	Id         bson.ObjectId        `bson:"_id" json:"id"` //,`json:"id"`
	FirstName  string               `json:"first_name"`
	LastName   string               `json:"last_name"`
	AddressOne string               `json:"address_one"`
	AddressTwo string               `json:"address_two"`
	City       string               `json:"city"`
	State      string               `json:"state"`
	Zip        string               `json:"zip"`
	Country    string               `json:"country"`
	Email      string               `json:"email"`
	Link       string               `json:"link"`
	Sig        []map[string]float64 `json:"sig"`
	EmailSent  bool                 `json:"email_sent"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
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

func (f *Form) FullNameForFile() string {
	return fmt.Sprintf("%s_%s", f.FirstName, f.LastName)
}

func (f *Form) FullAddress() string {
	return fmt.Sprintf("%s %s, %s %s %s %s",
		f.AddressOne, f.AddressTwo, f.City, f.State, f.Zip, f.Country)
}

func (f *Form) SetSignature(sigString string) error {
	rawSig := []byte(sigString)
	var sig []map[string]float64
	err := json.Unmarshal(rawSig, &sig)
	f.Sig = sig
	return err
}

func (f *Form) GetSignature() []map[string]float64 {
	return f.Sig
}

type ArtistForm struct {
	Form  `bson:",inline"`
	Works []Work `bson:"works" json:"works"`
}

func (a *ArtistForm) IsArtist() bool {
	return true
}

func (a *ArtistForm) IsModel() bool {
	return false
}

func (a *ArtistForm) GetDataAsString() string {
	var worksBuffer bytes.Buffer

	index := 0
	for _, work := range a.Works {
		if index < (len(a.Works)-1) || len(a.Works) == 1 {
			worksBuffer.WriteString("[" + work.Name + "] ([" + photosAsString(work.Photos) + "]) (\"Images\"), ")
		} else {
			worksBuffer.WriteString("and [" + work.Name + "] ([" + photosAsString(work.Photos) + "]) (\"Images\"),")
		}
		index++
	}
	return worksBuffer.String()
}

func (a *ArtistForm) SetModelSigById(id bson.ObjectId, sig string) {
	for i := 0; i < len(a.Works); i++ {
		for j := 0; j < len(a.Works[i].Photos); j++ {
			for k := 0; k < len(a.Works[i].Photos[j].Models); k++ {
				if a.Works[i].Photos[j].Models[k].Id == id {
					a.Works[i].Photos[j].Models[k].SetSignature(sig)
					break
				}
			}
		}
	}
}

func (a *ArtistForm) SetModelSentById(id bson.ObjectId, sent bool) {
	for i := 0; i < len(a.Works); i++ {
		for j := 0; j < len(a.Works[i].Photos); j++ {
			for k := 0; k < len(a.Works[i].Photos[j].Models); k++ {
				if a.Works[i].Photos[j].Models[k].Id == id {
					a.Works[i].Photos[j].Models[k].EmailSent = sent
					break
				}
			}
		}
	}
}

func (a *ArtistForm) UpdateModelById(id bson.ObjectId, form map[string][]string) {
	for i := 0; i < len(a.Works); i++ {
		for j := 0; j < len(a.Works[i].Photos); j++ {
			for k := 0; k < len(a.Works[i].Photos[j].Models); k++ {
				if a.Works[i].Photos[j].Models[k].Id == id {
					a.Works[i].Photos[j].Models[k].FirstName = form["firstName"][0]
					a.Works[i].Photos[j].Models[k].LastName = form["lastName"][0]
					a.Works[i].Photos[j].Models[k].Email = form["emailAddress"][0]
					a.Works[i].Photos[j].Models[k].AddressOne = form["addressOne"][0]
					a.Works[i].Photos[j].Models[k].AddressTwo = form["addressTwo"][0]
					a.Works[i].Photos[j].Models[k].City = form["city"][0]
					a.Works[i].Photos[j].Models[k].State = form["state"][0]
					a.Works[i].Photos[j].Models[k].Zip = form["zip"][0]
					a.Works[i].Photos[j].Models[k].Country = form["country"][0]
					break
				}
			}
		}
	}
}

func (a *ArtistForm) ModelById(id bson.ObjectId) *ModelForm {
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

	return &m
}

func (a *ArtistForm) WorkByContentId(id string) Work {
	var w Work
	for i := 0; i < len(a.Works); i++ {
		if a.Works[i].ContentId == id {
			w = a.Works[i]
			break
		}
	}
	return w
}

/**
 * This (and the other array setters) smell a bit. They should be
 * agnostic w/r/t the original form structure. DECOUPLE THIS STUFF ASAP!
 */
func (a *ArtistForm) SetWorks(form map[string][]string) {
	numItems := getItemCount("descOfWork", form)
	workIndices := getIndices("descOfWork", form)
	a.Works = make([]Work, numItems)

	for i, e := range workIndices {
		a.Works[i] = Work{
			ContentId:   bson.NewObjectId().Hex(),
			Name:        form[fmt.Sprintf("nameOfWork%d", e)][0],
			Description: form[fmt.Sprintf("descOfWork%d", e)][0],
			Extra:       form[fmt.Sprintf("extraForWork%d", e)][0],
		}
		a.Works[i].SetPhotos(form, e)
		writeNewMetaData(&a.Works[i])
	}
}

type ModelForm struct {
	Form   `bson:",inline"`
	WorkId string `json:"work_id"`
}

func (m *ModelForm) IsArtist() bool {
	return false
}

func (m *ModelForm) IsModel() bool {
	return true
}

func (m *ModelForm) GetWork() (error, Work) {
	var work Work
	err, artistForms := makeOrGetCollection("artistForms")
	artistForm := ArtistForm{}
	query := bson.M{"works.contentid": m.WorkId}
	if err == nil {
		artistForms.Find(query).One(&artistForm)
		work = artistForm.WorkByContentId(m.WorkId)
	}
	return err, work
}

func (m *ModelForm) GetDataAsString() string {
	_, work := m.GetWork()
	return "[" + work.Name + "] ([" + photosAsString(work.Photos) + "]) (\"Images\"), "
}

type Work struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	ContentId   string        `json:"content_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Extra       string        `json:"extra"`
	Photos      []Photo       `bson:"photos" json:"photos"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
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

/**
 * This (and the other array setters) smell a bit. They should be
 * agnostic w/r/t the original form structure. DECOUPLE THIS STUFF ASAP!
 */
func (w *Work) SetPhotos(form map[string][]string, workIndex int) {
	filter := fmt.Sprintf("nameOfPhoto%d", workIndex)
	numItems := getItemCount(filter, form)
	photoIndices := getIndices(filter, form)
	w.Photos = make([]Photo, numItems)

	for i, e := range photoIndices {
		w.Photos[i] = Photo{
			Name:   form[fmt.Sprintf("nameOfPhoto%d%d", workIndex, e)][0],
			Title:  form[fmt.Sprintf("titleOfPhoto%d%d", workIndex, e)][0],
			WorkId: w.ContentId,
		}

		w.Photos[i].SetModels(form, workIndex, e)
		writeNewMetaData(&w.Photos[i])
	}
}

type Photo struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `json:"name"`
	Title     string        `json:"title"`
	WorkId    string        `json:"work_id"`
	Models    []ModelForm   `bson:"models" json:"models"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
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

/**
 * This (and the other array setters) smell a bit. They should be
 * agnostic w/r/t the original form structure. DECOUPLE THIS STUFF ASAP!
 */
func (p *Photo) SetModels(form map[string][]string, workIndex int,
	photoIndex int) {

	filter := fmt.Sprintf("firstNameOfModel%d%d-", workIndex, photoIndex)
	numItems := getItemCount(filter, form)
	modelIndices := getIndices(filter, form)
	p.Models = make([]ModelForm, numItems)

	for i, e := range modelIndices {
		p.Models[i] = ModelForm{
			Form: Form{
				FirstName: form[fmt.Sprintf("firstNameOfModel%d%d-%d",
					workIndex, photoIndex, e)][0],
				LastName: form[fmt.Sprintf("lastNameOfModel%d%d-%d",
					workIndex, photoIndex, e)][0],
				Email: form[fmt.Sprintf("emailOfModel%d%d-%d",
					workIndex, photoIndex, e)][0],
			},
		}

		p.Models[i].WorkId = p.WorkId

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
				//panic(err)
			}

			indices[i] = j
			i++
		}
	}
	sort.Ints(indices)
	return indices
}

func makeArtistForm(form map[string][]string) (error, *ArtistForm) {
	artistForm := &ArtistForm{
		Form: Form{
			FirstName:  form["firstName"][0],
			LastName:   form["lastName"][0],
			AddressOne: form["addressOne"][0],
			AddressTwo: form["addressTwo"][0],
			City:       form["city"][0],
			State:      form["state"][0],
			Zip:        form["zip"][0],
			Country:    form["country"][0],
			Email:      form["emailAddress"][0],
			Link:       form["downloadLink"][0],
		},
	}

	artistForm.SetWorks(form)
	writeNewMetaData(artistForm)

	err := artistForm.SetSignature(form["output"][0])
	return err, artistForm
}

//Probably will never be used..
func randomHex() string {
	var numbers = []rune("abcdef0123456789")
	b := make([]rune, 32)
	for i := range b {
		b[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(b)
}

func photosAsString(photos []Photo) string {
	var worksBuffer bytes.Buffer
	for i, photo := range photos {
		if i < len(photos)-1 {
			worksBuffer.WriteString(photo.Name + ", ")
		} else {
			worksBuffer.WriteString(photo.Name)
		}
	}
	return worksBuffer.String()
}

func writeNewMetaData(record Record) {
	record.SetId()
	record.SetCreatedAt()
	record.SetUpdatedAt()
}
