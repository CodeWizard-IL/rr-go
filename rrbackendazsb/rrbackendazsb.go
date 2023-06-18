package rrbackendazsb

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/google/uuid"
	"log"
	. "rrbackend"
)

type SessionReceiver struct {
	SessionID string
	Receiver  *azservicebus.SessionReceiver
}

type AzSBBackend struct {
	ConnectionString    string
	RequestQueueName    string
	ResponseQueueName   string
	MinSessionReceivers int
	MaxSessionReceivers int

	client                 *azservicebus.Client
	requestQueueSender     *azservicebus.Sender
	responseQueueSender    *azservicebus.Sender
	requestQueueReceiver   *azservicebus.Receiver
	responseQueueReceivers []SessionReceiver
	sessionReceiversMap    map[string]SessionReceiver

	obtainReceiver  chan SessionReceiver
	releaseReceiver chan SessionReceiver
}

func (backend *AzSBBackend) getObtainReceiverChannel() <-chan SessionReceiver {
	return backend.obtainReceiver
}

func (backend *AzSBBackend) getReleaseReceiverChannel() chan<- SessionReceiver {
	return backend.releaseReceiver
}

func (backend *AzSBBackend) newSessionReceiver() SessionReceiver {
	sessionID := uuid.New().String()
	receiver, err := backend.client.AcceptSessionForQueue(context.TODO(), backend.ResponseQueueName, sessionID, nil)
	if err != nil {
		log.Printf("Error creating new session receiver: %v", err)
		return SessionReceiver{}
	}

	sessionReceiver := SessionReceiver{SessionID: sessionID, Receiver: receiver}
	backend.sessionReceiversMap[sessionID] = sessionReceiver
	return sessionReceiver
}
func (backend *AzSBBackend) Connect() error {

	if backend.client != nil {
		log.Println("Already connected")
		return nil
	}

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

	if backend.MinSessionReceivers == 0 {
		backend.MinSessionReceivers = 10
	}

	if backend.MaxSessionReceivers == 0 {
		backend.MaxSessionReceivers = 20
	}

	backend.obtainReceiver = make(chan SessionReceiver)
	backend.releaseReceiver = make(chan SessionReceiver)
	backend.responseQueueReceivers = make([]SessionReceiver, backend.MinSessionReceivers)
	backend.sessionReceiversMap = make(map[string]SessionReceiver)

	log.Println("Creating initial session receivers")
	for i := 0; i < backend.MinSessionReceivers; i++ {
		backend.responseQueueReceivers[i] = backend.newSessionReceiver()
		//backend.responseQueueReceivers = append(backend.responseQueueReceivers, backend.newSessionReceiver())
		log.Printf("Created session receiver for session %s", backend.responseQueueReceivers[i].SessionID)
	}

	go func() {
		for {
			releasedReceiver := <-backend.releaseReceiver
			log.Printf("Releasing receiver for session %s", releasedReceiver.SessionID)
			backend.responseQueueReceivers = append(backend.responseQueueReceivers, releasedReceiver)
		}
	}()

	go func() {
		for {
			if len(backend.responseQueueReceivers) < 1 {
				backend.responseQueueReceivers = append(backend.responseQueueReceivers, backend.newSessionReceiver())
				log.Printf("Created new session receiver for session %s", backend.responseQueueReceivers[0].SessionID)
			}
			nextReceiver := backend.responseQueueReceivers[0]
			backend.responseQueueReceivers = backend.responseQueueReceivers[1:]
			log.Printf("Preparing next receiver for session %s", nextReceiver.SessionID)
			backend.obtainReceiver <- nextReceiver
			log.Printf("Obtained next receiver for session %s", nextReceiver.SessionID)
		}
	}()

	return err
}

func (backend *AzSBBackend) getOrCreateRequestReceiverForServer() (*azservicebus.Receiver, error) {
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
func (backend *AzSBBackend) GetRequestReadChannelByID(ID string) (<-chan TransportEnvelope, string) {
	receiver, err := backend.getOrCreateRequestReceiverForServer()

	if err != nil {
		return nil, ""
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

	return channel, ID
}

// GetResponseReadChannelByID returns a channel that can be used to read responses from the specified ID
func (backend *AzSBBackend) GetResponseReadChannelByID(ID string) (<-chan TransportEnvelope, string) {

	channel := make(chan TransportEnvelope)
	sessionID := ID
	syncForId := make(chan string)

	go func() {
		preparedSessionReceiver := <-backend.getObtainReceiverChannel()

		receiver := preparedSessionReceiver.Receiver
		sessionID = preparedSessionReceiver.SessionID

		syncForId <- sessionID

		for {
			msg, err := receiver.ReceiveMessages(context.TODO(), 1, nil)
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				close(channel)
				return
			}

			if len(msg) == 0 {
				log.Printf("No message received. Retrying...")
				continue
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

	sessionID = <-syncForId

	return channel, sessionID
}

// GetRequestWriteChannelByID returns a channel that can be used to write requests to the specified ID
func (backend *AzSBBackend) GetRequestWriteChannelByID(ID string) (chan<- TransportEnvelope, string) {
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

	return channel, ID
}

// GetResponseWriteChannelByID returns a channel that can be used to write responses to the specified ID
func (backend *AzSBBackend) GetResponseWriteChannelByID(ID string) (chan<- TransportEnvelope, string) {
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

	return channel, ID
}

func (backend *AzSBBackend) ReleaseChannelByID(ID string) error {

	sessionReceiver := backend.sessionReceiversMap[ID]

	backend.releaseReceiver <- sessionReceiver
	return nil
}

func (backend *AzSBBackend) GetEnvelopeSerdes() EnvelopeSerdes {
	return &AzSBTransportEnvelopeSerdes{}
}
