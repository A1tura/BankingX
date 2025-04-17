package service

import (
	"2fa/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type Request struct {
	UserId   int    `json:"userId"`
	Email    string `json:"email"`
	ActionId string `json:"actionId"`
	Queue    string `json:"queue"`
}

type EmailTemplate struct {
	TemplateName string            `json:"template_name"`
	Args         map[string]string `json:"args"`
	To           string            `json:"to"`
}

func Start(rabbitmq *amqp.Connection, redis *redis.Client) error {
	channel, err := rabbitmq.Channel()
	if err != nil {
		return err
	}

	errChan := make(chan error)

	msgs, err := channel.Consume(
		"2fa",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	emailQueue, err := channel.QueueDeclare(
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

	if err != nil {
		return err
	}

	var forever chan struct{}

	go func() {
		for msg := range msgs {
			var req Request
			if err := json.Unmarshal(msg.Body, &req); err != nil {
				continue
			}
			code, err := utils.GenerateConfirmationCode()
			if err != nil {
				continue
			}

			// send code to user throw email
			args := map[string]string{
				"code": code,
			}
			var template EmailTemplate
			template.To = req.Email
			template.TemplateName = "2fa"
			template.Args = args

			templateBytes, err := json.Marshal(template)
			if err != nil {
				continue
			}

			if err := channel.PublishWithContext(context.Background(),
				"",
				emailQueue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        templateBytes,
				}); err != nil {
				continue
			}

			// if it's 2fa from auth service, there is no userId becuse user it's not auth yet, so we store it without user id, but only for that case
			if strings.Split(req.ActionId, ":")[0] == "auth" {
				key := fmt.Sprintf("auth:%s", code)
				cmd := redis.Set(context.Background(), key, req.UserId, time.Minute*10)
				if cmd.Err() != nil {
					log.Fatal(cmd.Err())
					continue
				}
			} else {
				// Save code and some more info in redis
				key := fmt.Sprintf("%d:%d", req.UserId, code)
				cmd := redis.Set(context.Background(), key, req.ActionId, time.Minute*10)
				if cmd.Err() != nil {
					log.Fatal(cmd.Err())
					continue
				}
			}

		}
		close(errChan)
	}()

	select {
	case err := <-errChan:
		return err
	case <-forever:
		return nil
	}
}
