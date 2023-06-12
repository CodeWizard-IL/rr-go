package rrbackend

import "log"

type LocalBackend struct {
	channels map[string]chan TransportEnvelope
}

func (t *LocalBackend) Connect() error {

	if t.channels != nil {
		log.Default().Println("Local backend is already connected")
		return nil
	}

	log.Default().Println("Connecting local backend")

	t.channels = make(map[string]chan TransportEnvelope)

	return nil
}

func (t *LocalBackend) getOrCreateChannelByID(ID string) chan TransportEnvelope {
	channel := t.channels[ID]

	if channel == nil {
		log.Default().Printf("Creating channel %s", ID)
		channel = make(chan TransportEnvelope, 100)
		t.channels[ID] = channel
	}

	return channel
}

func (t *LocalBackend) GetRequestReadChannelByID(ID string) <-chan TransportEnvelope {
	return t.getOrCreateChannelByID(ID)
}

func (t *LocalBackend) GetResponseReadChannelByID(ID string) <-chan TransportEnvelope {
	return t.getOrCreateChannelByID(ID)
}
func (t *LocalBackend) GetRequestWriteChannelByID(ID string) chan<- TransportEnvelope {
	return t.getOrCreateChannelByID(ID)
}

func (t *LocalBackend) GetResponseWriteChannelByID(ID string) chan<- TransportEnvelope {
	return t.getOrCreateChannelByID(ID)
}

func (t *LocalBackend) ReleaseChannelByID(ID string) error {
	log.Default().Printf("Releasing channel %s", ID)
	delete(t.channels, ID)
	return nil
}

func (t *LocalBackend) GetEnvelopeSerdes() EnvelopeSerdes {
	return &LocalTransportEnvelopeSerdes{}
}
