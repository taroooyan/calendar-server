package esa

import (
	"fmt"
	"github.com/taroooyan/go-esa/esa"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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

	for _, post := range articles {
		fmt.Printf("%#v\n", post.Category)
		fmt.Println(strings.Join(strings.Split(post.Category, "/")[1:], ""))

		dateSplit := strings.Split(post.Category, "/")[1:]
		y, _ := strconv.Atoi(dateSplit[0])
		m, _ := strconv.Atoi(dateSplit[1])
		d, _ := strconv.Atoi(dateSplit[2])
		t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
		t = t.AddDate(0, 0, 1)
		fmt.Println(y, m, d)
		fmt.Println(t.Format("20060102"))

	}

	return articles
}
