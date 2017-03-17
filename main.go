package main

import (
	"fmt"
	"github.com/taroooyan/go-esa/esa"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
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

func main() {
	// takeArticle()

	http.HandleFunc("/calendar.ics", printICS)
	http.ListenAndServe(":8080", nil)
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
	articles := takeArticle()
	ical := ICalnedar{}

	// initialize iCalendar
	ical.Begin = "VCALENDAR"
	ical.Prodid = "esa/taroooyan"
	ical.Version = "2.0"
	ical.Calscale = "GREGORIAN"
	ical.Method = "PUBLISH"
	ical.Xwrtimezone = "UTC"

	// crate events
	for _, post := range articles {
		fmt.Println(post.Category)

		event := Vevent{}
		event.Begin = "VEVENT"

		// convert "日報/2016/09/13" to 20160913
		dateSplit := strings.Split(post.Category, "/")[1:]
		event.Dtstart = strings.Join(dateSplit, "")

		y, _ := strconv.Atoi(dateSplit[0])
		m, _ := strconv.Atoi(dateSplit[1])
		d, _ := strconv.Atoi(dateSplit[2])
		t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
		t = t.AddDate(0, 0, 1)
		dateEnd := t.Format("20060102")
		event.Dtend = dateEnd

		event.Uid = event.Dtstart
		event.Class = "PUBLISH"
		event.Description = post.BodyMd
		event.Sequence = "0"
		event.Status = "CONFIRMED"
		event.Summary = "日報"
		event.Transp = "TRANSPARENT"
		event.End = "VEVENT"

		ical.Vevent = append(ical.Vevent, event)
	}

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
