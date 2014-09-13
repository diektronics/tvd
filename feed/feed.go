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
	return time.Parse(format, timestamp)
}

func parenthesize(str string) string {
	// RlsBB doesn't use parenthesis when a Series name has a year attached to it,
	// eg. Castle (2009), but the DB has them.
	// So, if "title" ends with four digits, we are going to add
	// parenthesis around it.
	stuff := `\d{4}$`
	epsRegexp := regexp.MustCompile(stuff)
	return epsRegexp.ReplaceAllString(str, "($0)")
}

func (d data) IsNewerThan(time.Time timestamp) (bool, error) { 
	parsedTime, err := d.Date()
	if err != nil {
		return false, err
	}

	return parsedTime.After(timestamp), nil
}

func (f Feed) Update(timestamp time.Time) (titles []string, newTimestamp time.Time, err error) {
	newTimestamp = timestamp
	stuff, err := http.Get(f.url)
	if err != nil {
		return
	}
	defer stuff.Body.Close()

	body, err := ioutil.ReadAll(stuff.Body)
	if err != nil {
		return
	}

	var d *data
	err = xml.Unmarshal([]byte(string(body)), &d)
	if err != nil {
		return
	}

	if !d.IsNewerThan(timestamp) {
		return
	}

	newTimestamp := date(d.DateStamp)
	for _, entry := range data.ItemList {
		title, eps := entry.Tokenize()
		title = parenthesize(title)
		f.shows[title] = f.append(shows[title], &show{eps, entry})
		titles = append(titles, fmt.Sprintf("%q", title))
	}
	return
}
// func (f Feed) SetLinks(shows []*common.Episode) ([]*common.Episode, error) {}
// 	// We range the array in reverse because episodes are added on the top of the feed,
// 	// and when a show has two episodes back to back, we will first find the newest one.
// 	for i := len(shows[name]) - 1; i >= 0; i-- {
// 		s := shows[name][i]
// 		if latest_ep < s.eps {
// 			log.Printf("title: %q episode: %q latest_ep: %q\n", name, s.eps, latest_ep)
// 			link := s.it.Link(d.linkRegexp)

// 			if len(link) != 0 {
// 				log.Printf("link: %q\n", link)
// 				log.Println("update latest_ep in DB")
//
// 				log.Println("download the thing")

// 			}
// 		}
// 	}
// }
