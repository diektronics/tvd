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
	dl   *downloader.Downloader
}

const waitingTime = time.Duration(20) * time.Minute

func New(c *common.Configuration) *Tvd {
	return &Tvd{
		db:   db.New(c),
		feed: feed.New(c),
		dl:   downloader.New(c),
	}
}

func logAndWait(err error) {
	log.Println("err: ", err)
	time.Sleep(waitingTime)
}

func (t *Tvd) Run() {
	// Start just 4 workers to not kill bandwidth.
	t.dl.Start(4)
	var oldData *feed.Data
	for {
		data, err := t.feed.Get()
		if err != nil {
			logAndWait(err)
			continue
		}

		newer, err := data.IsNewerThan(oldData)
		if err != nil {
			logAndWait(err)
			continue
		}

		if newer {
			shows, err := t.db.GetMyShows(data)
			if err != nil {
				logAndWait(err)
				continue
			}

			for _, show := range shows {
				t.dl.Queue <- show
			}

			oldData = data
		}
		time.Sleep(waitingTime)
	}
}
