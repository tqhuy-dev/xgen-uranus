package kafka_provider

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

type KafkaConfig struct {
	Brokers  []string
	Producer struct {
		Return struct {
			Successes bool
			Errors    bool
		}
		RequiredAcks sarama.RequiredAcks
		Retry        struct {
			Max int
		}
		Flush struct {
			Frequency time.Duration
			Bytes     int
			Messages  int
		}
		Compression sarama.CompressionCodec
	}
}

func (c *KafkaConfig) GetBrokers() []string {
	if c == nil {
		return nil
	}
	return c.Brokers
}

func (c *KafkaConfig) ParseSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = c.Producer.Return.Successes
	config.Producer.Return.Errors = c.Producer.Return.Errors

	config.Producer.RequiredAcks = c.Producer.RequiredAcks
	config.Producer.Retry.Max = c.Producer.Retry.Max

	config.Producer.Flush.Frequency = c.Producer.Flush.Frequency
	config.Producer.Flush.Bytes = c.Producer.Flush.Bytes
	config.Producer.Flush.Messages = c.Producer.Flush.Messages

	config.Producer.Compression = c.Producer.Compression
	return config
}

type AsyncProducer struct {
	sarama.AsyncProducer
}

func NewAsyncProducer(config *KafkaConfig) (*AsyncProducer, error) {
	producer, err := sarama.NewAsyncProducer(config.GetBrokers(), config.ParseSaramaConfig())
	if err != nil {
		return nil, err
	}
	return &AsyncProducer{producer}, nil
}

func (p *AsyncProducer) Run() {
	go func() {
		for err := range p.Errors() {
			log.Printf("❌ produce error: %v", err)
		}
	}()
}

func (p *AsyncProducer) SendMessages(topic string, key string, data []byte) error {

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	}
	p.Input() <- msg

	return nil
}
