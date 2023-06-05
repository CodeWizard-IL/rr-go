package rrbackend

import "log"

type LocalBackend struct {
	channels map[string]chan RREnvelope
}

func (t *LocalBackend) Connect() error {

	if t.channels != nil {
		log.Default().Println("Local backend is already connected")
		return nil
	}

	log.Default().Println("Connecting local backend")

	t.channels = make(map[string]chan RREnvelope)

	return nil
}

func (t *LocalBackend) getOrCreateChannelByID(ID string) chan RREnvelope {
	channel := t.channels[ID]

	if channel == nil {
		log.Default().Printf("Creating channel %s", ID)
		channel = make(chan RREnvelope, 100)
		t.channels[ID] = channel
	}

	return channel
}

func (t *LocalBackend) GetReadChannelByID(ID string) <-chan RREnvelope {
	return t.getOrCreateChannelByID(ID)
}
func (t *LocalBackend) GetWriteChannelByID(ID string) chan<- RREnvelope {
	return t.getOrCreateChannelByID(ID)
}

func (t *LocalBackend) ReleaseChannelByID(ID string) error {
	log.Default().Printf("Releasing channel %s", ID)
	delete(t.channels, ID)
	return nil
}
