package local

import (
	"log"
	"rrbackend"
)

type RequestResponseBackend struct {
	channels map[string]chan rrbackend.TransportEnvelope
}

func (t *RequestResponseBackend) Connect() error {

	if t.channels != nil {
		log.Default().Println("Local backend is already connected")
		return nil
	}

	log.Default().Println("Connecting local backend")

	t.channels = make(map[string]chan rrbackend.TransportEnvelope)

	return nil
}

func (t *RequestResponseBackend) getOrCreateChannelByID(ID string) chan rrbackend.TransportEnvelope {
	channel := t.channels[ID]

	if channel == nil {
		log.Default().Printf("Creating channel %s", ID)
		channel = make(chan rrbackend.TransportEnvelope, 100)
		t.channels[ID] = channel
	}

	return channel
}

func (t *RequestResponseBackend) GetRequestReadChannelByID(ID string) <-chan rrbackend.TransportEnvelope {
	return t.getOrCreateChannelByID(ID)
}

func (t *RequestResponseBackend) GetResponseReadChannelByID(ID string) <-chan rrbackend.TransportEnvelope {
	return t.getOrCreateChannelByID(ID)
}
func (t *RequestResponseBackend) GetRequestWriteChannelByID(ID string) chan<- rrbackend.TransportEnvelope {
	return t.getOrCreateChannelByID(ID)
}

func (t *RequestResponseBackend) GetResponseWriteChannelByID(ID string) chan<- rrbackend.TransportEnvelope {
	return t.getOrCreateChannelByID(ID)
}

func (t *RequestResponseBackend) ReleaseChannelByID(ID string) error {
	log.Default().Printf("Releasing channel %s", ID)
	delete(t.channels, ID)
	return nil
}

func (t *RequestResponseBackend) GetEnvelopeSerdes() rrbackend.EnvelopeSerdes {
	return &TransportEnvelopeSerdes{}
}
