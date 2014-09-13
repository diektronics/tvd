package datamanager

import (
	"time"

	"diektronics.com/carter/tvd/common"
	"diektronics.com/carter/tvd/db"
	"diektronics.com/carter/tvd/feed"
)

type DataManager struct {
	db         *db.Db
	feedUrl    string
	linkRegexp string
}

func New(c *common.Configuration) *DataManager {
	return &DataManager{
		db:         db.New(c),
		feedUrl:    c.Feed,
		linkRegexp: c.LinkRegexp,
	}
}

func (dm *DataManager) GetMyShows(timestamp *time.Time) ([]*common.Episode, *time.Time, error) {
	feed := feed.New(dm.feedUrl)

	titles, timestamp, err := dm.feed.Get(timestamp)
	if err != nil || titles == nil {
		return nil, timestamp, err
	}

	// 1.Use titles to get db.MyShows
	// 2.With results loop over shows and generate links
	// 3. return

}
