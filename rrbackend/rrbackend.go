package rrbackend

type RREnvelope struct {
	ID          string
	ContentType string
	Headers     map[string]interface{}
	Payload     []byte
}

type UnsupportedTransportEnvelopeError struct {
	Reason string
}

func (e UnsupportedTransportEnvelopeError) Error() string {
	return "Unsupported transport envelope: " + e.Reason
}

type TransportEnvelope interface {
}
type EnvelopeSerdes interface {
	SerializeForRequest(envelope RREnvelope) (TransportEnvelope, error)
	SerializeForResponse(envelope RREnvelope) (TransportEnvelope, error)
	DeserializeForRequest(serialized TransportEnvelope) (RREnvelope, error)
	DeserializeForResponse(serialized TransportEnvelope) (RREnvelope, error)
}

type RequestResponseBackend interface {
	Connect() error
	GetRequestReadChannelByID(ID string) <-chan TransportEnvelope
	GetResponseReadChannelByID(ID string) <-chan TransportEnvelope
	GetRequestWriteChannelByID(ID string) chan<- TransportEnvelope
	GetResponseWriteChannelByID(ID string) chan<- TransportEnvelope
	ReleaseChannelByID(ID string) error
	GetEnvelopeSerdes() EnvelopeSerdes
}

type ClientRequestChannel interface {
	GetChannel() chan Request
}

type ResponseChannel interface {
	GetResponseId() string
	GetChannel() chan Response
}

type Request struct {
	ResponseId  string
	ContentType string
	Headers     map[string]interface{}
	Payload     []byte
}

type Response struct {
	ContentType string
	Headers     map[string]interface{}
	Payload     []byte
}
