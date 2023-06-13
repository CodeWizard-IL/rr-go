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

// NewReverseProxyProcessor Creates a new ReverseProxyProcessor
func NewReverseProxyProcessor(urlMapper URLMapper) *ReverseProxyProcessor {
	return &ReverseProxyProcessor{UrlMapper: urlMapper}
}

// Process Processes the request and returns a response
func (processor *ReverseProxyProcessor) ProcessRequest(request rrbackend.RREnvelope) (rrbackend.RREnvelope, error) {

	var requestContainer data.RequestContainer
	json.Unmarshal(request.Payload, &requestContainer)

	// TODO: Implement error handling

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

	if err != nil {
		// TODO: Handle error
	}

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		// TODO: Handle error
	}

	responseContainer := data.ResponseContainer{
		StatusCode: response.StatusCode,
		Headers:    response.Header,
		Body:       bodyBytes,
	}

	responseBytes, err := json.Marshal(responseContainer)

	if err != nil {
		// TODO: Handle error
	}

	return rrbuilder.ReplyTo(request).WithPayload(responseBytes).Build(), nil
}
