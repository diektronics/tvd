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

func (dm *DataManager) GetMyShows(timestamp time.Time) (myShows []*common.Episode, newTimestamptime.Time, err error) {
	feed := feed.New(dm.feedUrl)
	var titles []string
	titles, newTimestamp, err = dm.feed.Update(timestamp)
	if err != nil || len(titles) == 0 {
		return
	}

	myShows, err = dm.db.GetMyShows(titles)
	if err != nil || len(myShows) == 0 {
		return
	}

	myShows, err = dm.feed.SetLinks(myShows)
	f err != nil || len(myShows) == 0 {
		return
	}

	if err = dm.db.UpdateMyShows(myShows); err != nil {
		return
	}

	return
}
