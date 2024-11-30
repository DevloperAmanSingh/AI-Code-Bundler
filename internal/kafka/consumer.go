package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topic    string
}

func NewKafkaConsumer(brokers, groupID, topic string) (*KafkaConsumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	return &KafkaConsumer{
		consumer: consumer,
		topic:    topic,
	}, nil
}

func (kc *KafkaConsumer) Receive() (string, error) {
	msg, err := kc.consumer.ReadMessage(-1)
	if err != nil {
		return "", err
	}
	return string(msg.Value), nil
}

func (kc *KafkaConsumer) Close() {
	if kc.consumer != nil {
		kc.consumer.Close()
	}
}
