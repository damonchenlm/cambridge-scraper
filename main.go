package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

type Article struct {
	Title string
}

func main() {
	// Instantiate default collector
	sum := 0
	fName := "cambridge.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	goQuery := func(pageNum int) {
		url := "https://www.cambridge.org/core/search?q=cambrian%20brachiopod*&aggs%5BproductJournal%5D%5Bfilters%5D=56B1B6F705BBEC4F8958383925A06535&pageNum=" + strconv.Itoa(pageNum)
		c := colly.NewCollector()
		// On every a element which has href attribute call callback

		//列表爬取
		c.OnHTML("a.part-link", func(e *colly.HTMLElement) {
			// title := strings.Trim(e.Text, "\n")
			link := e.Attr("href")
			// Print link
			//fmt.Printf("Link found: %q -> %s\n", title, link)
			c.Visit(e.Request.AbsoluteURL(link))

		})

		// 进入列表
		// 爬取 DOI 号
		// c.OnHTML("div.doi-data>div>a>span.text", func(e *colly.HTMLElement) {
		// 	fmt.Println(e.Text)
		// })
		// 爬取 Title
		c.OnHTML("div#maincontent>h1", func(e *colly.HTMLElement) {
			fmt.Println(e.Text)
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})

		// Set error handler
		c.OnError(func(r *colly.Response, err error) {
			fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		})

		// c.OnHTML("a[aria-label~=Next]", func(e *colly.HTMLElement) {
		// 	link := e.Attr("href")
		// 	fmt.Println(e)
		// 	c.Visit(e.Request.AbsoluteURL(link))
		// })

		// Start scraping on https://hackerspaces.org
		c.Visit(url)
	}

	for i := 1; i <= 34; i++ {
		goQuery(i)
	}
	fmt.Println("===================")
	fmt.Println(sum)
}
