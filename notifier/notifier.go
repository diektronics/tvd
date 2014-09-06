package notifier

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"diektronics.com/carter/tvd/common"
)

type Notifier struct {
	addr      string
	port      string
	recipient string
	sender    string
	password  string
}

func New(c *common.Configuration) *Notifier {
	return &Notifier{
		addr:      c.MailAddr,
		port:      c.MailPort,
		recipient: c.MailRecipient,
		sender:    c.MailSender,
		password:  c.MailPassword,
	}
}

func (n Notifier) Notify(filename string) {
	parts := strings.Split(filename, "/")
	episode := parts[len(parts)-1]

	header := fmt.Sprintf("From: %s\nTo: %s\nSubject: New file! %s\n\n",
		n.sender, n.recipient, episode)
	body := fmt.Sprintf("New download complete: %q", filename)
	content := []byte(header + body)

	addrPort := n.addr
	if n.port != "" {
		addrPort += ":" + n.port
	}

	auth := smtp.PlainAuth("", n.sender, n.password, n.addr)
	to := []string{n.recipient}
	if err := smtp.SendMail(addrPort, auth, n.sender, to, content); err != nil {
		log.Println("err: ", err)
	}
}
