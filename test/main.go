package main

import (
	"encoding/json"
	"fmt"
	"log"
	"rrbackend"
	"rrclient"
	"rrserver"
)

type UnsupportedContentTypeError struct {
	ContentType string
}

func (e UnsupportedContentTypeError) Error() string {
	return fmt.Sprintf("Content type %s is not supported", e.ContentType)
}

type TestRequestPayload struct {
	Content string `json:"content"`
}

type TestResponsePayload struct {
	Length int `json:"length"`
}

type TestProcessor struct {
}

func (processor *TestProcessor) ProcessRequest(request rrbackend.RREnvelope) (rrbackend.RREnvelope, error) {
	if request.ContentType != "application/json" {
		return rrbackend.RREnvelope{}, UnsupportedContentTypeError{}
	}

	payloadBytes := request.Payload

	var payload TestRequestPayload

	err := json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return rrbackend.RREnvelope{}, err
	}

	content := payload.Content

	responsePayload := TestResponsePayload{
		Length: len(content),
	}

	responsePayloadBytes, _ := json.Marshal(responsePayload)

	response := rrbackend.RREnvelope{
		ID:          request.ID,
		ContentType: "application/json",
		Payload:     responsePayloadBytes,
	}

	return response, nil
}

func main() {
	fmt.Println("Starting RR tests")

	testBackend := rrbackend.LocalBackend{}

	processor := TestProcessor{}

	rrServer := rrserver.SimpleRequestResponseServer{
		RequestChannelID: "test-requests",
		Backend:          &testBackend,
		Processor:        &processor,
	}

	err := rrServer.Start()
	if err != nil {
		log.Fatal(err)
	}

	request := rrbackend.RREnvelope{
		ContentType: "application/json",
		Payload:     []byte(`{"content": "Hello world!"}`),
	}

	rrClient := rrclient.SimpleRequestResponseClient{
		RequestChannelID: "test-requests",
		Backend:          &testBackend,
		TimeoutMillis:    1000,
	}

	response, err := rrClient.SendRequest(request)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", response.Payload)

	secondResponse, err := rrClient.SendRequest(rrbackend.RREnvelope{
		ContentType: "application/json",
		Payload:     []byte(`{"content": "Goodbye world!"}`),
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Second response: %s\n", secondResponse.Payload)

}
