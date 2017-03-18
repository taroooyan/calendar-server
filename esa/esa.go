package esa

import (
	"github.com/taroooyan/go-esa/esa"
	"net/url"
	"os"
	"strconv"
)

// esa.ioから日報カテゴリのすべての記事を取得
func TakeArticle() []esa.PostResponse {
	team := os.Getenv("ESA_TEAM")
	client := esa.NewClient(os.Getenv("ESA_API"))
	articles := []esa.PostResponse{}

	page := "1"
	for {
		query := url.Values{}
		query.Add("in", "日報")
		query.Add("page", page)
		postsResponse, err := client.Post.GetPosts(team, query)
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

	return articles
}
