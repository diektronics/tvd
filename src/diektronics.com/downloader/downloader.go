package downloader

import (
	"diektronics.com/notifier"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Episode struct {
	Title    string
	Episode  string
	Link     string
	Location string
}

func Download(queue chan *Episode) {
	for {
		// wait for data
		ep := <-queue

		parts := strings.Split(ep.Episode, "E")
		season, _ := strconv.Atoi(strings.Trim(parts[0], "S"))

		destination := fmt.Sprintf("%s/%s/Season%d",
			ep.Location,
			ep.Title,
			season)
		filename := fmt.Sprintf("%s - %s.mkv", ep.Title, ep.Episode)
		fmt.Printf("getting %q %q via %q to be stored in %q\n",
			ep.Title,
			ep.Episode,
			ep.Link,
			destination)
		cmd := []string{"/usr/local/bin/plowdown",
			"--output-directory=" + destination,
			ep.Link}

		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			fmt.Println("err: ", err)
			return
		}
		parts = strings.Split(ep.Link, "/")
		oldFilename := fmt.Sprintf("%s/%s", destination, strings.Replace(parts[len(parts)-1], ".htm", "", 1))
		newFilename := fmt.Sprintf("%s/%s", destination, filename)
		if err := os.Rename(oldFilename, newFilename); err != nil {
			fmt.Println("err: ", err)
			return
		}

		fmt.Printf("%q download complete\n", filename)
		notifier.Notify(newFilename)
	}
}
