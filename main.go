package main

import (
    "flag"
    "fmt"
    "bufio"
    "net/http"
)

func main() {

    urlptr := flag.String("url", "http://google.com", "url to be crawled")

    flag.Parse()

    resp, err := http.Get(*urlptr)
    if err != nil {
	panic(err)
    }

    defer resp.Body.Close()

    fmt.Println("Resposnse Status", resp.Status)

    scanner := bufio.NewScanner(resp.Body)
    for i := 0; scanner.Scan(); i++ {
	fmt.Println(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
	panic(err)
    }
}
