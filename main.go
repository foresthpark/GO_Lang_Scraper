package main

import (
	"errors"
	"fmt"
	"net/http"
)

type urlResult struct {
	url    string
	status string
}

var errRequesFailed = errors.New("request failed")

func main() {

	results := make(map[string]string)
	c := make(chan urlResult)

	urls := []string{
		"https://www.amazon.ca/",
		"https://www.google.com/",
		"https://www.airbnb.ca/",
		"https://reddit.com",
		"https://facebook.com",
		"https://www.youtube.com/",
		"https://www.linkedin.com/",
	}

	for _, url := range urls {
		go checkUrl(url, c)
	}

	for i := 0; i < len(urls); i++ {
		result := <-c
		results[result.url] = result.status
	}

	for url, status := range results {
		fmt.Println(url, status)
	}

}

func checkUrl(url string, c chan<- urlResult) {
	fmt.Println("Checking: " + url)
	resp, err := http.Get(url)
	status := "OK"
	if err != nil || resp.StatusCode >= 400 {
		status = "FAILED"
	}
	c <- urlResult{url: url, status: status}
}
