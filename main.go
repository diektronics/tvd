package main

import (
	"flag"
	"log"
	"os"

	"diektronics.com/carter/tvd/common"
	"diektronics.com/carter/tvd/tvd"
)

var cfgFile = flag.String(
	"cfg",
	os.Getenv("HOME")+"/.config/tvd/config.json",
	"Configuration file in JSON format indicating DB credentials and mailing details.",
)

func main() {
	flag.Parse()
	c, err := common.GetConfig(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	tvd.New(c).Run()
}
