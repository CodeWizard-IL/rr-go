package rrserver

import (
	"log"
	. "rr-lib/rrbackend"
	"time"
)

type RequestProcessor interface {
	ProcessRequest(request RREnvelope) (RREnvelope, error)
}
type RequestResponseServer interface {
	Start() error
	Stop() error
}

type SimpleRequestResponseServer struct {
	RequestChannelID string
	Backend          RequestResponseBackend
	Processor        RequestProcessor
	requests         <-chan TransportEnvelope
	running          bool
}

func (server *SimpleRequestResponseServer) Start() error {
	err := server.Backend.Connect()
	if err != nil {
		return err
	}
	server.requests, _ = server.Backend.GetRequestReadChannelByID(server.RequestChannelID)

	server.running = true

	go server.listenForRequests()

	return nil
}

func (server *SimpleRequestResponseServer) listenForRequests() {

	log.Println("Listening for requests...")

	ticker := time.NewTicker(1000 * time.Millisecond)

	done := make(chan bool)

	for {
		select {
		case <-done:
			log.Println("Stopped listening for requests")
			return
		case <-ticker.C:
			log.Println("Tick")
			if server.running == false {
				done <- true
			}
		case request := <-server.requests:

			envelopeSerdes := server.Backend.GetEnvelopeSerdes()
			rrEnvelope, err := envelopeSerdes.DeserializeForRequest(request)
			if err != nil {
				log.Printf("Error processing request: %s", err)
				continue
			}

			log.Printf("Received request: %s", rrEnvelope)

			response, err := server.Processor.ProcessRequest(rrEnvelope)

			if err != nil {
				log.Printf("Error processing request: %s", err)
				continue
			}

			responseChannel, _ := server.Backend.GetResponseWriteChannelByID(rrEnvelope.ID)

			transportEnvelope, err := envelopeSerdes.SerializeForResponse(response)

			if err != nil {
				log.Printf("Error serializing response: %s", err)
				continue
			}

			responseChannel <- transportEnvelope
		}

	}

}

func (server *SimpleRequestResponseServer) Stop() error {
	server.running = false

	return nil
}
