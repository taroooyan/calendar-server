package main

import (
	"fmt"
	"net/http"
	"os"
)

func ShowICS(w http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/calendar.ics", ShowICS)
	http.ListenAndServe(":80", nil)
}
