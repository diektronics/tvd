package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
)

type Episode struct {
	Title    string
	Episode  string
	Link     string
	Location string
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
	LinkRegexp    string
	Feed          string
}

func GetConfig(cfgFile string) (*Configuration, error) {
	cfg, err := os.Open(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("Open: %v", err)
	}
	decoder := json.NewDecoder(cfg)
	c := &Configuration{}
	if err := decoder.Decode(c); err != nil {
		return nil, fmt.Errorf("Decode: %v", err)
	}

	return c, nil
}

func Match(reStr string, s string) (map[string]string, error) {
	re := regexp.MustCompile(reStr)
	matches := re.FindStringSubmatch(s)
	if len(matches) == 0 {
		return nil, errors.New("no matches found")
	}
	ret := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if len(name) == 0 {
			continue
		}
		ret[name] = matches[i]
	}

	return ret, nil
}
