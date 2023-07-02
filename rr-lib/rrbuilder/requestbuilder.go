package rrbuilder

import . "github.com/CodeWizard-IL/rr-go/rr-lib/rrclient"
import "github.com/CodeWizard-IL/rr-go/rr-lib/rrbackend"

type RequestBuilder struct {
	envelope rrbackend.RREnvelope
}

func NewRequest() *RequestBuilder {
	return &RequestBuilder{
		envelope: rrbackend.RREnvelope{},
	}
}

func ReplyTo(request rrbackend.RREnvelope) *RequestBuilder {
	return &RequestBuilder{
		envelope: rrbackend.RREnvelope{
			ID: request.ID,
		},
	}
}

func (builder *RequestBuilder) WithID(id string) *RequestBuilder {
	builder.envelope.ID = id
	return builder
}

func (builder *RequestBuilder) WithContentType(contentType string) *RequestBuilder {
	builder.envelope.ContentType = contentType
	return builder
}

func (builder *RequestBuilder) WithHeaders(headers map[string]interface{}) *RequestBuilder {
	builder.envelope.Headers = headers
	return builder
}

func (builder *RequestBuilder) WithHeader(key string, value interface{}) *RequestBuilder {
	if builder.envelope.Headers == nil {
		builder.envelope.Headers = make(map[string]interface{})
	}
	builder.envelope.Headers[key] = value
	return builder
}

func (builder *RequestBuilder) WithPayload(payload []byte) *RequestBuilder {
	builder.envelope.Payload = payload
	return builder
}

func (builder *RequestBuilder) Build() rrbackend.RREnvelope {
	return builder.envelope
}

func (builder *RequestBuilder) SendAsync(client RequestResponseClient) (ResponseHandler, error) {
	return client.SendRequestAsync(builder.Build())
}

func (builder *RequestBuilder) Send(client RequestResponseClient) (rrbackend.RREnvelope, error) {
	return client.SendRequest(builder.Build())
}
