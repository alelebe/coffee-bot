package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//Menu :
type Menu struct {
	Title string
	Entry Beverages
}

// Beverages :
type Beverages struct {
	Question string
	Items    []Drink
}

//Drink : hot beverages available at the caffee
type Drink struct {
	ID      string
	Display string
	Price   float64 `json:",omitempty"`
	Entry   Beverages
}

//DrinkIterFunc :
type DrinkIterFunc func(item Drink) bool

func traverse(items []Drink, callback DrinkIterFunc) *Drink {
	for _, it := range items {
		if callback(it) {
			return &it
		}
		r := traverse(it.Entry.Items, callback)
		if r != nil {
			return r
		}
	}
	return nil
}

func (entry Beverages) getAllEntries() []Drink {
	output := make([]Drink, 0, len(entry.Items))
	f := func(item Drink) bool {
		output = append(output, item)
		return false
	}
	traverse(entry.Items, f)
	return output
}

func (entry Beverages) getDrinkByID(ID string) *Drink {
	var item *Drink
	f := func(it Drink) bool {
		if it.ID == ID {
			item = &it
			return true
		}
		return false
	}
	traverse(entry.Items, f)
	return item
}

func loadBeverages(filePath string) (*Menu, error) {
	var err error

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	menu := new(Menu)
	err = json.Unmarshal(data, &menu)
	if err != nil {
		return nil, err
	}
	return menu, nil
}
