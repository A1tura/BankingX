package email

import (
	"email/config"
	"fmt"
	"net/smtp"
	"strconv"
)

type Email struct {
	From    string
	To      string
	Subject string
	Message string
}

func (email *Email) SendEmail() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

    msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html\n\r\r\n%s\r\n", email.From, email.To, email.Subject, email.Message)

	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)

	if err := smtp.SendMail(cfg.Host+":"+strconv.Itoa(cfg.Port), auth, email.From, []string{email.To}, []byte(msg)); err != nil {
		return err
	}

	return nil
}
