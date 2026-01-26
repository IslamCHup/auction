package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
)

type BidPlacedEvent struct {
	LotID            uint64 `json:"lot_id"`
	PreviousLeaderID uint64 `json:"previous_leader_id"`
	NewBidAmount     int64  `json:"new_bid_amount"`
}

type LotCompletedEvent struct {
	LotID      uint64   `json:"lot_id"`
	Winner     uint64   `json:"winner"`
	FinalPrice int64    `json:"final_price"`
	LoserIDs   []uint64 `json:"loser_ids"`
}

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer() (*Producer, error) {
	brokerStr := os.Getenv("KAFKA_BROKERS")
	var brokers []string
	if brokerStr == "" {
		brokers = []string{"kafka:9092"}
	} else {
		brokers = []string{brokerStr}
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	log.Printf("Kafka producer connected to brokers: %v", brokers)
	return &Producer{producer: producer}, nil
}

func (p *Producer) SendMessage(topic string, key string, value interface{}) error {
	if p == nil || p.producer == nil {
		return fmt.Errorf("producer is not initialized")
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(jsonValue),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("Message sent to topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}

func (p *Producer) Close() error {
	if p == nil || p.producer == nil {
		return nil
	}
	return p.producer.Close()
}
