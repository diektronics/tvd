package feed

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"diektronics.com/carter/tvd/common"
)

type data struct {
	DateStamp string `xml:"channel>lastBuildDate"`
	ItemList  []Item `xml:"channel>item"`
}

type Item struct {
	Title   string `xml:"title"`
	Content string `xml:"encoded"`
}

type show struct {
	eps string
	it  Item
}

type Feed struct {
	url    string
	shows  map[string][]*show
}

func New(url string) *Feed {
	return &Feed{url, make(map[string][]*show)}
}

func (i Item) Tokenize() (string, string) {
	reStr := `(?P<name>.*)\s+(?P<eps>S\d{2}E\d{2})`
	ret, err := common.Match(reStr, i.Title)
	if err != nil {
		return i.Title, ""
	}
	return ret["name"], ret["eps"]
}

func (i Item) Link(linkRegexp string) string {
	name, eps := i.Tokenize()
	titleEp := fmt.Sprintf("%s\\.%s\\.720p.*\\.mkv",
		strings.ToLower(strings.Replace(name, " ", "\\.", -1)),
		strings.ToLower(eps))
	reStr := "(?i)(?P<link>" + linkRegexp + titleEp + ")"
	ret, err := common.Match(reStr, i.Content)
	if err != nil {
		return ""
	}
	return ret["link"]
}

func date(timestamp string) (time.Time, error) {
	format := "Mon, 02 Jan 2006 15:04:05 -0700"
	if timestamp == "" {
		timestamp = format
	}
	return &time.Parse(format, timestamp)
}

func (d data) IsNewerThan(time.Time *timestamp) (bool, error) {
	if timestamp == nil {
		return true, nil
	}

	parsedTime, err := d.Date()
	if err != nil {
		return false, err
	}

	return parsedTime.After(timestamp), nil
}

func (f Feed) Update(timestamp *time.Time) ([]string, *time.Time, error) {
	stuff, err := http.Get(f.url)
	if err != nil {
		return nil, timestamp, err
	}
	defer stuff.Body.Close()

	body, err := ioutil.ReadAll(stuff.Body)
	if err != nil {
		return nil, timestamp, err
	}

	var d *data
	err = xml.Unmarshal([]byte(string(body)), &d)
	if err != nil {
		return nil, timestamp, err
	}

	if !d.IsNewerThan(timestamp) {
		return nil, timestamp, nil
	}

	newTimestamp := date(d.DateStamp)
	for _, entry := range data.ItemList {
		title, eps := entry.Tokenize()
		title = parenthesize(title)
		f.shows[title] = f.append(shows[title], &show{eps, entry})
		titles = append(titles, fmt.Sprintf("%q", title))
	}
	return titles, newTimestamp, nil
}
