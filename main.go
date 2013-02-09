package main

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

	name    string
	episode string
	link    string
}

func (i Item) Tokenize() {
	if i.episode == "" {
		stuff := `S\d\dE\d\d`
		epsRegexp, _ := regexp.Compile(stuff)
		start := epsRegexp.FindIndex([]byte(i.Title))
		i.name = i.Title[:start[0]-1]
		fmt.Printf("%q\n", i.name)
		parts := strings.Fields(i.Title[start[0]:])
		i.episode = parts[0]
		fmt.Printf("%q\n", i.episode)
	}
}

func (q Query) Date() (time.Time, error) {
	format := "Mon, 02 Jan 2006 15:04:05 -0700"
	date := q.DateStamp
	if q.DateStamp == "" {
		date = format
	}
	return time.Parse(format, date)
}
func main() {
	var oldQuery Query
	//for {
	stuff, err := http.Get("http://www.rlsbb.com/category/tv-shows/feed/")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	defer stuff.Body.Close()
	body, err := ioutil.ReadAll(stuff.Body)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Printf("%s\n", body)
	var q Query
	err = xml.Unmarshal([]byte(string(body)), &q)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	parsedTime, err := q.Date()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	oldParsedTime, err := oldQuery.Date()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	if parsedTime.After(oldParsedTime) {
		fmt.Println(parsedTime)
		for _, show := range q.ItemList {
			fmt.Println(show.Title)
			show.Tokenize()
		}
		oldQuery = q
	}
	//time.Sleep(20 * time.Minute)
	//}
}
