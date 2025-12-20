package kafka_provider

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

type IConsumerService interface {
	HandleMessage(ctx context.Context, message *sarama.ConsumerMessage) error
}
type ConsumeAssignor string

const (
	RoundRobinAssignor ConsumeAssignor = "roundrobin"
	StickyAssignor     ConsumeAssignor = "sticky"
	RangeAssignor      ConsumeAssignor = "range"
)

type ConsumerConfig struct {
	Brokers       []string
	Topics        []string
	GroupID       string
	Assignor      ConsumeAssignor
	OffsetInitial int64
}

func (c *ConsumerConfig) ToKafkaConsumerConfig() *sarama.Config {
	kafkaConsumerConfig := sarama.NewConfig()
	switch c.Assignor {
	case StickyAssignor:
		kafkaConsumerConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case RoundRobinAssignor:
		kafkaConsumerConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case RangeAssignor:
		kafkaConsumerConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", c.Assignor)
	}
	kafkaConsumerConfig.Consumer.Offsets.Initial = c.OffsetInitial
	return kafkaConsumerConfig
}

type ConsumerMessageHandle struct {
	fHandlerError func(error)
	fReceive      func(message *sarama.ConsumerMessage)
}
type ConsumerApp struct {
	Config             ConsumerConfig
	ready              chan bool
	consumerGroup      sarama.ConsumerGroup
	handler            map[string]IConsumerService
	consumerMsgHandler *ConsumerMessageHandle
}

func NewConsumerApp(config ConsumerConfig) *ConsumerApp {
	kafkaConsumerConfig := config.ToKafkaConsumerConfig()
	client, err := sarama.NewConsumerGroup(config.Brokers, config.GroupID, kafkaConsumerConfig)
	if err != nil {
		panic(err)
	}
	return &ConsumerApp{
		Config:        config,
		ready:         make(chan bool),
		consumerGroup: client,
		handler:       make(map[string]IConsumerService),
		consumerMsgHandler: &ConsumerMessageHandle{
			fHandlerError: func(error) {},
			fReceive: func(message *sarama.ConsumerMessage) {
				log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, offset = %d", string(message.Value),
					message.Timestamp, message.Topic, message.Offset)
			},
		},
	}
}
func (c *ConsumerApp) Ready() chan bool {
	return c.ready
}

func (c *ConsumerApp) Run() error {
	keepRunning := true
	config := sarama.NewConfig()

	switch c.Config.Assignor {
	case StickyAssignor:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case RoundRobinAssignor:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case RangeAssignor:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", c.Config.Assignor)
	}

	config.Consumer.Offsets.Initial = c.Config.OffsetInitial

	ctx, cancel := context.WithCancel(context.Background())
	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := c.consumerGroup.Consume(ctx, c.Config.Topics, c); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			c.ready = make(chan bool)
		}
	}()

	<-c.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(c.consumerGroup, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err := c.consumerGroup.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
	return nil
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}

func (c *ConsumerApp) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *ConsumerApp) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (c *ConsumerApp) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			c.consumerMsgHandler.fReceive(message)
			session.MarkMessage(message, "")
			if handlerService, valid := c.handler[message.Topic]; !valid || handlerService == nil {
				continue
			} else {
				err := handlerService.HandleMessage(context.Background(), message)
				if err != nil {
					c.consumerMsgHandler.fHandlerError(err)
				}
			}
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *ConsumerApp) RegisterHandler(topic string, handler IConsumerService) *ConsumerApp {
	c.handler[topic] = handler
	return c
}

func (c *ConsumerApp) SetErrorMessageHandler(handler func(error)) *ConsumerApp {
	c.consumerMsgHandler.fHandlerError = handler
	return c
}

func (c *ConsumerApp) SetReceiveMessageHandler(handler func(message *sarama.ConsumerMessage)) *ConsumerApp {
	c.consumerMsgHandler.fReceive = handler
	return c
}
