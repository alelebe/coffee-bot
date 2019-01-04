package main

import (
	"strings"
	"testing"
)

const filePath string = "./data/benugo.json"

func TestBewerages(t *testing.T) {
	benugo, err := loadBewerages(filePath)
	if err != nil {
		t.Errorf("'%s' can't be loaded: %v", filePath, err)
		return
	}
	if false == strings.HasPrefix(benugo.Title, "Benugo") {
		t.Errorf("'%s' has wrong 'Title'", filePath)
	}
	if 6 != len(benugo.Bewerages) {
		t.Errorf("'%s' has wrong number of items: %d != 6", filePath, len(benugo.Bewerages))
	}

	total := numberOfDrinks(benugo.Bewerages)
	if len(benugo.Bewerages) == total {
		t.Errorf("'%s' has wrong number of items in total: %d == %d", filePath,
			total, len(benugo.Bewerages))
	}

	all := allDrinks(benugo.Bewerages)
	if total != len(all) {
		t.Errorf("'%s' has wrong number of all drinks: %d == %d", filePath,
			total, len(all))
	}
	if len(benugo.Bewerages) == len(all) {
		t.Errorf("'%s' has wrong number of all drinks, no groups detected: %d == %d", filePath,
			len(benugo.Bewerages), len(all))
	}
}
