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

const jobCsv string = "jobs.csv"

func handleScrape(c echo.Context) error {
	defer os.Remove(jobCsv)

	searchTerm := strings.ToLower(scrapper.CleanString(c.FormValue("searchTerm")))
	scrapper.ScrapeShit(searchTerm)

	return c.Attachment(jobCsv, jobCsv)
}

func main() {

	scrapper.ScrapeShit("flutter")

	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}
