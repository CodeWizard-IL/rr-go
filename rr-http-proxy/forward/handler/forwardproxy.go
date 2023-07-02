package handler

import (
	"bytes"
	"common/data"
	"encoding/json"
	"net/http"
	"rr-lib/rrbuilder"
	"rr-lib/rrclient"
)

type ProxyHandler struct {
	RRClient rrclient.RequestResponseClient
}

// Implements http.Handler
func (handler *ProxyHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {

	container := handler.convertRequest(request)

	// Marshall as json
	containerBytes, err2 := json.Marshal(container)
	if err2 != nil {
		panic(err2)
	}

	response, err := rrbuilder.NewRequest().WithPayload(containerBytes).Send(handler.RRClient)
	if err != nil {
		panic(err)
	}

	// Unmarshall response
	var responseContainer data.ResponseContainer
	err3 := json.Unmarshal(response.Payload, &responseContainer)
	if err3 != nil {
		panic(err3)
	}

	responseWriter.WriteHeader(responseContainer.StatusCode)
	for key, values := range responseContainer.Headers {
		for _, value := range values {
			responseWriter.Header().Add(key, value)
		}
	}

	responseWriter.Write(responseContainer.Body)

}

// Convert http.Request to RequestContainer
func (handler *ProxyHandler) convertRequest(request *http.Request) data.RequestContainer {

	buffer := bytes.Buffer{}
	buffer.ReadFrom(request.Body)

	container := data.RequestContainer{
		RawURI:    request.RequestURI,
		Host:      request.Host,
		ParsedURL: *request.URL,
		Method:    request.Method,
		Headers:   request.Header,
		Body:      buffer.Bytes(),
	}
	return container
}

func NewProxyHandler(rrClient rrclient.RequestResponseClient) *ProxyHandler {
	return &ProxyHandler{RRClient: rrClient}
}
