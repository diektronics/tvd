package notifier

import (
	"fmt"
	"net/smtp"
)

type Notifier struct {
	Addr      string
	Port      string
	Recipient string
	Sender    string
	Password  string
}

func (n Notifier) Notify(filename string) {
	auth := smtp.PlainAuth("", n.Sender, n.Password, n.Addr)
	to := []string{n.Recipient}
	header := fmt.Sprintf("From: %s\nTo: %s\nSubject: New file!\n\n", n.Sender, n.Recipient)
	body := fmt.Sprintf("New download complete: %q", filename)
	content := []byte(header + body)

	addrPort := n.Addr
	if n.Port != "" {
		addrPort += ":" + n.Port
	}
	if err := smtp.SendMail(addrPort, auth, n.Sender, to, content); err != nil {
		fmt.Println("err: ", err)
	}
}
