package mydict

import (
	"errors"
)

// Dictionary type
type Dictionary map[string]string

var errNotFound = errors.New("Not found")
var errWordExists = errors.New("Word already in Dictionary")
var errCantUpdate = errors.New("Can't update a work that doesn't exist")

// Search for a word in the Dictionary
func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word]
	if exists {
		return value, nil
	}
	return "", errNotFound
}

// Add word to Dictionary
func (d Dictionary) Add(word string, def string) error {
	_, err := d.Search(word)
	if err == errNotFound {
		d[word] = def
	} else if err == nil {
		return errWordExists
	}
	return nil
}

// Update existing word
func (d Dictionary) Update(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = def
	case errNotFound:
		return errCantUpdate
	}
	return nil

}

// Delete word in Dictionary
func (d Dictionary) Delete(word string) {
	delete(d, word)
}
