package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func colect(urlptr *string) {
	resp, err := http.Get(*urlptr)
	if err != nil {
		log.Fatal("Error on main url request!")
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error:", err)
	}

	body1, body2 := b, bytes.NewReader(b)

	clear()
	writeHtmlFile(body1, "index.html")
	bodyLinks := colectlinks(body2)
	handleBodylinks(bodyLinks)
}

func colectlinks(body *bytes.Reader) []string {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal("Error while parsing body to goqyery")
	}

	s := []string{}
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		s = append(s, href)
	})

	return s
}

func handleBodylinks(bodyLinks []string) {
	c := make(chan string)
	go downloadBodyLinks(bodyLinks, c)

	for l := range c {
		fmt.Println(l)
	}
}

func downloadBodyLinks(bodyLinks []string, c chan string) {
	for i, link := range bodyLinks {
		resp, err := http.Get(link)
		if err != nil {
			c <- "Error downloading: " + link
			continue
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error:", err)
		}

		writeHtmlFile(b, strconv.Itoa(i)+".html")

		c <- "Download complete: " + link
	}

	close(c)
}
