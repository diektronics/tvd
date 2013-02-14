package downloader

import (
	"diektronics.com/notifier"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Download(title, episode, link, location string) {
	parts := strings.Split(episode, "E")
	season, _ := strconv.Atoi(strings.Trim(parts[0], "S"))

	destination := fmt.Sprintf("%s/%s/Season%d",
		location,
		title,
		season)
	filename := fmt.Sprintf("%s - %s.mkv", title, episode)
	fmt.Printf("getting %q %q via %q to be stored in %q\n",
		title,
		episode,
		link,
		destination)
	cmd := []string{"/usr/local/bin/plowdown",
		"--output-directory=" + destination,
		link}

	if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
		fmt.Println("err: ", err)
		return
	}
	parts = strings.Split(link, "/")
	oldFilename := fmt.Sprintf("%s/%s", destination, strings.Replace(parts[len(parts)-1], ".htm", "", 1))
	newFilename := fmt.Sprintf("%s/%s", destination, filename)
	if err := os.Rename(oldFilename, newFilename); err != nil {
		fmt.Println("err: ", err)
		return
	}

	fmt.Printf("%q download complete\n", filename)
	notifier.Notify(newFilename)
}
