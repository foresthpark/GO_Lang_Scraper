package main

import (
	"os"
	"strings"

	"github.com/foresthpark/learngo/scrapper"
	"github.com/labstack/echo"
)

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

// Name of saved CSV file
const jobCsv string = "jobs.csv"

// Change the Search Term here to search for different jobs in the Calgary, AB Area
const searchTerm string = "flutter" // For example, this will search for Flutter

func handleScrape(c echo.Context) error {
	defer os.Remove(jobCsv)

	searchTerm := strings.ToLower(scrapper.CleanString(c.FormValue("searchTerm")))
	scrapper.ScrapeShit(searchTerm)

	return c.Attachment(jobCsv, jobCsv)
}

func main() {

	scrapper.ScrapeShit(searchTerm)

	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}
