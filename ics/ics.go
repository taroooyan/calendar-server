package ics

import (
	"fmt"
	"github.com/taroooyan/esa2ics/esa"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ICalnedar struct {
	Begin       string `ics:"BEGIN"`
	Prodid      string `ics:"PRODID"`
	Version     string `ics:"VERSION"`
	Calscale    string `ics:"CALSCALE"`
	Method      string `ics:"METHOD"`
	Xwrtimezone string `ics:"X-WR-TIMEZONE"`
	Vevent      []Vevent
	End         string `ics:"END"`
}

type Vevent struct {
	Begin       string `ics:"BEGIN"`
	Dtstart     string `ics:"DTSTART"`
	Dtend       string `ics:"DTEND"`
	Uid         string `ics:"UID"`
	Class       string `ics:"CLASS"`
	Description string `ics:"DESCRIPTION"`
	Sequence    string `ics:"SEQUENCE"`
	Status      string `ics:"STATUS"`
	Summary     string `ics:"SUMMARY"`
	Transp      string `ics:"TRANSP"`
	End         string `ics:"END"`
}

func createICS() ICalnedar {
	articles := esa.TakeArticle()
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
		event.Description = strings.Replace(post.BodyMd, "\r\n", "\\n", -1)
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

func PrintICS(w http.ResponseWriter, r *http.Request) {
	ical := createICS()
	icalType := reflect.TypeOf(ical)
	icalValue := reflect.ValueOf(ical)
	// print ICalnedar
	for i := 0; i < icalType.NumField(); i++ {
		icalTag := icalType.Field(i).Tag.Get("ics")
		if icalTag != "" {
			fmt.Fprintf(w, "%s:%s\n", icalTag, icalValue.Field(i).Interface())
		} else {
			// print multi Vevent
			for _, event := range ical.Vevent {
				eventType := reflect.TypeOf(event)
				eventValue := reflect.ValueOf(event)
				for j := 0; j < eventType.NumField(); j++ {
					eventTag := eventType.Field(j).Tag.Get("ics")
					fmt.Fprintf(w, "%s:%s\n", eventTag, eventValue.Field(j).Interface())
				}
			}
		}
	}
}
/*
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
*/
