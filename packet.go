package goeventbus

type packet struct {
	from    Sender
	message Message
}

func newSubscriberPacket(message Message) packet {
	return packet{from: SUBSCRIBER, message: message}
}

func newPublisherPacket(message Message) packet {
	return packet{from: PUBLISHER, message: message}
}
