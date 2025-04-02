package rabbitmq

import (
	"email/email"
	"email/templates"
	"encoding/json"
	"errors"
	"os"

	ampq "github.com/rabbitmq/amqp091-go"
)

type Template struct {
	TemplateName string            `json:"template_name"`
	Args         map[string]string `json:"args"`
	To           string            `json:"to"`
}

func Listen() error {
	conStr, exist := os.LookupEnv("RABBITMQ")
	if !exist {
		return errors.New("ENV RABBITMQ do not exist")
	}

	con, err := ampq.Dial(conStr)
	if err != nil {
		return err
	}

	channel, err := con.Channel()
	if err != nil {
		return err
	}

	errChan := make(chan error)

	msgs, err := channel.Consume(
		"email",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	var forever chan struct{}

	go func() {
		for msg := range msgs {
			var templateRequest Template

			if err := json.Unmarshal(msg.Body, &templateRequest); err != nil {
				continue
			}

			template, err := templates.ApplyTemplate(templateRequest.TemplateName, templateRequest.Args)
			if err != nil {
				errChan <- err
				return
			}

			mail := email.Email{
				From:    template.Sender,
				To:      templateRequest.To,
				Subject: template.Subject,
				Message: template.Message,
			}

			if err := mail.SendEmail(); err != nil {
				errChan <- err
				return
			}

		}
		close(errChan)
	}()

    select {
    case err := <- errChan:
        return err
    case <- forever:
        return nil
    }
}
