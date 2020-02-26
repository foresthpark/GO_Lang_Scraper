package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type urlResult struct {
	url    string
	status string
}

type extractedJob struct {
	id      string
	title   string
	salary  string
	company string
	summary string
}

var baseURL string = "https://www.indeed.ca/jobs?q=full+stack+developer&l=Calgary%2C+AB"

var errRequesFailed = errors.New("request failed")

func main() {
	var jobs []extractedJob
	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}

	fmt.Println(jobs)
}

func getPage(page int) []extractedJob {
	var jobs []extractedJob
	pageURL := baseURL + "&start=" + strconv.Itoa(page*20)
	fmt.Println("Requestiong pageURL: " + pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, jobcard *goquery.Selection) {
		job := extractJobInfo(jobcard)
		jobs = append(jobs, job)
	})

	return jobs

}

func getPages() int {
	pages := 0

	res, err := http.Get(baseURL)

	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func extractJobInfo(jobcard *goquery.Selection) extractedJob {
	id, _ := jobcard.Attr("data-jk")
	title := cleanString(jobcard.Find(".title>a").Text())
	company := jobcard.Find(".company").Text()
	salary := cleanString(jobcard.Find(".salaryText").Text())
	summary := cleanString(jobcard.Find(".summary").Text())

	return extractedJob{
		id:      id,
		title:   title,
		salary:  salary,
		company: company,
		summary: summary,
	}

	// fmt.Println(id, title, company, salary, summary)
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with status: ", res.StatusCode)
	}

}
