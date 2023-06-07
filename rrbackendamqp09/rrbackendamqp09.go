package rrbackendamqp09

import (
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

func (backend *Amqp09Backend) GetReadChannelByID(ID string) <-chan TransportEnvelope {
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

	trc := make(chan TransportEnvelope)

	go func() {
		for d := range msgs {
			trc <- TransportEnvelope(d)
		}
	}()

	return trc
}

func (backend *Amqp09Backend) GetWriteChannelByID(ID string) chan<- TransportEnvelope {

	q := backend.getOrDeclareQueue(ID)

	envelops := make(chan TransportEnvelope)

	go func() {
		for envelope := range envelops {
			publishing, ok := envelope.(amqp.Publishing)
			if !ok {
				log.Default().Printf("Error publishing to queue %s: %s", q.Name, "envelope is not of type amqp.Publishing")
				continue
			}
			err := backend.amqpChannel.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				publishing,
			)
			if err != nil {
				log.Default().Printf("Error publishing to queue %s: %s", q.Name, err)
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

func (backend *Amqp09Backend) GetEnvelopeSerdes() EnvelopeSerdes {
	return &StreadwayEnvelopeSerdes{}
}
