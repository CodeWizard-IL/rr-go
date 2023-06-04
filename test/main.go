package main

import (
	"backend"
	"client"
	"encoding/json"
	"fmt"
	uuid "github.com/google/uuid"
	"log"
	"server"
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

func (processor *TestProcessor) ProcessRequest(request backend.Request) (backend.Response, error) {
	if request.ContentType != "application/json" {
		return backend.Response{}, UnsupportedContentTypeError{}
	}

	payloadBytes := request.Payload

	var payload TestRequestPayload

	err := json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return backend.Response{}, err
	}

	content := payload.Content

	responsePayload := TestResponsePayload{
		Length: len(content),
	}

	responsePayloadBytes, _ := json.Marshal(responsePayload)

	response := backend.Response{
		ContentType: "application/json",
		Payload:     responsePayloadBytes,
	}

	return response, nil
}

func main() {
	fmt.Println("Starting RR tests")

	testBackend := backend.LocalBackend{}

	processor := TestProcessor{}

	rrServer := server.SimpleRequestResponseServer{
		Backend:   &testBackend,
		Processor: &processor,
	}

	err := rrServer.Start()
	if err != nil {
		log.Fatal(err)
	}

	request := backend.Request{
		ResponseId:  uuid.New().String(),
		ContentType: "application/json",
		Payload:     []byte(`{"content": "Hello world!"}`),
	}

	rrClient := client.SimpleRequestResponseClient{
		Backend: &testBackend,
	}

	response, err := rrClient.SendRequest(request)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s", response.Payload)
}
