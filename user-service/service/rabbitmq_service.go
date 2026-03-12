package service

import (
	"context"
	"encoding/json"
	"fmt"
	"micro-warehouse/user-service/configs"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type EmailPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
}

type RabbitMQServiceInterface interface {
	PublishEmail(ctx context.Context, payload EmailPayload) error
	Close() error
}

type rabbitMQService struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	config configs.Config
}

// Close implements RabbitMQServiceInterface.
func (r *rabbitMQService) Close() error {
	if r.ch != nil {
		r.ch.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// PublishEmail implements RabbitMQServiceInterface.
func (r *rabbitMQService) PublishEmail(ctx context.Context, payload EmailPayload) error {
	// Convert payload ke JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %v", err)
	}

	// Declare queue if not exists
	queue, err := r.ch.QueueDeclare(
		"email_queue", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare email queue: %v", err)
	}

	// Publish ke email queue langsung (tanpa exchange)
	err = r.ch.Publish(
		"",         // exchange (empty for default)
		queue.Name, // routing key (queue name)
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish email message: %v", err)
	}

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
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.RabbitMQ.Username, config.RabbitMQ.Password, config.RabbitMQ.Host, config.RabbitMQ.Port)
	conn, err := dialWithRetry(url)
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
		conn:   conn,
		ch:     ch,
		config: config,
	}, nil
}
