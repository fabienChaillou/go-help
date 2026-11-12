package saramatest

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"
)

// OrderProcessor est la logique métier, isolée de Kafka.
type OrderProcessor interface {
	Process(ctx context.Context, evt OrderEvent) error
}

// OrderHandler implémente sarama.ConsumerGroupHandler.
type OrderHandler struct {
	processor OrderProcessor
	logger    *slog.Logger
}

func NewOrderHandler(p OrderProcessor, l *slog.Logger) *OrderHandler {
	if l == nil {
		l = slog.Default()
	}
	return &OrderHandler{processor: p, logger: l}
}

func (h *OrderHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *OrderHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *OrderHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			var evt OrderEvent
			if err := json.Unmarshal(msg.Value, &evt); err != nil {
				// Message empoisonné : on le log et on l'acquitte pour ne pas bloquer la partition.
				h.logger.Error("invalid payload", "offset", msg.Offset, "err", err)
				sess.MarkMessage(msg, "")
				continue
			}
			if err := h.processor.Process(sess.Context(), evt); err != nil {
				// Erreur métier : on n'acquitte PAS, le message sera relu après rebalance/redémarrage.
				h.logger.Error("process failed", "id", evt.ID, "err", err)
				return err
			}
			sess.MarkMessage(msg, "")
		case <-sess.Context().Done():
			return nil
		}
	}
}

// RunConsumer boucle sur Consume tant que le contexte n'est pas annulé
// (Consume doit être rappelé après chaque rebalance).
func RunConsumer(ctx context.Context, group sarama.ConsumerGroup, topics []string, h sarama.ConsumerGroupHandler) error {
	for {
		if err := group.Consume(ctx, topics, h); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return nil
		}
	}
}
