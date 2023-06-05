package rrbackend

type RREnvelope struct {
	ID          string
	ContentType string
	Headers     map[string]interface{}
	Payload     []byte
}

type RequestResponseBackend interface {
	Connect() error
	//GetRequestChannel() ClientRequestChannel
	//GetResponseChannel(request Request) ResponseChannel
	//ReleaseResponseChannel(response ResponseChannel)
	GetReadChannelByID(ID string) <-chan RREnvelope
	GetWriteChannelByID(ID string) chan<- RREnvelope
	ReleaseChannelByID(ID string) error
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
