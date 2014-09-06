package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

func getConfig() (*common.Configuration, error) {
	cfg, err := os.Open(*cfgFile)
	if err != nil {
		return nil, fmt.Errorf("Open: %v", err)
	}
	decoder := json.NewDecoder(cfg)
	c := &common.Configuration{}
	if err := decoder.Decode(c); err != nil {
		return nil, fmt.Errorf("Decode: %v", err)
	}

	return c, nil
}

func main() {
	flag.Parse()
	c, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	tvd.New(c).Run()
}
