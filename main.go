package main

import (
	"fmt"
	"github.com/taroooyan/esa2ics/ics"
	"net/http"
	"os"
)

func readICS(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("basic.ics")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, string(buf))
	}
}

func main() {
	// takeArticle()

	http.HandleFunc("/calendar.ics", ics.PrintICS)
	http.ListenAndServe(":8080", nil)
}
