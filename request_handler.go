package main

import (
	"bytes"
	"fmt"
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

	clear()
	writeHtml(body1, "index.html")
	bodyLinks := colectlinks(body2)
	handleBodylinks(bodyLinks)
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

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal("Error while parsing body to goqyery")
	}

	s := []string{}
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		s = append(s, href)
	})

	s = validateLinks(s)

	return s
}

func validateLinks(s []string) []string {
	// TODO: adaptar downloadBodyLinks para fazer isso aqui
	// TODO: fazer a func de go routine parar
	r := []string{}
	for _, link := range s {
		_, err := http.Get(link)
		if err == nil {
			r = append(r, link)
		}
	}

	return r
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
			return
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
		}

		body := bytes.NewReader(b)
		writeHtml(body, strconv.Itoa(i)+".html")

		c <- "Download complete: " + link
	}

	close(c)
}
