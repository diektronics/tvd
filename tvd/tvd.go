package tvd

import (
	"log"
	"time"

	"diektronics.com/carter/tvd/common"
	"diektronics.com/carter/tvd/data"
	"diektronics.com/carter/tvd/downloader"
	"diektronics.com/carter/tvd/notifier"
)

type Tvd struct {
	db   *data.Db
	n    *notifier.Notifier
	feed *data.Feed
}

const waitingTime = time.Duration(20) * time.Minute

func New(c *common.Configuration) *Tvd {
	return &Tvd{
		db:   data.NewDb(c),
		n:    notifier.New(c),
		feed: data.NewFeed(c),
	}
}

func reportAndWait(err error) {
	log.Println("err: ", err)
	time.Sleep(waitingTime)
}

func (t *Tvd) Run() {
	// we are not going to get more than 10 eps to download...
	q := make(chan *common.Episode, 10)
	// prepare the downloaders, 4 to not destroy BW
	for i := 0; i < 4; i++ {
		go downloader.Download(q, i, t.n)
	}

	var oldQuery *data.Query
	for {
		query, err := t.feed.Get()
		if err != nil {
			reportAndWait(err)
			continue
		}

		newer, err := query.IsNewerThan(oldQuery)
		if err != nil {
			reportAndWait(err)
			continue
		}

		if newer {
			interestingShows, err := t.db.GetInterestingShows(query)
			if err != nil {
				reportAndWait(err)
				continue
			}

			for _, show := range interestingShows {
				q <- show
			}

			oldQuery = query
		}
		time.Sleep(waitingTime)
	}
}
