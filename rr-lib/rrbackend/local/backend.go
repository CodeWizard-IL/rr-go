package local

import (
	"github.com/CodeWizard-IL/rr-go/rr-lib/rrbackend"
	"log"
)

type RRBackendLocal struct {
	channels map[string]chan rrbackend.TransportEnvelope
}

func (t *RRBackendLocal) Connect() error {

	if t.channels != nil {
		log.Default().Println("Local backend is already connected")
		return nil
	}

	log.Default().Println("Connecting local backend")

	t.channels = make(map[string]chan rrbackend.TransportEnvelope)

	return nil
}

func (t *RRBackendLocal) getOrCreateChannelByID(ID string) chan rrbackend.TransportEnvelope {
	channel := t.channels[ID]

	if channel == nil {
		log.Default().Printf("Creating channel %s", ID)
		channel = make(chan rrbackend.TransportEnvelope, 100)
		t.channels[ID] = channel
	}

	return channel
}

func (t *RRBackendLocal) GetRequestReadChannelByID(ID string) (<-chan rrbackend.TransportEnvelope, string) {
	return t.getOrCreateChannelByID(ID), ID
}

func (t *RRBackendLocal) GetResponseReadChannelByID(ID string) (<-chan rrbackend.TransportEnvelope, string) {
	return t.getOrCreateChannelByID(ID), ID
}
func (t *RRBackendLocal) GetRequestWriteChannelByID(ID string) (chan<- rrbackend.TransportEnvelope, string) {
	return t.getOrCreateChannelByID(ID), ID
}

func (t *RRBackendLocal) GetResponseWriteChannelByID(ID string) (chan<- rrbackend.TransportEnvelope, string) {
	return t.getOrCreateChannelByID(ID), ID
}

func (t *RRBackendLocal) ReleaseChannelByID(ID string) error {
	log.Default().Printf("Releasing channel %s", ID)
	delete(t.channels, ID)
	return nil
}

func (t *RRBackendLocal) GetEnvelopeSerdes() rrbackend.EnvelopeSerdes {
	return &TransportEnvelopeSerdes{}
}
