package main

import (
	"flag"
)

func main() {

	urlptr := flag.String("url", "https://google.com", "url to be crawled")

	flag.Parse()

	colect(urlptr)

}
