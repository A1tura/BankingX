package mql

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"sharedTypes"
)

type Template struct {
	TemplateName string            `json:"template_name"`
	Args         map[string]string `json:"args"`
	To           string            `json:"to"`
}

func SendEmailConfirmationEmail(rabbitmq *amqp.Connection, link, to string) error {
	args := map[string]string{
		"link": link,
	}

	var template Template
	template.TemplateName = "emailConfirmation"
	template.To = to
	template.Args = args

	res, err := json.Marshal(template)
	if err != nil {
		return err
	}

	channel, err := rabbitmq.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	q, err := channel.QueueDeclare(
		"email",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = channel.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(res),
		}); err != nil {
		return err
	}

	return nil
}

func Send2FSignIn(rabbitmq *amqp.Connection, userId int, email string) error {
	channel, err := rabbitmq.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	q, err := channel.QueueDeclare(
		"2fa",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var res sharedtypes.ActionRequest
	res.Email = email
	res.UserId = userId
	res.ActionId = "auth:"
	res.Queue = ""

	resBytes, err := json.Marshal(res)
	if err != nil {
		return err
	}

	if err = channel.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        resBytes,
		}); err != nil {
		return err
	}

	return nil
}
