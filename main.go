package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Article struct {
	ID          int `gorm:"primaryKey"`
	Title       string
	Author      string
	Link        string `gorm:"type:text"`
	JournalInfo string
	DOI         string
}

func main() {

	// 数据库相关
	var err error
	db, err := gorm.Open("mysql", "root:Cyl851106@(127.0.0.1:3306)/cambridge_scraper?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 自动迁移
	db.AutoMigrate(&Article{})

	goQuery := func(pageNum int) {
		url := "https://www.cambridge.org/core/search?q=cambrian%20brachiopod*&aggs%5BproductJournal%5D%5Bfilters%5D=56B1B6F705BBEC4F8958383925A06535&pageNum=" + strconv.Itoa(pageNum)
		c := colly.NewCollector()

		//列表爬取
		c.OnHTML("a.part-link", func(e *colly.HTMLElement) {
			// title := strings.Trim(e.Text, "\n")
			link := e.Attr("href")
			// Print link
			//fmt.Printf("Link found: %q -> %s\n", title, link)
			c.Visit(e.Request.AbsoluteURL(link))

		})

		c.OnHTML("div.column__main__left", func(e *colly.HTMLElement) {
			//fmt.Println(e)
			article := Article{
				Title:       "",
				Link:        "",
				DOI:         "",
				Author:      "",
				JournalInfo: "",
			}
			// 爬取 URL 和 Title
			e.ForEach("div#maincontent>h1", func(i int, e *colly.HTMLElement) {
				article.Link = e.Request.URL.String()
				article.Title = e.Text
			})
			// 爬取 DOI
			e.ForEach("div.doi-data>div>a>span.text", func(i int, e *colly.HTMLElement) {
				article.DOI = e.Text
				//fmt.Println(e)
			})
			// 爬取 Author
			var author string
			e.ForEach(".contributor-type__contributor", func(i int, e *colly.HTMLElement) {
				//article.DOI = e.Text

				e.ForEach("a", func(i int, h *colly.HTMLElement) {
					author += h.Text + ", "
				})
			})
			article.Author = strings.TrimRight(author, ", ")

			// 爬取 journalInfo
			var journalInfo string
			e.ForEach("div.content__journal", func(i int, e *colly.HTMLElement) {
				e.ForEach("a", func(i int, h *colly.HTMLElement) {
					journalInfo += h.Text + ", "
				})
				//fmt.Println(e)
			})
			article.JournalInfo = strings.TrimRight(journalInfo, ", ")

			db.Create(&article)
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
}
