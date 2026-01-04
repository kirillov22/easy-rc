package messages

type Processable interface {
	Process() (returnPayload *Packet, err error)
	Debug()
}
