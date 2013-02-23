package main

import (
	"diektronics.com/data"
	"diektronics.com/downloader"
	"diektronics.com/episode"
	"diektronics.com/notifier"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func reportAndWait(err error) {
	fmt.Println("err: ", err)
	time.Sleep(20 * time.Minute)
}

type Message struct {
	DbUser        string
	DbServer      string
	DbPassword    string
	DbDatabase    string
	MailAddr      string
	MailPort      string
	MailRecipient string
	MailSender    string
	MailPassword  string
}

func main() {
	b, err := ioutil.ReadFile(os.Getenv("HOME") + "/.tvd/config.json")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	var m Message
	err = json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Printf("%#v\n", m)
	// we are not going to get more than 10 eps to download...
	queue := make(chan *episode.Episode, 10)
	n := notifier.Notifier{m.MailAddr, m.MailPort, m.MailRecipient,
		m.MailSender, m.MailPassword}
	// prepare the downloaders, 4 to not destroy BW
	for i := 0; i < 4; i++ {
		go downloader.Download(queue, i, n)
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
			interestingShows, err := data.InterestingShows(query, m.DbUser,
				m.DbPassword, m.DbServer, m.DbDatabase)
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
