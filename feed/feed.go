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

type Data struct {
	DateStamp string `xml:"channel>lastBuildDate"`
	ItemList  []Item `xml:"channel>item"`
}

type Item struct {
	Title   string `xml:"title"`
	Content string `xml:"encoded"`
}

type Feed struct {
	url string
}

func New(c *common.Configuration) *Feed {
	return &Feed{c.Feed}
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

func (d Data) Date() (time.Time, error) {
	format := "Mon, 02 Jan 2006 15:04:05 -0700"
	date := d.DateStamp
	if d.DateStamp == "" {
		date = format
	}
	return time.Parse(format, date)
}

func (d Data) IsNewerThan(otherD *Data) (bool, error) {
	if otherD == nil {
		return true, nil
	}

	parsedTime, err := d.Date()
	if err != nil {
		return false, err
	}

	otherParsedTime, err := otherD.Date()
	if err != nil {
		return false, err
	}

	return parsedTime.After(otherParsedTime), nil
}

func (f Feed) Get() (d *Data, err error) {
	stuff, err := http.Get(f.url)
	if err != nil {
		return
	}
	defer stuff.Body.Close()

	body, err := ioutil.ReadAll(stuff.Body)
	if err != nil {
		return
	}

	//fmt.Printf("%s\n", body)

	err = xml.Unmarshal([]byte(string(body)), &d)
	if err != nil {
		return
	}

	return
}
