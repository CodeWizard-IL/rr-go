package rrclient

import (
	"github.com/google/uuid"
	"log"
	. "rrbackend"
	"time"
)

type ResponseTimeoutError struct {
}

func (e ResponseTimeoutError) Error() string {
	return "Response timeout"
}

type RequestResponseClient interface {
	SendRequestAsync(request Request) (ResponseHandler, error)
	SendRequest(request Request) (Response, error)
}

type ResponseHandler interface {
	ReceiveResponse() (Response, error)
	ReleaseResponseChannel(rrBackend RequestResponseBackend)
}

type SimpleResponseHandler struct {
	ResponseChannel ResponseChannel
	TimeoutMillis   int
}

func (handler *SimpleResponseHandler) ReceiveResponse() (Response, error) {

	var timeout = make(<-chan time.Time)
	if handler.TimeoutMillis > 0 {
		log.Default().Printf("Setting timeout to %d ms", handler.TimeoutMillis)
		timeout = time.After(time.Duration(handler.TimeoutMillis) * time.Millisecond)
	}

	select {
	case response := <-handler.ResponseChannel.GetChannel():
		return response, nil
	case <-timeout:
		return Response{}, ResponseTimeoutError{}
	}

}

func (handler *SimpleResponseHandler) ReleaseResponseChannel(rrBackend RequestResponseBackend) {
	rrBackend.ReleaseResponseChannel(handler.ResponseChannel)
}

type SimpleRequestResponseClient struct {
	Backend       RequestResponseBackend
	TimeoutMillis int
}

func (client *SimpleRequestResponseClient) SendRequestAsync(request Request) (ResponseHandler, error) {
	if request.ResponseId == "" {
		request.ResponseId = NewUuid()
	}

	responseChannel := client.Backend.GetResponseChannel(request)
	channel := client.Backend.GetRequestChannel().GetChannel()

	channel <- request

	handler := SimpleResponseHandler{
		ResponseChannel: responseChannel,
		TimeoutMillis:   client.TimeoutMillis,
	}

	return &handler, nil
}

func NewUuid() string {
	return uuid.New().String()
}

func (client *SimpleRequestResponseClient) SendRequest(request Request) (Response, error) {
	handler, err := client.SendRequestAsync(request)
	if err != nil {
		return Response{}, err
	}
	response, err := handler.ReceiveResponse()

	handler.ReleaseResponseChannel(client.Backend)

	return response, err
}
