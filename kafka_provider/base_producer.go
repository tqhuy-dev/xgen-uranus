package kafka_provider

type IProducer interface {
	SendMessages(topic string, key string, data []byte) error
}
