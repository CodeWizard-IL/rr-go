package rrclient

import (
	"github.com/google/uuid"
	"log"
	. "rr-lib/rrbackend"
	"time"
)

type ResponseTimeoutError struct {
}

func (e ResponseTimeoutError) Error() string {
	return "Response timeout"
}

type RequestResponseClient interface {
	Start() error
	SendRequestAsync(request RREnvelope) (ResponseHandler, error)
	SendRequest(request RREnvelope) (RREnvelope, error)
}

type ResponseHandler interface {
	ReceiveResponse() (RREnvelope, error)
	ReleaseResponseChannel()
}

type SimpleResponseHandler struct {
	ResponseChannelID string
	Backend           RequestResponseBackend
	ResponseChannel   <-chan TransportEnvelope
	TimeoutMillis     int
}

func (handler *SimpleResponseHandler) ReceiveResponse() (RREnvelope, error) {

	var timeout = make(<-chan time.Time)
	if handler.TimeoutMillis > 0 {
		log.Default().Printf("Setting timeout to %d ms", handler.TimeoutMillis)
		timeout = time.After(time.Duration(handler.TimeoutMillis) * time.Millisecond)
	}

	select {
	case response := <-handler.ResponseChannel:
		envelope, err := handler.Backend.GetEnvelopeSerdes().DeserializeForResponse(response)
		if err != nil {
			return RREnvelope{}, err
		}
		return envelope, nil
	case <-timeout:
		return RREnvelope{}, ResponseTimeoutError{}
	}

}

func (handler *SimpleResponseHandler) ReleaseResponseChannel() {
	err := handler.Backend.ReleaseChannelByID(handler.ResponseChannelID)
	if err != nil {
		return
	}
}

type SimpleRequestResponseClient struct {
	RequestChannelID string
	Backend          RequestResponseBackend
	TimeoutMillis    int
}

// Start the client
func (client *SimpleRequestResponseClient) Start() error {
	err := client.Backend.Connect()
	return err
}

func (client *SimpleRequestResponseClient) SendRequestAsync(request RREnvelope) (ResponseHandler, error) {
	if request.ID == "" {
		request.ID = NewUuid()
	}

	responseChannel, actualResponseID := client.Backend.GetResponseReadChannelByID(request.ID)
	if request.ID != actualResponseID {
		log.Printf("Request ID %s was replaced by %s by backend implementation", request.ID, actualResponseID)
		request.ID = actualResponseID
	}

	channel, _ := client.Backend.GetRequestWriteChannelByID(client.RequestChannelID)

	transportEnvelope, err := client.Backend.GetEnvelopeSerdes().SerializeForRequest(request)

	if err != nil {
		return nil, err
	}

	channel <- transportEnvelope

	handler := SimpleResponseHandler{
		ResponseChannelID: request.ID,
		Backend:           client.Backend,
		ResponseChannel:   responseChannel,
		TimeoutMillis:     client.TimeoutMillis,
	}

	return &handler, nil
}

func NewUuid() string {
	return uuid.New().String()
}

func (client *SimpleRequestResponseClient) SendRequest(request RREnvelope) (RREnvelope, error) {
	handler, err := client.SendRequestAsync(request)
	if err != nil {
		return RREnvelope{}, err
	}
	response, err := handler.ReceiveResponse()

	handler.ReleaseResponseChannel()

	return response, err
}
