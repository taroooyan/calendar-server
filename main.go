package main

import (
	"fmt"
	"github.com/taroooyan/go-esa/esa"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

// esa.ioから日報カテゴリのすべての記事を取得
func takeArticle() []esa.PostResponse {
	client := esa.NewClient(os.Getenv("ESA_API"))
  articles := []esa.PostResponse{}

	page := "1"
	for {
		query := url.Values{}
		query.Add("in", "日報")
    query.Add("page", page)
		postsResponse, err := client.Post.GetPosts("taroooyan", query)
		if err != nil {
			panic(err)
		}

    articles = append(articles, postsResponse.Posts...)

		if postsResponse.NextPage == nil {
			break
		} else {
			t := postsResponse.NextPage
			page = strconv.FormatFloat(t.(float64), 'G', 4, 64)
		}
	}

  for _, post := range articles {
    fmt.Println(post.Category)
  }

  return articles
}

func main() {
  takeArticle()
	// http.HandleFunc("/calendar.ics", readICS)
	// http.ListenAndServe(":80", nil)
}
