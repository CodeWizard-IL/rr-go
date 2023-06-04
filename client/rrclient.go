package client

import (
	"backend"
	"time"
)

type ResponseTimeoutError struct {
}

func (e ResponseTimeoutError) Error() string {
	return "Response timeout"
}

type RequestResponseClient interface {
	SendRequestAsync(request backend.Request) (ResponseHandler, error)
	SendRequest(request backend.Request) (backend.Response, error)
}

type ResponseHandler interface {
	ReceiveResponse() (backend.Response, error)
	ReleaseResponseChannel(rrBackend backend.RequestResponseBackend)
}

type SimpleResponseHandler struct {
	ResponseChannel backend.ResponseChannel
	TimeoutMillis   int
}

func (handler *SimpleResponseHandler) ReceiveResponse() (backend.Response, error) {
	timeout := time.After(time.Duration(handler.TimeoutMillis) * time.Millisecond)

	select {
	case response := <-handler.ResponseChannel.GetChannel():
		return response, nil
	case <-timeout:
		return backend.Response{}, ResponseTimeoutError{}
	}

}

func (handler *SimpleResponseHandler) ReleaseResponseChannel(rrBackend backend.RequestResponseBackend) {
	rrBackend.ReleaseResponseChannel(handler.ResponseChannel)
}

type SimpleRequestResponseClient struct {
	Backend backend.RequestResponseBackend
}

func (client *SimpleRequestResponseClient) SendRequestAsync(request backend.Request) (ResponseHandler, error) {
	responseChannel := client.Backend.GetResponseChannel(request)
	channel := client.Backend.GetRequestChannel().GetChannel()

	channel <- request

	handler := SimpleResponseHandler{
		ResponseChannel: responseChannel,
		TimeoutMillis:   1000000,
	}

	return &handler, nil
}

func (client *SimpleRequestResponseClient) SendRequest(request backend.Request) (backend.Response, error) {
	handler, err := client.SendRequestAsync(request)
	if err != nil {
		return backend.Response{}, err
	}
	response, err := handler.ReceiveResponse()

	handler.ReleaseResponseChannel(client.Backend)

	return response, err
}
