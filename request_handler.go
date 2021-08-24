package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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
		panic(err)
	}

	body1, body2 := bytes.NewReader(b), bytes.NewReader(b)

	colectHtml(body1)
	colectlinks(body2)
}

func colectHtml(body1 *bytes.Reader) {
	index, err := os.Create("index.html")
	if err != nil {
		log.Fatal("Error criando o html")
	}

	defer index.Close()
	index.ReadFrom(body1)
}

func colectlinks(body2 *bytes.Reader) {
	doc, err := goquery.NewDocumentFromReader(body2)
	if err != nil {
		log.Fatal("Error no parser pro goquery")
	}

	s := []string{}
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		s = append(s, href)
	})

	vl := []string{}
	for _, href := range s {
		_, err := http.Get(href)
		if err == nil {
			vl = append(vl, href)
		}
	}

	file, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, link := range vl {
		_, _ = datawriter.WriteString(link + "\n")
	}

	datawriter.Flush()
	file.Close()
}
