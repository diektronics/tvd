package data

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Query struct {
	DateStamp string `xml:"channel>lastBuildDate"`
	ItemList  []Item `xml:"channel>item"`
}

type Item struct {
	Title   string `xml:"title"`
	Content string `xml:"encoded"`
}

func (i Item) Tokenize() (name, episode string) {
	stuff := `S\d\dE\d\d`
	epsRegexp, _ := regexp.Compile(stuff)
	start := epsRegexp.FindIndex([]byte(strings.ToUpper(i.Title)))
	if start == nil {
		name = i.Title
		return
	}
	name = i.Title[:start[0]-1]
	parts := strings.Fields(i.Title[start[0]:])
	episode = parts[0]

	return
}

func (i Item) Link() (link string) {
	name, episode := i.Tokenize()
	titleEp := fmt.Sprintf("%s\\.%s\\.720p",
		strings.ToLower(strings.Replace(name, " ", "\\.", -1)),
		strings.ToLower(episode))
	stuff := `http://netload.in/\w+/` + titleEp
	linkRegexp, _ := regexp.Compile(stuff)
	linkStart := linkRegexp.FindIndex([]byte(strings.ToLower(i.Content)))
	if len(linkStart) != 0 {
		parts := strings.Split(i.Content[linkStart[0]:], "\"")
		link = parts[0]
	}
	return
}

func (q Query) Date() (time.Time, error) {
	format := "Mon, 02 Jan 2006 15:04:05 -0700"
	date := q.DateStamp
	if q.DateStamp == "" {
		date = format
	}
	return time.Parse(format, date)
}

func (q Query) After(otherQ Query) (bool, error) {
	parsedTime, err := q.Date()
	if err != nil {
		return false, err
	}

	otherParsedTime, err := otherQ.Date()
	if err != nil {
		return false, err
	}

	return parsedTime.After(otherParsedTime), nil
}

func Shows() (Query, error) {
	var q Query
	stuff, err := http.Get("http://www.rlsbb.com/category/tv-shows/feed/")
	if err != nil {
		return q, err
	}
	defer stuff.Body.Close()

	body, err := ioutil.ReadAll(stuff.Body)
	if err != nil {
		return q, err
	}

	err = xml.Unmarshal([]byte(string(body)), &q)
	if err != nil {
		return q, err
	}

	return q, nil
}
