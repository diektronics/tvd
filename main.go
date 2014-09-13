package main

import (
	"flag"
	"log"
	"os"
	"time"

	"diektronics.com/carter/tvd/common"
	"diektronics.com/carter/tvd/datamanager"
	"diektronics.com/carter/tvd/downloader"
)

var cfgFile = flag.String(
	"cfg",
	os.Getenv("HOME")+"/.config/tvd/config.json",
	"Configuration file in JSON format indicating DB credentials and mailing details.",
)

const waitingTime = time.Duration(20) * time.Minute

func main() {
	flag.Parse()
	c, err := common.GetConfig(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	// Let's use only 4 downloaders to keep bandwidth sane.
	dl := downloader.New(c, 4)
	dm := datamanager.New(c)
	var timestamp *time.Time
	for {
		if shows, timestamp, err := dm.GetMyShows(timestamp); err != nil {
			log.Println(("err:", err))
		} else if len(shows) != 0 {
			dl.Download(shows)
		}

		time.Sleep(waitingTime)
	}
}
