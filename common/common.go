package common

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
