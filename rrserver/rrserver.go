package rrserver

import (
	"log"
	. "rrbackend"
	"time"
)

type RequestProcessor interface {
	ProcessRequest(request Request) (Response, error)
}
type RequestResponseServer interface {
	Start() error
	Stop() error
}

type SimpleRequestResponseServer struct {
	Backend RequestResponseBackend

	Processor RequestProcessor

	requests chan Request

	running bool
}

func (server *SimpleRequestResponseServer) Start() error {
	err := server.Backend.Connect()
	if err != nil {
		return err
	}
	server.requests = server.Backend.GetRequestChannel().GetChannel()

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

			log.Printf("Received request: %s", request)

			response, err := server.Processor.ProcessRequest(request)

			if err != nil {
				log.Printf("Error processing request: %s", err)
				continue
			}

			responseChannel := server.Backend.GetResponseChannel(request).GetChannel()
			responseChannel <- response
		}

	}

}

func (server *SimpleRequestResponseServer) Stop() error {
	server.running = false

	return nil
}
