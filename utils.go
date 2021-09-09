package main

import (
	"os"
	"path/filepath"
)

func writeHtmlFile(body []byte, filename string) {
	err := os.WriteFile("downloaded_htmls/"+filename, body, 0644)
	if err != nil {
		panic(err)
	}
}

func clear() {
	files, err := filepath.Glob("downloaded_htmls/*.html")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}
