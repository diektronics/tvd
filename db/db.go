package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"diektronics.com/carter/tvd/common"
	_ "github.com/Go-SQL-Driver/MySQL"
)

type Db struct {
	connectionString string
}

func New(c *common.Configuration) *Db {
	return &Db{
		connectionString: fmt.Sprintf("%s:%s@%s/%s?charset=utf8",
			c.DbUser, c.DbPassword, c.DbServer, c.DbDatabase),
	}
}

func (d *Db) GetMyShows(titles []string) (myShows []*common.Episode, err error) {
	db, err := sql.Open("mysql", d.connectionString)
	if err != nil {
		return
	}
	defer db.Close()

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
		myShows = append(myShows, &common.Episode{
			Title:    name,
			Episode:  s.eps,
			Location: location,
		})
	}
	return
}

func (d *Db) UpdateMyShows(shows []*common.Episode) error {
	db, err := sql.Open("mysql", d.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()
	var lastErr error = nil
	for _, s := range shows {
		dbQuery = fmt.Sprintf("UPDATE series SET latest_ep=%q WHERE name=%q", s.Episode, s.Title)
		_, err = db.Exec(dbQuery)
		if err != nil {
			log.Println(err)
			lastErr = err
			continue
		}
	}

	return lastErr
}
