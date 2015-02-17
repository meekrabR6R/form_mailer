package main

import (
	"fmt"
	"testing"
)

var item1 = map[string][]string{
	"firstName":           []string{"Nick"},
	"lastName":            []string{"Mujanjo"},
	"emailAddress":        []string{"freddy@aol.com"},
	"downloadLink":        []string{"http://www.dropbox.com/myawesomestuff"},
	"output":              []string{"[{ \"lx\" : 122, \"ly\" : 31, \"mx\" : 122, \"my\" : 30 }]"},
	"descOfWork1":         []string{"Blob"},
	"descOfWork3":         []string{"Slob"},
	"descOfWork5":         []string{"Corn on the cob"},
	"nameOfWork1":         []string{"Pic of Blob"},
	"nameOfWork3":         []string{"Pic of Slob"},
	"nameOfWork5":         []string{"Pic of Corn on the cob"},
	"nameOfPhoto10":       []string{"Photo of Blob"},
	"firstNameOfModel100": []string{"Troy"},
	"lastNameOfModel100":  []string{"McClure"},
	"emailOfModel100":     []string{"troy@mcclure.com"},
}

/**
 * Test for private function getItemCount
 */
func TestGetItemCount(t *testing.T) {
	var workCountOfThree = getItemCount("descOfWork", item1)
	var workCountOfZero = getItemCount("fart", item1)

	if workCountOfThree != 3 {
		t.Error("For", item1,
			"expected", 3,
			"got", workCountOfThree)
	}

	if workCountOfZero != 0 {
		t.Error("For", item1,
			"expected", 0,
			"got", workCountOfZero)
	}

}

/**
 * Test for private function getIndices
 */
func TestGetIndices(t *testing.T) {
	var workIndices = getIndices("descOfWork", item1)

	if len(workIndices) != 3 {
		t.Error("For", workIndices,
			"expected", 3,
			"got", len(workIndices))
	}

	if workIndices[0] != 1 {
		t.Error("For", workIndices,
			"expected", 1,
			"got", workIndices[0])
	}

	if workIndices[1] != 3 {
		t.Error("For", workIndices,
			"expected", 3,
			"got", workIndices[1])
	}

	if workIndices[2] != 5 {
		t.Error("For", workIndices,
			"expected", 5,
			"got", workIndices[2])
	}
}

/**
 * Test ArtistForm factory helper
 */
func TestMakeArtistForm(t *testing.T) {
	err, artistForm := makeArtistForm(item1)

	fullName := fmt.Sprintf("%s %s", item1["firstName"][0],
		item1["lastName"][0])

	fullNameForFile := fmt.Sprintf("%s_%s", item1["firstName"][0],
		item1["lastName"][0])

	if err != nil {
		panic(err)
		return
	}

	if artistForm.Email != item1["emailAddress"][0] {
		t.Error("For", artistForm,
			"expected", item1["emailAddress"][0],
			"got", artistForm.Email)
	}

	if artistForm.FirstName != item1["firstName"][0] {
		t.Error("For", artistForm,
			"expected", item1["firstName"][0],
			"got", artistForm.Email)
	}

	if artistForm.LastName != item1["lastName"][0] {
		t.Error("For", artistForm,
			"expected", item1["emailAddress"][0],
			"got", artistForm.Email)
	}

	if artistForm.FullName() != fullName {
		t.Error("For", artistForm,
			"expected", fullName,
			"got", artistForm.FullName())
	}

	if artistForm.FullNameForFile() != fullNameForFile {
		t.Error("For", artistForm,
			"expected", fullNameForFile,
			"got", artistForm.FullNameForFile())
	}

	if artistForm.EmailSent {
		t.Error("For", artistForm,
			"expected", false,
			"got", artistForm.EmailSent)
	}

	if artistForm.Link != item1["downloadLink"][0] {
		t.Error("For", artistForm,
			"expected", item1["downloadLink"][0],
			"got", artistForm.Link)
	}

	if len(artistForm.Works) != 3 {
		t.Error("For", artistForm,
			"expected", 3,
			"got", len(artistForm.Works))
	}

	if artistForm.Works[0].Name != "Pic of Blob" {
		t.Error("For", artistForm.Works[0],
			"expected", "Pic of Blob",
			"got", artistForm.Works[0].Name)
	}

	if artistForm.Works[1].Name != "Pic of Slob" {
		t.Error("For", artistForm.Works[0],
			"expected", "Pic of Slob",
			"got", artistForm.Works[1].Name)
	}

	if artistForm.Works[2].Name != "Pic of Corn on the cob" {
		t.Error("For", artistForm.Works[0],
			"expected", "Pic of Corn on the cob",
			"got", artistForm.Works[2].Name)
	}

	if artistForm.IsModel() {
		t.Error("For", artistForm,
			"expected", false,
			"got", artistForm.IsModel())
	}

	if !artistForm.IsArtist() {
		t.Error("For", artistForm,
			"expected", true,
			"got", artistForm.IsArtist())
	}
}

/**
 * Test Work struct
 **/
