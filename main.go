package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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
	c := make(chan []extractedJob)

	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)

	}

	defer fmt.Println("All requests are successfully done")
	writeJobs(jobs)
}

func getPage(page int, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)

	pageURL := baseURL + "&start=" + strconv.Itoa(page*20)
	fmt.Println("Requestiong pageURL: " + pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, jobcard *goquery.Selection) {
		go extractJobInfo(jobcard, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs

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

func extractJobInfo(jobcard *goquery.Selection, c chan<- extractedJob) {
	id, _ := jobcard.Attr("data-jk")
	title := cleanString(jobcard.Find(".title>a").Text())
	company := cleanString(jobcard.Find(".company").Text())
	salary := cleanString(jobcard.Find(".salaryText").Text())
	summary := cleanString(jobcard.Find(".summary").Text())

	c <- extractedJob{
		id:      id,
		title:   title,
		salary:  salary,
		company: company,
		summary: summary,
	}

}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Salary", "Company", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://www.indeed.ca/viewjob?jk=" + job.id, job.title, job.salary, job.company, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}

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
