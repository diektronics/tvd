package main

import (
	"database/sql"
	"diektronics.com/data"
	"diektronics.com/downloader"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"time"
)

func main() {
	var oldQuery data.Query
	for {
		query, err := data.Shows()
		if err != nil {
			fmt.Println("err: ", err)
			return
		}

		newer, err := query.After(oldQuery)
		if err != nil {
			fmt.Println("err: ", err)
			return
		}

		if newer {
			db, err := sql.Open("mysql", "tvd:tvd@/tvd?charset=utf8")
			if err != nil {
				fmt.Println("err: ", err)
				return
			}
			defer db.Close()

			for _, show := range query.ItemList {
				title, episode := show.Tokenize()
				dbQuery := fmt.Sprintf("SELECT latest_ep, location FROM series where name=%q", title)
				rows, err := db.Query(dbQuery)
				if err != nil {
					fmt.Println("err: ", err)
					return
				}

				var latest_ep string
				var location string

				// Fetch rows. Only one results, if any
				for rows.Next() {
					// Scan the value to string
					err = rows.Scan(&latest_ep, &location)
					if err != nil {
						fmt.Println("err: ", err)
						return
					}
					if latest_ep < episode {
						fmt.Printf("title: %q episode: %q latest_ep: %q\n", title, episode, latest_ep)
						link := show.Link()

						if link != "" {
							fmt.Printf("link: %q\n", link)
							fmt.Println("update latest_ep in DB")
							dbQuery = fmt.Sprintf("UPDATE series SET latest_ep=%q WHERE name=%q", episode, title)
							_, err = db.Exec(dbQuery)
							if err != nil {
								fmt.Println("err: ", err)
								return
							}
							fmt.Println("download the thing")
							go downloader.Download(title, episode, link, location)
						}
					}

				}

			}

			oldQuery = query
		}
		time.Sleep(20 * time.Minute)
	}
}
