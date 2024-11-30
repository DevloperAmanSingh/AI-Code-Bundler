package kafka

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(brokers, clientID, topic string) (*KafkaProducer, error) {
	log.Println("Initializing Kafka producer...")

	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"client.id":         clientID,
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (kp *KafkaProducer) Send(message string) error {
	if kp.producer == nil {
		return fmt.Errorf("producer not initialized")
	}

	err := kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kp.topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (kp *KafkaProducer) Close() {
	if kp.producer != nil {
		kp.producer.Close()
	}
}
