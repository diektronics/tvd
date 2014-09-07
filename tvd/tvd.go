package tvd

import (
	"log"
	"time"

	"diektronics.com/carter/tvd/common"
	"diektronics.com/carter/tvd/db"
	"diektronics.com/carter/tvd/downloader"
	"diektronics.com/carter/tvd/feed"
)

type Tvd struct {
	db   *db.Db
	feed *feed.Feed
	d    *downloader.Downloader
	q    chan *common.Episode
}

const waitingTime = time.Duration(20) * time.Minute

func New(c *common.Configuration) *Tvd {
	// we are not going to get more than 10 eps to download...
	q := make(chan *common.Episode, 10)
	return &Tvd{
		db:   db.New(c),
		feed: feed.New(c),
		d:    downloader.New(c, q),
		q:    q,
	}
}

func reportAndWait(err error) {
	log.Println("err: ", err)
	time.Sleep(waitingTime)
}

func (t *Tvd) Run() {
	// Start just 4 workers to not kill bandwidth.
	t.d.Start(4)
	var oldData *feed.Data
	for {
		data, err := t.feed.Get()
		if err != nil {
			reportAndWait(err)
			continue
		}

		newer, err := data.IsNewerThan(oldData)
		if err != nil {
			reportAndWait(err)
			continue
		}

		if newer {
			interestingShows, err := t.db.GetInterestingShows(data)
			if err != nil {
				reportAndWait(err)
				continue
			}

			for _, show := range interestingShows {
				t.q <- show
			}

			oldData = data
		}
		time.Sleep(waitingTime)
	}
}
