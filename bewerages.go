package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//Menu :
type Menu struct {
	Title     string
	Bewerages []Drink
}

//Drink : hot bewerages available at the caffee
type Drink struct {
	Display  string
	Name     string  `json:",omitempty"`
	Price    float64 `json:",omitempty"`
	SubItems []Drink
	Question string `json:",omitempty"`
}

//DrinkIterFunc :
type DrinkIterFunc func(item Drink) bool

func traverse(input []Drink, callback DrinkIterFunc) *Drink {
	for _, item := range input {
		if callback(item) {
			return &item
		}
		r := traverse(item.SubItems, callback)
		if r != nil {
			return r
		}
	}
	return nil
}

func allNamedDrinks(bewerages []Drink) []Drink {
	output := make([]Drink, 0, len(bewerages))
	f := func(item Drink) bool {
		if item.Name != "" {
			output = append(output, item)
		}
		return false
	}
	traverse(bewerages, f)
	return output
}

func allDrinks(bewerages []Drink) []Drink {
	output := make([]Drink, 0, len(bewerages))
	f := func(item Drink) bool {
		output = append(output, item)
		return false
	}
	traverse(bewerages, f)
	return output
}

func numberOfDrinks(drinks []Drink) int {
	num := len(drinks)
	for _, item := range drinks {
		num += numberOfDrinks(item.SubItems)
	}
	return num
}

func loadBewerages(filePath string) (*Menu, error) {
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
