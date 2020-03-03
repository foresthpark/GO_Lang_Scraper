package scrapper

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

var errRequesFailed = errors.New("request failed")

//ScrapeShit on indeed.ca for Calgary,AB
func ScrapeShit(searchTerm string) {
	var baseURL string = "https://www.indeed.ca/jobs?q=" + searchTerm + "&l=Calgary%2C+AB"

	var jobs []extractedJob
	c := make(chan []extractedJob)

	totalPages := getPages(baseURL)

	for i := 0; i < totalPages; i++ {
		go getPage(i, baseURL, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)

	}

	defer fmt.Println("All requests are successfully done")
	writeJobs(jobs)
}

func getPage(page int, url string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)

	pageURL := url + "&start=" + strconv.Itoa(page*20)
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

func getPages(url string) int {
	pages := 0

	res, err := http.Get(url)

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
	title := CleanString(jobcard.Find(".title>a").Text())
	company := CleanString(jobcard.Find(".company").Text())
	salary := CleanString(jobcard.Find(".salaryText").Text())
	summary := CleanString(jobcard.Find(".summary").Text())

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

//CleanString cleans a String
func CleanString(str string) string {
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
