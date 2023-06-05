package rrbackendamqp09

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	. "rrbackend"
)

type Amqp09Backend struct {
	ConnectString string

	amqpConnection *amqp.Connection
	amqpChannel    *amqp.Channel

	queues map[string]*amqp.Queue
}

func (backend *Amqp09Backend) Connect() error {
	// Connect to RabbitMQ server
	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(backend.ConnectString)
	if err != nil {
		return err
	}
	backend.amqpConnection = conn
	//defer backend.amqpConnection.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	backend.amqpChannel = ch
	//defer ch.Close()

	backend.queues = make(map[string]*amqp.Queue)

	return nil
}

func (backend *Amqp09Backend) getOrDeclareQueue(queueName string) *amqp.Queue {
	if _, ok := backend.queues[queueName]; !ok {
		q, err := backend.amqpChannel.QueueDeclare(
			queueName, // name
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			panic(err)
		}

		backend.queues[queueName] = &q
	}

	return backend.queues[queueName]
}

func (backend *Amqp09Backend) GetReadChannelByID(ID string) <-chan RREnvelope {
	q := backend.getOrDeclareQueue(ID)

	msgs, err := backend.amqpChannel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}

	envelops := make(chan RREnvelope)

	go func() {
		for msg := range msgs {
			var envelope RREnvelope
			err := json.Unmarshal(msg.Body, &envelope)
			if err != nil {
				log.Printf("Error decoding message: %s", err)
			}

			envelops <- envelope
		}
		close(envelops)
	}()

	return envelops
}

func (backend *Amqp09Backend) GetWriteChannelByID(ID string) chan<- RREnvelope {
	q := backend.getOrDeclareQueue(ID)

	envelops := make(chan RREnvelope)

	go func() {
		for envelope := range envelops {
			body, err := json.Marshal(envelope)
			if err != nil {
				panic(err)
			}

			err = backend.amqpChannel.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        body,
				})
			if err != nil {
				panic(err)
			}
		}
	}()

	return envelops
}

func (backend *Amqp09Backend) ReleaseChannelByID(ID string) error {
	queue := backend.queues[ID]
	if queue != nil {
		delete(backend.queues, ID)
	}
	return nil
}
