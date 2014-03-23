package main

import (
	"diektronics.com/carter/tvd/lib"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func reportAndWait(err error) {
	log.Println("err: ", err)
	time.Sleep(20 * time.Minute)
}

type Configuration struct {
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
		log.Println("err: ", err)
		return
	}

	var c Configuration
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Println("err: ", err)
		return
	}

	// we are not going to get more than 10 eps to download...
	queue := make(chan *lib.Episode, 10)
	n := lib.Notifier{c.MailAddr, c.MailPort, c.MailRecipient,
		c.MailSender, c.MailPassword}
	// prepare the downloaders, 4 to not destroy BW
	for i := 0; i < 4; i++ {
		go lib.Download(queue, i, n)
	}

	var oldQuery *lib.Query
	for {
		query, err := lib.AllShows()
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
			interestingShows, err := lib.InterestingShows(query, c.DbUser,
				c.DbPassword, c.DbServer, c.DbDatabase)
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
