package notifier

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type Notifier struct {
	Addr      string
	Port      string
	Recipient string
	Sender    string
	Password  string
}

func (n Notifier) Notify(filename string) {
	parts := strings.Split(filename, "/")
	episode := parts[len(parts)-1]

	header := fmt.Sprintf("From: %s\nTo: %s\nSubject: New file! %s\n\n",
		n.Sender, n.Recipient, episode)
	body := fmt.Sprintf("New download complete: %q", filename)
	content := []byte(header + body)

	addrPort := n.Addr
	if n.Port != "" {
		addrPort += ":" + n.Port
	}

	auth := smtp.PlainAuth("", n.Sender, n.Password, n.Addr)
	to := []string{n.Recipient}
	if err := smtp.SendMail(addrPort, auth, n.Sender, to, content); err != nil {
		log.Println("err: ", err)
	}
}
