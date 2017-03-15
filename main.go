package main

import (
	"fmt"
	"github.com/taroooyan/go-esa/esa"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

func saveArticle() {
	// Initialization
	client := esa.NewClient(os.Getenv("ESA_API"))

	page := "1"
	for {
		query := url.Values{}
		query.Add("in", "日報")
		query.Add("page", page)
		query.Add("order", "asc")

		postsResponse, err := client.Post.GetPosts("taroooyan", query)
		if err != nil {
			panic(err)
		}

		for _, post := range postsResponse.Posts {
			fmt.Println(post.Number)
		}

		if postsResponse.NextPage == nil {
			break
		} else {
			t := postsResponse.NextPage
			page = strconv.FormatFloat(t.(float64), 'G', 4, 64)
		}
	}
}

func main() {
	saveArticle()
	// http.HandleFunc("/calendar.ics", ShowICS)
	// http.ListenAndServe(":80", nil)
}
