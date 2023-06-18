package processor

import (
	"bytes"
	"common/data"
	"encoding/json"
	"io"
	"net/http"
	"rrbackend"
	"rrbuilder"
	"time"
)

// ReverseProxyProcessor Implements the RequestResponseProcessor interface
type ReverseProxyProcessor struct {
	UrlMapper URLMapper
}

// ProcessRequest Processes the request and returns a response
func (processor *ReverseProxyProcessor) ProcessRequest(request rrbackend.RREnvelope) (rrbackend.RREnvelope, error) {

	var requestContainer data.RequestContainer
	err := json.Unmarshal(request.Payload, &requestContainer)
	if err != nil {
		return rrbackend.RREnvelope{}, err
	}

	forwardUrl := processor.UrlMapper.MapURL(requestContainer.Host, requestContainer.RawURI)

	newRequest, err := http.NewRequest(
		requestContainer.Method,
		forwardUrl,
		bytes.NewReader(requestContainer.Body),
	)

	if err != nil {
		//TODO: Handle error
	}

	for key, values := range requestContainer.Headers {
		for _, value := range values {
			newRequest.Header.Add(key, value)
		}
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	response, err := client.Do(newRequest)

	var responseContainer data.ResponseContainer

	if err != nil {
		responseContainer = data.ResponseContainer{
			StatusCode: http.StatusBadGateway,
			Headers:    http.Header{},
			Body:       []byte(err.Error()),
		}
	} else {

		bodyBytes, err := io.ReadAll(response.Body)

		if err != nil {
			// TODO: Handle error
		}

		responseContainer = data.ResponseContainer{
			StatusCode: response.StatusCode,
			Headers:    response.Header,
			Body:       bodyBytes,
		}
	}

	responseBytes, err := json.Marshal(responseContainer)
	if err != nil {
		// TODO: Handle error
	}

	return rrbuilder.ReplyTo(request).WithPayload(responseBytes).Build(), nil
}
