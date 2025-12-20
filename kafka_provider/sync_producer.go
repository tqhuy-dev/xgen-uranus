package kafka_provider

import "github.com/IBM/sarama"

type SyncProducer struct {
	sarama.SyncProducer
}

func NewSyncProducer(config *KafkaConfig) (*SyncProducer, error) {
	config.Producer.Return.Errors = true
	producer, err := sarama.NewSyncProducer(config.GetBrokers(), config.ParseSaramaConfig())
	if err != nil {
		return nil, err
	}
	return &SyncProducer{producer}, nil
}

func (p *SyncProducer) Run() {
}

func (p *SyncProducer) SendMessages(topic string, key string, data []byte) (int32, int64, error) {
	partition, offset, err := p.SyncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	})
	return partition, offset, err
}
