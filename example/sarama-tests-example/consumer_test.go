package saramatest

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Fakes : sarama/mocks ne fournit pas de mock pour ConsumerGroup,
// mais ConsumerGroupSession et ConsumerGroupClaim sont des interfaces.

type fakeSession struct {
	sarama.ConsumerGroupSession // embarqué : les méthodes non utilisées paniquent si appelées
	ctx                         context.Context
	mu                          sync.Mutex
	marked                      []int64
}

func (s *fakeSession) Context() context.Context { return s.ctx }
func (s *fakeSession) MarkMessage(msg *sarama.ConsumerMessage, _ string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.marked = append(s.marked, msg.Offset)
}

type fakeClaim struct {
	sarama.ConsumerGroupClaim
	msgs chan *sarama.ConsumerMessage
}

func (c *fakeClaim) Messages() <-chan *sarama.ConsumerMessage { return c.msgs }

func newClaim(msgs ...*sarama.ConsumerMessage) *fakeClaim {
	ch := make(chan *sarama.ConsumerMessage, len(msgs))
	for _, m := range msgs {
		ch <- m
	}
	close(ch) // la boucle se termine quand le canal est vidé
	return &fakeClaim{msgs: ch}
}

type fakeProcessor struct {
	got []OrderEvent
	err error
}

func (p *fakeProcessor) Process(_ context.Context, e OrderEvent) error {
	p.got = append(p.got, e)
	return p.err
}

func msg(offset int64, value string) *sarama.ConsumerMessage {
	return &sarama.ConsumerMessage{Topic: "orders", Offset: offset, Value: []byte(value)}
}

// --- Tests (table-driven) ---

func TestOrderHandler_ConsumeClaim(t *testing.T) {
	tests := []struct {
		name       string
		msgs       []*sarama.ConsumerMessage
		procErr    error
		wantErr    bool
		wantMarked []int64
		wantCalls  int
	}{
		{
			name:       "messages valides acquittés",
			msgs:       []*sarama.ConsumerMessage{msg(0, `{"id":"a","amount":1}`), msg(1, `{"id":"b","amount":2}`)},
			wantMarked: []int64{0, 1},
			wantCalls:  2,
		},
		{
			name:       "payload invalide ignoré mais acquitté",
			msgs:       []*sarama.ConsumerMessage{msg(0, `not-json`), msg(1, `{"id":"b"}`)},
			wantMarked: []int64{0, 1},
			wantCalls:  1,
		},
		{
			name:       "erreur métier : pas d'acquittement, on remonte l'erreur",
			msgs:       []*sarama.ConsumerMessage{msg(7, `{"id":"a"}`)},
			procErr:    errors.New("db down"),
			wantErr:    true,
			wantMarked: nil,
			wantCalls:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proc := &fakeProcessor{err: tt.procErr}
			sess := &fakeSession{ctx: context.Background()}
			h := NewOrderHandler(proc, nil)

			err := h.ConsumeClaim(sess, newClaim(tt.msgs...))

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.wantMarked, sess.marked)
			assert.Len(t, proc.got, tt.wantCalls)
		})
	}
}

func TestOrderHandler_StopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sess := &fakeSession{ctx: ctx}
	claim := &fakeClaim{msgs: make(chan *sarama.ConsumerMessage)} // jamais fermé

	done := make(chan error, 1)
	go func() { done <- NewOrderHandler(&fakeProcessor{}, nil).ConsumeClaim(sess, claim) }()

	cancel()
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("ConsumeClaim ne s'est pas arrêté après l'annulation du contexte")
	}
}
