package main

import (
	"fmt"

	"github.com/foresthpark/learngo/mydict"
)

func main() {
	dictionary := mydict.Dictionary{"first": "first word"}
	// defintion, err := dictionary.Search("second")
	baseWord := "hello"

	dictionary.Add(baseWord, "First")
	dictionary.Search(baseWord)

	dictionary.Delete(baseWord)

	word, err := dictionary.Search(baseWord)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(word)
	}

}
