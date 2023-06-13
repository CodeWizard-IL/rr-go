package local

import "rrbackend"

type TransportEnvelopeSerdes struct {
}

func (serdes *TransportEnvelopeSerdes) SerializeForRequest(envelope rrbackend.RREnvelope) (rrbackend.TransportEnvelope, error) {
	return envelope, nil
}

func (serdes *TransportEnvelopeSerdes) SerializeForResponse(envelope rrbackend.RREnvelope) (rrbackend.TransportEnvelope, error) {
	return envelope, nil
}

func (serdes *TransportEnvelopeSerdes) DeserializeForRequest(envelope rrbackend.TransportEnvelope) (rrbackend.RREnvelope, error) {
	return deserializeAsIs(envelope)
}

func (serdes *TransportEnvelopeSerdes) DeserializeForResponse(envelope rrbackend.TransportEnvelope) (rrbackend.RREnvelope, error) {
	return deserializeAsIs(envelope)
}

func deserializeAsIs(envelope rrbackend.TransportEnvelope) (rrbackend.RREnvelope, error) {
	rrEnvelope, ok := envelope.(rrbackend.RREnvelope)
	if !ok {
		return rrbackend.RREnvelope{}, rrbackend.UnsupportedTransportEnvelopeError{Reason: "Not an RREnvelope"}
	}

	return rrEnvelope, nil
}
