package backend

import "log"

type LocalRequestChannel struct {
	Channel chan Request
}

func (t LocalRequestChannel) GetChannel() chan Request {
	return t.Channel
}

type LocalResponseChannel struct {
	ResponseId string
	Channel    chan Response
}

func (t LocalResponseChannel) GetChannel() chan Response {
	return t.Channel
}

func (t LocalResponseChannel) GetResponseId() string {
	return t.ResponseId
}

type LocalBackend struct {
	requests  chan Request
	responses map[string]chan Response
}

func (test *LocalBackend) Connect() error {
	if test.requests != nil {
		log.Default().Println("Test backend is already connected")
		return nil
	}

	log.Default().Println("Connecting test backend")
	test.requests = make(chan Request, 100)
	test.responses = make(map[string]chan Response)

	return nil
}

func (test *LocalBackend) GetRequestChannel() RequestChannel {
	return LocalRequestChannel{
		Channel: test.requests,
	}
}

func (test *LocalBackend) GetResponseChannel(request Request) ResponseChannel {

	if test.responses[request.ResponseId] == nil {
		log.Printf("Creating response channel for request %s", request.ResponseId)
		test.responses[request.ResponseId] = make(chan Response, 1)
	} else {
		log.Printf("Reusing response channel for request %s", request.ResponseId)
	}

	responseChannel := test.responses[request.ResponseId]

	return LocalResponseChannel{
		ResponseId: request.ResponseId,
		Channel:    responseChannel,
	}
}

func (test *LocalBackend) ReleaseResponseChannel(response ResponseChannel) {
	responseChannel := response.GetChannel()
	responseId := response.GetResponseId()

	log.Printf("Releasing response channel for request %s", responseId)

	delete(test.responses, responseId)
	close(responseChannel)
}
