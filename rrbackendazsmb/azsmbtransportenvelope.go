package rrbackendazsmb

import (
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	. "rrbackend"
)

// Implements EnvelopeSerdes
type AzSMBTransportEnvelopeSerses struct {
}

// SerializeForRequest serializes the specified envelope for sending as a request
func (serdes *AzSMBTransportEnvelopeSerses) SerializeForRequest(envelope RREnvelope) (TransportEnvelope, error) {
	msg := azservicebus.Message{
		Body:                  envelope.Payload,
		ContentType:           &envelope.ContentType,
		ReplyToSessionID:      &envelope.ID,
		ApplicationProperties: envelope.Headers,
	}
	return msg, nil
}

// SerializeForResponse serializes the specified envelope for sending as a response
func (serdes *AzSMBTransportEnvelopeSerses) SerializeForResponse(envelope RREnvelope) (TransportEnvelope, error) {
	msg := azservicebus.Message{
		Body:                  envelope.Payload,
		ContentType:           &envelope.ContentType,
		SessionID:             &envelope.ID,
		ApplicationProperties: envelope.Headers,
	}
	return msg, nil
}

// DeserializeForRequest deserializes the specified envelope as a request
func (serdes *AzSMBTransportEnvelopeSerses) DeserializeForRequest(envelope TransportEnvelope) (RREnvelope, error) {
	msg, ok := envelope.(azservicebus.ReceivedMessage)
	if !ok {
		return RREnvelope{}, UnsupportedTransportEnvelopeError{Reason: "Not an azservicebus.Message"}
	}

	return RREnvelope{
		ID:          *msg.ReplyToSessionID,
		Payload:     msg.Body,
		ContentType: *msg.ContentType,
		Headers:     msg.ApplicationProperties,
	}, nil
}

// DeserializeForResponse deserializes the specified envelope as a response
func (serdes *AzSMBTransportEnvelopeSerses) DeserializeForResponse(envelope TransportEnvelope) (RREnvelope, error) {
	msg, ok := envelope.(azservicebus.ReceivedMessage)
	if !ok {
		return RREnvelope{}, UnsupportedTransportEnvelopeError{Reason: "Not an azservicebus.Message"}
	}

	return RREnvelope{
		ID:          *msg.SessionID,
		Payload:     msg.Body,
		ContentType: *msg.ContentType,
	}, nil
}
