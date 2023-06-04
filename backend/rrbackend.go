package backend

type RequestResponseBackend interface {
	Connect() error
	GetRequestChannel() RequestChannel
	GetResponseChannel(request Request) ResponseChannel
	ReleaseResponseChannel(response ResponseChannel)
}

type RequestChannel interface {
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
