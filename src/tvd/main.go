package main

import (
	"diektronics.com/data"
	"diektronics.com/downloader"
	"diektronics.com/episode"
	"fmt"
	"time"
)

func reportAndWait(err error) {
	fmt.Println("err: ", err)
	time.Sleep(20 * time.Minute)
}

func main() {
	// we are not going to get more than 10 eps to download...
	var queue = make(chan *episode.Episode, 10)
	// prepare the downloaders, 4 to not destroy BW
	for i := 0; i < 4; i++ {
		go downloader.Download(queue, i)
	}

	var oldQuery *data.Query
	for {
		query, err := data.AllShows()
		if err != nil {
			reportAndWait(err)
			continue
		}

		newer := true
		if oldQuery != nil {
			newer, err = query.After(*oldQuery)
			if err != nil {
				reportAndWait(err)
				continue
			}
		}

		if newer {
			interestingShows, err := data.InterestingShows(query)
			if err != nil {
				reportAndWait(err)
				continue
			}

			for _, show := range interestingShows {
				queue <- show
			}

			oldQuery = query
		}
		time.Sleep(20 * time.Minute)
	}
}
