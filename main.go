package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/leekchan/timeutil"
)

func checkError(err error) {
	if err != nil {
		log.Fatalln("Error!:", err)
		panic(err)
	}
}

func writeContent(url string, writer *csv.Writer, page string) {
	response, err := http.Get(url)
	checkError(err)
	defer response.Body.Close()
	if statusCode := response.StatusCode; statusCode > 400 {
		log.Fatalln("Status code:", statusCode)
		return
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	checkError(err)

	doc.Find("td.title").Find("span.titleline").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		link, _ := s.Find("a").Attr("href")
		// fmt.Printf("idx: %v\ttext: %v\tlink: %v\n", i, text, link)
		writer.Write([]string{fmt.Sprintf("%d", i+1), text, link, page})
	})
}

func main() {
	website := "https://news.ycombinator.com"
	sites := []string{"newest", "front", "jobs"}

	ex, err := os.Executable()
	checkError(err)
	dir := path.Dir(ex)
	now := time.Now()
	timeString := timeutil.Strftime(&now, `%B-%d-%Y`)
	path := path.Join(dir, fmt.Sprintf("hacker-news_%s.csv", timeString))

	file, err := os.Create(path)
	checkError(err)
	writer := csv.NewWriter(file)

	for _, site := range sites {
		url := fmt.Sprintf("%s/%s", website, site)
		writeContent(url, writer, site)
	}

	writer.Flush()
}
