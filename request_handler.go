package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func colect(urlptr *string) {
	resp, err := http.Get(*urlptr)
	if err != nil {
		log.Fatal("Error no request")
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error:", err)
	}

	body1, body2 := bytes.NewReader(b), bytes.NewReader(b)

	writeHtml(body1, "index.html")
	bodyLinks := colectlinks(body2)
	downloadlinks(bodyLinks)
}

func writeHtml(body *bytes.Reader, filename string) {
	index, err := os.Create(filename)
	if err != nil {
		log.Fatal("Error while creating html:", filename)
	}

	defer index.Close()
	index.ReadFrom(body)
}

func clear() {
	files, err := filepath.Glob("*.html")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func colectlinks(body *bytes.Reader) []string {
	clear()

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

func downloadlinks(bodyLinks []string) {
	for i, href := range bodyLinks {
		resp, err := http.Get(href)
		if err == nil {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Error:", err)
			}
			body := bytes.NewReader(b)
			writeHtml(body, strconv.Itoa(i)+".html")
		}
	}
}
