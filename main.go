package main

import "fmt"

func main() {
	nico := map[string]string{"name": "nico", "age": "12"}

	for key := range nico {
		fmt.Println(key)
	}

	// fmt.Println(nico)
}
