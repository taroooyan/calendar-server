package main

import (
	"fmt"
	"github.com/taroooyan/go-esa/esa"
	"net/http"
	"net/url"
	"os"
	"reflect"
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
	// createICS()
	// takeArticle()
	http.HandleFunc("/calendar.ics", printICS)
	http.ListenAndServe(":80", nil)
}

type ICalnedar struct {
	Begin       string   `ick:"BEGIN"`
	Prodid      string   `ick:"PRODID"`
	Version     string   `ick:"VERSION"`
	Calscale    string   `ick:"CALSCALE"`
	Method      string   `ick:"METHOD"`
	Xwrtimezone string   `ick:"X-WR-TIMEZONE"`
	Vevent      []Vevent `ick:"_VEVENT"`
	End         string   `ick:"END"`
}

type Vevent struct {
	Begin   string `ick:"BEGIN"`
	Dtstart string `ick:"DTSTART"`
	Dtend   string `ick:"DTEND"`
	// Dtstamp      string `ick:"DTSTAMP`
	Uid   string `ick:"UID"`
	Class string `ick:"CLASS"`
	// Created      string `ick:"CREATED"`
	Description string `ick:"DESCRIPTION"`
	// LastModified string `ick:"LAST-MODIFIED"`
	Sequence string `ick:"SEQUENCE"`
	Status   string `ick:"STATUS"`
	Summary  string `ick:"SUMMARY"`
	Transp   string `ick:"TRANSP"`
	End      string `ick:"END"`
}

func createICS() ICalnedar {
	ical := ICalnedar{}

	ical.Begin = "VCALENDAR"
	ical.Prodid = "taroooyan"
	ical.Version = "2.0"
	ical.Calscale = "GREGORIAN"
	ical.Method = "PUBLISH"
	ical.Xwrtimezone = "UTC"

	// crate event
	event := Vevent{}
	event.Begin = "VEVENT"
	event.Dtstart = "20170320"
	event.Dtend = "20170321"
	// event.Dtstamp = "20170313T223209Z"
	event.Uid = "aaaaaaaa"
	event.Class = "PUBLISH"
	// event.Created = "20150421T224828Z"
	event.Description = "test1"
	// event.LastModified = "20150421T224828Z"
	event.Sequence = "0"
	event.Status = "CONFIRMED"
	event.Summary = "TEST"
	event.Transp = "TRANSPARENT"
	event.End = "VEVENT"

	ical.Vevent = append(ical.Vevent, event)

	// crate event
	event.Dtstart = "20170321"
	event.Dtend = "20170322"
	event.Uid = "bbbbbbbbb"
	event.Description = "test2"

	ical.Vevent = append(ical.Vevent, event)

	ical.End = "VCALENDAR"
	return ical
}

func printICS(w http.ResponseWriter, r *http.Request) {
	ical := createICS()
	icalType := reflect.TypeOf(ical)
	icalValue := reflect.ValueOf(ical)
	for i := 0; i < icalType.NumField(); i++ {
		icalTag := icalType.Field(i).Tag.Get("ick")
		if icalTag != "_VEVENT" {
			fmt.Fprintf(w, "%s:%s\n", icalTag, icalValue.Field(i).Interface())
		} else {
			for _, event := range ical.Vevent {
				eventType := reflect.TypeOf(event)
				eventValue := reflect.ValueOf(event)
				for j := 0; j < eventType.NumField(); j++ {
					eventTag := eventType.Field(j).Tag.Get("ick")
					fmt.Fprintf(w, "%s:%s\n", eventTag, eventValue.Field(j).Interface())
				}
			}
		}
	}
}
