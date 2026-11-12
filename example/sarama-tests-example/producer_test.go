package saramatest

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublishOrder_Success(t *testing.T) {
	// mocks.SyncProducer échoue le test (via t) si une attente n'est pas consommée
	// ou si SendMessage est appelé sans attente déclarée.
	prod := mocks.NewSyncProducer(t, mocks.NewTestConfig())
	defer func() { require.NoError(t, prod.Close()) }()

	// Vérifie le contenu complet du message (clé, headers, topic, payload).
	prod.ExpectSendMessageWithMessageCheckerFunctionAndSucceed(func(msg *sarama.ProducerMessage) error {
		assert.Equal(t, "orders", msg.Topic)

		key, _ := msg.Key.Encode()
		assert.Equal(t, "order-42", string(key))

		require.Len(t, msg.Headers, 1)
		assert.Equal(t, "order.created", string(msg.Headers[0].Value))

		raw, _ := msg.Value.Encode()
		var got OrderEvent
		require.NoError(t, json.Unmarshal(raw, &got))
		assert.Equal(t, OrderEvent{ID: "order-42", Amount: 99.5}, got)
		return nil
	})

	pub := NewPublisher(prod, "orders")
	require.NoError(t, pub.PublishOrder(context.Background(), OrderEvent{ID: "order-42", Amount: 99.5}))
}

func TestPublishOrder_BrokerError(t *testing.T) {
	prod := mocks.NewSyncProducer(t, mocks.NewTestConfig())
	defer func() { _ = prod.Close() }()

	prod.ExpectSendMessageAndFail(sarama.ErrNotLeaderForPartition)

	err := NewPublisher(prod, "orders").PublishOrder(context.Background(), OrderEvent{ID: "x"})
	require.Error(t, err)
	assert.True(t, errors.Is(err, sarama.ErrNotLeaderForPartition), "l'erreur sarama doit être wrappée")
}

func TestPublishOrder_Validation(t *testing.T) {
	// Aucune attente : si le code appelait SendMessage, le mock ferait échouer le test.
	prod := mocks.NewSyncProducer(t, mocks.NewTestConfig())
	defer func() { _ = prod.Close() }()

	err := NewPublisher(prod, "orders").PublishOrder(context.Background(), OrderEvent{})
	assert.EqualError(t, err, "order id is required")
}
