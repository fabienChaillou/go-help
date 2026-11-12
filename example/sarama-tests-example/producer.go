package saramatest

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

// OrderEvent est un exemple de message métier publié sur Kafka.
type OrderEvent struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

// Publisher dépend de l'interface sarama.SyncProducer (et non d'une
// implémentation concrète) : c'est ce qui permet d'injecter
// mocks.SyncProducer dans les tests unitaires.
type Publisher struct {
	producer sarama.SyncProducer
	topic    string
}

func NewPublisher(p sarama.SyncProducer, topic string) *Publisher {
	return &Publisher{producer: p, topic: topic}
}

func (p *Publisher) PublishOrder(ctx context.Context, evt OrderEvent) error {
	if evt.ID == "" {
		return fmt.Errorf("order id is required")
	}
	payload, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshal order: %w", err)
	}
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(evt.ID), // même clé => même partition => ordre garanti
		Value: sarama.ByteEncoder(payload),
		Headers: []sarama.RecordHeader{
			{Key: []byte("event-type"), Value: []byte("order.created")},
		},
	}
	if _, _, err := p.producer.SendMessage(msg); err != nil {
		return fmt.Errorf("send order %s: %w", evt.ID, err)
	}
	return nil
}
