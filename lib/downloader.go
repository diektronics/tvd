package lib

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Download(queue chan *Episode, i int, n Notifier) {
	log.Printf("%d ready for action!\n", i)
	// wait for data
	for ep := range queue {
		parts := strings.Split(ep.Episode, "E")
		season, _ := strconv.Atoi(strings.Trim(parts[0], "S"))

		destination := fmt.Sprintf("%s/%s/Season%d",
			ep.Location,
			ep.Title,
			season)
		filename := fmt.Sprintf("%s - %s.mkv", ep.Title, ep.Episode)
		log.Printf("%d: getting %q %q via %q to be stored in %q\n",
			i,
			ep.Title,
			ep.Episode,
			ep.Link,
			destination)
		cmd := []string{"/home/carter/bin/plowdown",
			"--engine=xfilesharing",
			"--output-directory=" + destination,
			"--printf=%F",
			"--temp-rename",
			ep.Link}
		output, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
		if err != nil {
			log.Println(i, " err: ", err)
                        log.Println(i, " output: ", string(output))
			continue
		}
		parts = strings.Split(strings.TrimSpace(string(output)), "\n")
		oldFilename := parts[len(parts) - 1]
		newFilename := fmt.Sprintf("%s/%s", destination, filename)
		if err := os.Rename(oldFilename, newFilename); err != nil {
			log.Println(i, " err: ", err)
			continue
		}

		log.Printf("%d: %q download complete\n", i, filename)
		n.Notify(newFilename)
	}
}
