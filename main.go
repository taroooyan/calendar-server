package main

import (
	"fmt"
	"net/http"
	"os"
)

func ShowICS(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("basic.ics")
	defer file.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := file.Read(buf)
		if n == 0 {
			break
		}
		fmt.Fprintf(w, string(buf[:n]))
	}
}
func main() {
	http.HandleFunc("/calendar.ics", ShowICS)
	http.ListenAndServe(":80", nil)
}
