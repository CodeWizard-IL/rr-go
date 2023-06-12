package rrbackendazsmb

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"log"
	. "rrbackend"
)

type RRBackendAzSMB struct {
	ConnectionString  string
	RequestQueueName  string
	ResponseQueueName string

	client               *azservicebus.Client
	requestQueueSender   *azservicebus.Sender
	responseQueueSender  *azservicebus.Sender
	requestQueueReceiver *azservicebus.Receiver
}

func (backend *RRBackendAzSMB) Connect() error {
	var err error

	backend.client, err = azservicebus.NewClientFromConnectionString(backend.ConnectionString, nil)
	if err != nil {
		return err
	}

	backend.requestQueueSender, err = backend.client.NewSender(backend.RequestQueueName, nil)
	if err != nil {
		return err
	}

	backend.responseQueueSender, err = backend.client.NewSender(backend.ResponseQueueName, nil)
	if err != nil {
		return err
	}

	return err
}

func (backend *RRBackendAzSMB) getOrCreateRequestReceiverForServer() (*azservicebus.Receiver, error) {
	var err error
	if backend.requestQueueReceiver == nil {
		backend.requestQueueReceiver, err = backend.client.NewReceiverForQueue(backend.RequestQueueName, nil)
		if err != nil {
			return &azservicebus.Receiver{}, err
		}
	}

	return backend.requestQueueReceiver, nil
}

// GetRequestReadChannelByID returns a channel that can be used to read requests from the specified ID
func (backend *RRBackendAzSMB) GetRequestReadChannelByID(ID string) <-chan TransportEnvelope {
	receiver, err := backend.getOrCreateRequestReceiverForServer()

	if err != nil {
		return nil
	}

	channel := make(chan TransportEnvelope)

	go func() {
		for {
			msg, err := receiver.ReceiveMessages(context.TODO(), 1, nil)
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				close(channel)
				return
			}

			message := *msg[0]

			err = receiver.CompleteMessage(context.TODO(), &message, nil)
			if err != nil {
				log.Printf("Error completing message: %v", err)
				close(channel)
				return
			}

			channel <- message
		}
	}()

	return channel
}

// GetResponseReadChannelByID returns a channel that can be used to read responses from the specified ID
func (backend *RRBackendAzSMB) GetResponseReadChannelByID(ID string) <-chan TransportEnvelope {

	channel := make(chan TransportEnvelope)

	go func() {
		receiver, err := backend.client.AcceptSessionForQueue(context.TODO(), backend.ResponseQueueName, ID, nil)
		if err != nil {
			log.Printf("Error accepting session: %v", err)
		}
		defer receiver.Close(context.TODO())

		for {
			msg, err := receiver.ReceiveMessages(context.TODO(), 1, nil)
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				close(channel)
				return
			}

			message := *msg[0]

			err = receiver.CompleteMessage(context.TODO(), &message, nil)
			if err != nil {
				log.Printf("Error completing message: %v", err)
				close(channel)
				return
			}

			channel <- message
			// Single message responses are supported currently. Close the channel after the first message is received.
			close(channel)
			return
		}
	}()

	return channel
}

// GetRequestWriteChannelByID returns a channel that can be used to write requests to the specified ID
func (backend *RRBackendAzSMB) GetRequestWriteChannelByID(ID string) chan<- TransportEnvelope {
	channel := make(chan TransportEnvelope)

	go func() {
		for envelope := range channel {
			message, ok := envelope.(azservicebus.Message)
			if !ok {
				log.Printf("Error converting envelope to message")
				continue
			}
			if message.ReplyToSessionID == nil {
				log.Printf("Error: message.ReplyToSessionID is nil")
				continue
			}
			err := backend.requestQueueSender.SendMessage(context.TODO(), &message, nil)
			if err != nil {
				log.Printf("Error sending message: %v", err)
				continue
			}
		}
	}()

	return channel
}

// GetResponseWriteChannelByID returns a channel that can be used to write responses to the specified ID
func (backend *RRBackendAzSMB) GetResponseWriteChannelByID(ID string) chan<- TransportEnvelope {
	channel := make(chan TransportEnvelope)

	go func() {
		for envelope := range channel {
			message, ok := envelope.(azservicebus.Message)
			if !ok {
				log.Printf("Error converting envelope to message")
				continue
			}
			if message.SessionID == nil {
				message.SessionID = &ID
			} else if *message.SessionID != ID {
				log.Printf("GetResponseWriteChannelByID: Error: message.SessionID[%s] != ID[%s]", *message.SessionID, ID)
				continue
			}
			err := backend.responseQueueSender.SendMessage(context.TODO(), &message, nil)
			if err != nil {
				log.Printf("Error sending message: %v", err)
				continue
			}
		}
	}()

	return channel
}

func (backend *RRBackendAzSMB) ReleaseChannelByID(ID string) error {
	// Nothing to do here
	return nil
}

func (backend *RRBackendAzSMB) GetEnvelopeSerdes() EnvelopeSerdes {
	return &AzSMBTransportEnvelopeSerses{}
}
