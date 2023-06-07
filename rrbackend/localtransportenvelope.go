package rrbackend

type LocalTransportEnvelopeSerdes struct {
}

func (serdes *LocalTransportEnvelopeSerdes) SerializeForRequest(envelope RREnvelope) (TransportEnvelope, error) {
	return envelope, nil
}

func (serdes *LocalTransportEnvelopeSerdes) SerializeForResponse(envelope RREnvelope) (TransportEnvelope, error) {
	return envelope, nil
}

func (serdes *LocalTransportEnvelopeSerdes) DeserializeForRequest(envelope TransportEnvelope) (RREnvelope, error) {
	return deserializeAsIs(envelope)
}

func (serdes *LocalTransportEnvelopeSerdes) DeserializeForResponse(envelope TransportEnvelope) (RREnvelope, error) {
	return deserializeAsIs(envelope)
}

func deserializeAsIs(envelope TransportEnvelope) (RREnvelope, error) {
	rrEnvelope, ok := envelope.(RREnvelope)
	if !ok {
		return RREnvelope{}, UnsupportedTransportEnvelopeError{Reason: "Not an RREnvelope"}
	}

	return rrEnvelope, nil
}
