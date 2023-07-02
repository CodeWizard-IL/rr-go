package main

import (
	"encoding/json"
	"fmt"
	"log"
	. "rr-lib/rrbackend"
	. "rr-lib/rrbackend/local"
	. "rr-lib/rrbuilder"
	. "rr-lib/rrclient"
	. "rr-lib/rrserver"
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

func (processor *TestProcessor) ProcessRequest(request RREnvelope) (RREnvelope, error) {
	if request.ContentType != "application/json" {
		return RREnvelope{}, UnsupportedContentTypeError{}
	}

	payloadBytes := request.Payload

	var payload TestRequestPayload

	err := json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return RREnvelope{}, err
	}

	content := payload.Content

	responsePayload := TestResponsePayload{
		Length: len(content),
	}

	responsePayloadBytes, _ := json.Marshal(responsePayload)

	response := ReplyTo(request).
		WithContentType("application/json").
		WithPayload(responsePayloadBytes).
		Build()

	return response, nil
}

func main() {
	fmt.Println("Starting RR tests")

	testBackend := RRBackendLocal{}

	//testBackend := Amqp09Backend{
	//	ConnectString: "amqp://guest:guest@localhost:5672/",
	//}

	//testBackend := RRBackendAzSMB{
	//	ConnectionString:  "Endpoint=sb://cwalexeyrr.servicebus.windows.net/;SharedAccessKeyName=rrgo;SharedAccessKey=sKMyUVlVxhjG62QrJh3mLlS/zXLpIK/a9+ASbLD88Xc=",
	//	RequestQueueName:  "myrequest",
	//	ResponseQueueName: "myrequest-response",
	//}

	processor := TestProcessor{}

	rrServer := SimpleRequestResponseServer{
		RequestChannelID: "myrequest",
		Backend:          &testBackend,
		Processor:        &processor,
	}

	err := rrServer.Start()
	if err != nil {
		log.Fatal(err)
	}

	rrClient := SimpleRequestResponseClient{
		RequestChannelID: "myrequest",
		Backend:          &testBackend,
		TimeoutMillis:    10000,
	}

	response, err := NewRequest().
		WithContentType("application/json").
		WithPayload([]byte(`{"content": "Hello world!"}`)).
		Send(&rrClient)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", response.Payload)

	secondResponse, err := rrClient.SendRequest(RREnvelope{
		ContentType: "application/json",
		Payload:     []byte(`{"content": "Goodbye world!"}`),
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Second response: %s\n", secondResponse.Payload)

}
