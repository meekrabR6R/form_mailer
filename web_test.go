package main

import (
	"testing"
)

var item1 = map[string][]string{
	"name":        []string{"Nick"},
	"email":       []string{"freddy@aol.com"},
	"descOfWork1": []string{"Blob"},
	"descOfWork3": []string{"Slob"},
	"descOfWork5": []string{"Corn on the cob"},
	"nameOfWork1": []string{"Pic of Blob"},
	"nameOfWork3": []string{"Pic of Slob"},
	"nameOfWork5": []string{"Pic of Corn on the cob"},
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
		t.Error("For", item1,
			"expected", 3,
			"got", len(workIndices))
	}

	if workIndices[0] != 1 {
		t.Error("For", item1,
			"expected", 1,
			"got", workIndices[0])
	}

	if workIndices[1] != 3 {
		t.Error("For", item1,
			"expected", 3,
			"got", workIndices[1])
	}

	if workIndices[2] != 5 {
		t.Error("For", item1,
			"expected", 5,
			"got", workIndices[2])
	}
}

/**
 * Test for private function extraWorks
 */
func TestExtractWorks(t *testing.T) {
	var works = extractWorks(item1)

	if len(works) != 3 {
		t.Error("For", item1,
			"expected", 3,
			"got", len(works))
	}
}
