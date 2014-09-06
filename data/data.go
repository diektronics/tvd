package data

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"diektronics.com/carter/tvd/common"
	_ "github.com/Go-SQL-Driver/MySQL"
)

type Query struct {
	DateStamp string `xml:"channel>lastBuildDate"`
	ItemList  []Item `xml:"channel>item"`
}

type Item struct {
	Title   string `xml:"title"`
	Content string `xml:"encoded"`
}

type Feed struct {
	s string
}

type Db struct {
	connectionString string
	linkRegexp       string
}

func NewFeed(c *common.Configuration) *Feed {
	return &Feed{c.Feed}
}

func NewDb(c *common.Configuration) *Db {
	return &Db{
		connectionString: fmt.Sprintf("%s:%s@%s/%s?charset=utf8",
			c.DbUser, c.DbPassword, c.DbServer, c.DbDatabase),
		linkRegexp: c.LinkRegexp,
	}
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

func (q Query) Date() (time.Time, error) {
	format := "Mon, 02 Jan 2006 15:04:05 -0700"
	date := q.DateStamp
	if q.DateStamp == "" {
		date = format
	}
	return time.Parse(format, date)
}

func (q Query) IsNewerThan(otherQ *Query) (bool, error) {
	if otherQ == nil {
		return true, nil
	}

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

func (f Feed) Get() (q *Query, err error) {
	stuff, err := http.Get(f.s)
	if err != nil {
		return
	}
	defer stuff.Body.Close()

	body, err := ioutil.ReadAll(stuff.Body)
	if err != nil {
		return
	}

	//fmt.Printf("%s\n", body)

	err = xml.Unmarshal([]byte(string(body)), &q)
	if err != nil {
		return
	}

	return
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

func (d *Db) GetInterestingShows(query *Query) (interestingShows []*common.Episode, err error) {
	db, err := sql.Open("mysql", d.connectionString)
	if err != nil {
		return
	}
	defer db.Close()

	type show struct {
		eps string
		it  Item
	}
	shows := make(map[string][]*show)
	titles := []string{}
	for _, s := range query.ItemList {
		title, eps := s.Tokenize()
		title = parenthesize(title)
		shows[title] = append(shows[title], &show{eps, s})
		titles = append(titles, fmt.Sprintf("%q", title))
	}

	dbQuery := fmt.Sprintf("SELECT name, latest_ep, location FROM series where name IN (%s)", strings.Join(titles, ","))
	var rows *sql.Rows
	rows, err = db.Query(dbQuery)
	if err != nil {
		return
	}

	for rows.Next() {
		var name string
		var latest_ep string
		var location string
		err = rows.Scan(&name, &latest_ep, &location)
		if err != nil {
			return
		}
		// We range the array in reverse because episodes are added on the top of the feed,
		// and when a show has two episodes back to back, we will first find the newest one.
		for i := len(shows[name]) - 1; i >= 0; i-- {
			s := shows[name][i]
			if latest_ep < s.eps {
				log.Printf("title: %q episode: %q latest_ep: %q\n", name, s.eps, latest_ep)
				link := s.it.Link(d.linkRegexp)

				if len(link) != 0 {
					log.Printf("link: %q\n", link)
					log.Println("update latest_ep in DB")
					dbQuery = fmt.Sprintf("UPDATE series SET latest_ep=%q WHERE name=%q", s.eps, name)
					_, err = db.Exec(dbQuery)
					if err != nil {
						return
					}
					log.Println("download the thing")
					interestingShows = append(interestingShows, &common.Episode{
						Title:    name,
						Episode:  s.eps,
						Link:     link,
						Location: location,
					})
				}
			}
		}
	}

	return
}
