package amqp09

import (
	"errors"
	"github.com/streadway/amqp"
	. "rr-lib/rrbackend"
)

type StreadwayEnvelopeSerdes struct {
}

func (serdes *StreadwayEnvelopeSerdes) SerializeForRequest(request RREnvelope) (TransportEnvelope, error) {

	return amqp.Publishing{
		Headers:       request.Headers,
		ContentType:   request.ContentType,
		CorrelationId: request.ID,
		ReplyTo:       request.ID,
		Body:          request.Payload,
	}, nil
}

func (serdes *StreadwayEnvelopeSerdes) DeserializeForRequest(envelope TransportEnvelope) (RREnvelope, error) {

	delivery, ok := envelope.(amqp.Delivery)

	if !ok {
		return RREnvelope{}, errors.New("envelope is not a streadway amqp.Delivery")
	}

	return RREnvelope{
		ID:          delivery.ReplyTo,
		ContentType: delivery.ContentType,
		Headers:     delivery.Headers,
		Payload:     delivery.Body,
	}, nil

}

func (serdes *StreadwayEnvelopeSerdes) SerializeForResponse(response RREnvelope) (TransportEnvelope, error) {
	return amqp.Publishing{
		Headers:       response.Headers,
		ContentType:   response.ContentType,
		CorrelationId: response.ID,
		Body:          response.Payload,
	}, nil
}

func (serdes *StreadwayEnvelopeSerdes) DeserializeForResponse(envelope TransportEnvelope) (RREnvelope, error) {
	delivery, ok := envelope.(amqp.Delivery)

	if !ok {
		return RREnvelope{}, errors.New("envelope is not a streadway amqp.Delivery")
	}
	return RREnvelope{
		ID:          delivery.CorrelationId,
		ContentType: delivery.ContentType,
		Headers:     delivery.Headers,
		Payload:     delivery.Body,
	}, nil
}
