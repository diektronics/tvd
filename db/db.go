package db

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	"diektronics.com/carter/tvd/common"
	"diektronics.com/carter/tvd/feed"
	_ "github.com/Go-SQL-Driver/MySQL"
)

type Db struct {
	connectionString string
	linkRegexp       string
}

func New(c *common.Configuration) *Db {
	return &Db{
		connectionString: fmt.Sprintf("%s:%s@%s/%s?charset=utf8",
			c.DbUser, c.DbPassword, c.DbServer, c.DbDatabase),
		linkRegexp: c.LinkRegexp,
	}
}

func (d *Db) GetMyShows(data *feed.Data) (myShows []*common.Episode, err error) {
	db, err := sql.Open("mysql", d.connectionString)
	if err != nil {
		return
	}
	defer db.Close()

	type show struct {
		eps string
		it  feed.Item
	}
	shows := make(map[string][]*show)
	titles := []string{}
	for _, s := range data.ItemList {
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
					myShows = append(myShows, &common.Episode{
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

func parenthesize(str string) string {
	// RlsBB doesn't use parenthesis when a Series name has a year attached to it,
	// eg. Castle (2009), but the DB has them.
	// So, if "title" ends with four digits, we are going to add
	// parenthesis around it.
	stuff := `\d{4}$`
	epsRegexp := regexp.MustCompile(stuff)
	return epsRegexp.ReplaceAllString(str, "($0)")
}
