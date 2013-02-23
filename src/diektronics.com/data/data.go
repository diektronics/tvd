package data

import (
	"database/sql"
	"diektronics.com/episode"
	"encoding/xml"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
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

func (i Item) Tokenize() (name, eps string) {
	stuff := `S\d\dE\d\d`
	epsRegexp, _ := regexp.Compile(stuff)
	start := epsRegexp.FindIndex([]byte(strings.ToUpper(i.Title)))
	if start == nil {
		name = i.Title
		return
	}
	name = i.Title[:start[0]-1]
	parts := strings.Fields(i.Title[start[0]:])
	eps = parts[0]

	return
}

func (i Item) Link() (link string) {
	name, eps := i.Tokenize()
	titleEp := fmt.Sprintf("%s\\.%s.*\\.720p",
		strings.ToLower(strings.Replace(name, " ", "\\.", -1)),
		strings.ToLower(eps))
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

func AllShows() (q *Query, err error) {
	stuff, err := http.Get("http://www.rlsbb.com/category/tv-shows/feed/")
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
	// So, it "title" ends with four digits, we are going to add
	// parenthesis around it.
	stuff := `\d\d\d\d$`
	epsRegexp, _ := regexp.Compile(stuff)
	return epsRegexp.ReplaceAllString(str, "($0)")
}

func InterestingShows(query *Query, user, password, server, database string) (interestingShows []*episode.Episode, err error) {
	connectionString := fmt.Sprintf("%s:%s@%s/%s?charset=utf8",
		user, password, server, database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return
	}
	defer db.Close()

	for _, show := range query.ItemList {
		title, eps := show.Tokenize()
		title = parenthesize(title)

		dbQuery := fmt.Sprintf("SELECT latest_ep, location FROM series where name=%q", title)
		var rows *sql.Rows
		rows, err = db.Query(dbQuery)
		if err != nil {
			return
		}

		var latest_ep string
		var location string

		// Fetch rows. Only one results, if any
		for rows.Next() {
			// Scan the value to string
			err = rows.Scan(&latest_ep, &location)
			if err != nil {
				return
			}
			if latest_ep < eps {
				fmt.Printf("title: %q episode: %q latest_ep: %q\n", title, eps, latest_ep)
				link := show.Link()

				if link != "" {
					fmt.Printf("link: %q\n", link)
					fmt.Println("update latest_ep in DB")
					dbQuery = fmt.Sprintf("UPDATE series SET latest_ep=%q WHERE name=%q", eps, title)
					_, err = db.Exec(dbQuery)
					if err != nil {
						return
					}
					fmt.Println("download the thing")
					episodeData := episode.Episode{title, eps, link, location}
					interestingShows = append(interestingShows, &episodeData)
				}
			}

		}

	}

	return
}
