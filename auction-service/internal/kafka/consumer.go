package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/IBM/sarama"
)

type MessageHandler func(topic string, key string, value map[string]interface{}) error

type Consumer struct {
	consumer sarama.ConsumerGroup
	handler  MessageHandler
	topics   []string
	wg       sync.WaitGroup
}

func NewConsumer(groupID string, topics []string, handler MessageHandler) (*Consumer, error) {
	brokerStr := os.Getenv("KAFKA_BROKERS")
	var brokers []string
	if brokerStr == "" {
		brokers = []string{"localhost:9092"}
	} else {
		brokers = []string{brokerStr}
	}

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.V2_6_0_0

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer group: %w", err)
	}

	log.Printf("Kafka consumer group '%s' connected to brokers: %v", groupID, brokers)
	return &Consumer{
		consumer: consumerGroup,
		handler:  handler,
		topics:   topics,
	}, nil
}

func (c *Consumer) Start() error {
	if c == nil || c.consumer == nil {
		return fmt.Errorf("consumer is not initialized")
	}

	ctx := context.Background()
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			err := c.consumer.Consume(ctx, c.topics, c)
			if err != nil {
				log.Printf("Error from consumer: %v", err)
				return
			}
		}
	}()

	log.Printf("Kafka consumer started, listening to topics: %v", c.topics)
	return nil
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			var value map[string]interface{}
			if err := json.Unmarshal(message.Value, &value); err != nil {
				log.Printf("Failed to unmarshal message from topic %s: %v", message.Topic, err)
				session.MarkMessage(message, "")
				continue
			}

			if c.handler != nil {
				if err := c.handler(message.Topic, string(message.Key), value); err != nil {
					log.Printf("Error handling message from topic %s: %v", message.Topic, err)
				}
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *Consumer) Close() error {
	if c == nil || c.consumer == nil {
		return nil
	}
	c.wg.Wait()
	return c.consumer.Close()
}
