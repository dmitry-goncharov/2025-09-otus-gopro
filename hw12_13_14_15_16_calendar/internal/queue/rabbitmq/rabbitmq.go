package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/queue"
	"github.com/streadway/amqp"
)

const (
	exchangeName = ""
	contextType  = "application/json"
)

type RabbitMQ struct {
	log        app.Logger
	dsn        string
	queueName  string
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQ(log app.Logger, dsn, queueName string) queue.MessageQueue {
	return &RabbitMQ{
		log:       log,
		dsn:       dsn,
		queueName: queueName,
	}
}

func (q *RabbitMQ) Produce(_ context.Context, message app.Message) error {
	q.log.Debug("produce message to rabbitmq")
	if q.channel == nil {
		if err := q.connectToExchange(); err != nil {
			return fmt.Errorf("error getting connection to exchange: %w", err)
		}
	}
	body, err := encondeMessage(message)
	if err != nil {
		return fmt.Errorf("error encoding message: %w", err)
	}
	if err = q.channel.Publish(
		exchangeName, // publish to an exchange
		q.queueName,  // routing to 0 or more queues
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     contextType,
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("error publishing to exchange: %w", err)
	}
	return nil
}

func (q *RabbitMQ) Consume(ctx context.Context) (<-chan app.Message, error) {
	q.log.Debug("cosume message from rabbitmq")
	if q.channel == nil {
		if err := q.connectToExchange(); err != nil {
			return nil, fmt.Errorf("error getting connection to exchange: %w", err)
		}
	}
	out := make(chan app.Message)
	deliveries, err := q.channel.Consume(q.queueName, "", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("error consuming: %w", err)
	}
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case delivery := <-deliveries:
				msg, err := decodeMessage(delivery.Body)
				if err != nil {
					q.log.Error("error decoding message")
				} else {
					out <- *msg
				}
				delivery.Ack(false)
			}
		}
	}()
	return out, nil
}

func (q *RabbitMQ) Close() error {
	q.log.Debug("close rabbitmq")
	var chErr error
	if q.channel != nil {
		chErr = q.channel.Close()
	}
	var conErr error
	if q.connection != nil {
		conErr = q.connection.Close()
	}
	if chErr != nil || conErr != nil {
		return fmt.Errorf("error closing rabbitmq, chErr: %w, conErr: %w", chErr, conErr)
	}
	return nil
}

func (q *RabbitMQ) connectToExchange() error {
	if q.connection == nil {
		q.log.Debug("getting connection " + q.dsn)
		connection, err := amqp.Dial(q.dsn)
		if err != nil {
			return fmt.Errorf("error getting connection: %w", err)
		}
		q.connection = connection
	}
	if q.channel == nil {
		q.log.Debug("got connection, getting channel")
		channel, err := q.connection.Channel()
		if err != nil {
			return fmt.Errorf("error getting channel: %w", err)
		}
		q.log.Debug("got channel, declaring queue")
		_, err = channel.QueueDeclare(
			q.queueName, // name of the queue
			true,        // durable
			false,       // delete when unused
			false,       // exclusive
			false,       // noWait
			nil,         // arguments
		)
		if err != nil {
			return fmt.Errorf("error declaring queue: %w", err)
		}
		q.log.Debug("declared queue")
		q.channel = channel
	}
	return nil
}

func encondeMessage(message app.Message) ([]byte, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func decodeMessage(b []byte) (*app.Message, error) {
	var message app.Message
	err := json.Unmarshal(b, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}
