package notifier

import (
	"fmt"
	"net/smtp"
)

func Notify(filename string) {
	addr := "smtp.gmail.com:587"
	auth := smtp.PlainAuth("", "tvd@diektronics.com", "1r3nkA52!!", "smtp.gmail.com")
	from := "tvd@diektronics.com"
	to := []string{"diego.carretero@gmail.com"}
	header := "From: tvd@diektronics.com\nTo: diego.carretero@gmail.com\nSubject: New file!\n\n"
	body := fmt.Sprintf("New download complete: %q", filename)
	content := []byte(header + body)

	if err := smtp.SendMail(addr, auth, from, to, content); err != nil {
		fmt.Println("err: ", err)
	}
}
