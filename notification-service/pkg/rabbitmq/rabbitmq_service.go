package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"micro-warehouse/notificaiton-service/configs"
	"micro-warehouse/notificaiton-service/pkg/email"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type RabbitMQServiceInterface interface {
	ConsumeEmail(ctx context.Context, emailService email.EmailServiceInterface) error
	Close() error
}

type rabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  configs.Config
}

// Close implements RabbitMQServiceInterface.
func (r *rabbitMQService) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// ConsumeEmail implements RabbitMQServiceInterface.
func (r *rabbitMQService) ConsumeEmail(ctx context.Context, emailService email.EmailServiceInterface) error {
	// Declare queue if not exists
	queue, err := r.channel.QueueDeclare(
		"email_queue", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Errorf("[RabbitMQService] ConsumeEmail - Queue declaration error: %v", err)
		return err
	}

	msgs, err := r.channel.Consume(
		queue.Name, // use declared queue name
		"",
		false, // auto-ack false, kita akan ack manual
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("[RabbitMQService] ConsumeEmail - 1: %v", err)
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("Email consumer context cancelled")
				return
			case msg := <-msgs:
				// Validasi pesan tidak kosong
				if len(msg.Body) == 0 {
					log.Warnf("[RabbitMQService] ConsumeEmail - Empty message received")
					msg.Nack(false, false)
					continue
				}

				var emailPayload email.EmailPayload
				if err := json.Unmarshal(msg.Body, &emailPayload); err != nil {
					log.Errorf("[RabbitMQService] ConsumeEmail - 2: JSON unmarshal error: %v, Raw message: %s", err, string(msg.Body))
					msg.Nack(false, false)
					continue
				}

				// Validasi payload
				if emailPayload.Email == "" {
					log.Errorf("[RabbitMQService] ConsumeEmail - Invalid payload: email is empty")
					msg.Nack(false, false)
					continue
				}

				// Process email berdasarkan type
				var err error
				switch emailPayload.Type {
				case "welcome", "welcome_email":
					err = emailService.SendWelcomeEmail(ctx, emailPayload)
				default:
					log.Errorf("[RabbitMQService] ConsumeEmail - 3: Unknown email type: %s", emailPayload.Type)
					msg.Nack(false, false)
					continue
				}

				if err != nil {
					log.Errorf("[RabbitMQService] ConsumeEmail - 4: %v", err)
					msg.Nack(false, true) // requeue
				} else {
					log.Infof("[RabbitMQService] ConsumeEmail - 5: Email sent successfully to %s", emailPayload.Email)
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

// dialWithRetry attempts to connect to RabbitMQ with exponential backoff.
// It retries up to maxRetries times before returning an error.
func dialWithRetry(url string) (*amqp.Connection, error) {
	const maxRetries = 10
	delay := 2 * time.Second
	maxDelay := 30 * time.Second

	for i := 1; i <= maxRetries; i++ {
		conn, err := amqp.Dial(url)
		if err == nil {
			log.Infof("[RabbitMQ] Connected successfully on attempt %d/%d", i, maxRetries)
			return conn, nil
		}
		log.Warnf("[RabbitMQ] Connection attempt %d/%d failed: %v. Retrying in %v...", i, maxRetries, err, delay)
		if i < maxRetries {
			time.Sleep(delay)
			delay *= 2
			if delay > maxDelay {
				delay = maxDelay
			}
		}
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts", maxRetries)
}

func NewRabbitMQService(config configs.Config) (RabbitMQServiceInterface, error) {
	conn, err := dialWithRetry(config.RabbitMQ.URL())
	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService - 1: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService - 2: %v", err)
		return nil, err
	}

	return &rabbitMQService{
		conn:    conn,
		channel: ch,
		config:  config,
	}, nil
}
