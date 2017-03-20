package main

import (
	"fmt"
	"github.com/taroooyan/esa2ics/ics"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/calendar.ics", ics.PrintICS)
	http.ListenAndServe(":8080", nil)
}
