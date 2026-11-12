//go:build integration

// Lancer avec : go test -tags=integration -race -v ./...
// Un seul test :  go test -tags=integration -v -run 'TestKafkaIntegration/TestPublishThenConsume' ./...
package saramatest

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/suite"
)

// KafkaIntegrationSuite partage un seul Kafka entre tous les tests
// (SetupSuite / TearDownSuite) et donne à chaque test un topic,
// un groupe et un producer neufs (SetupTest / TearDownTest).
type KafkaIntegrationSuite struct {
	suite.Suite

	env   *KafkaEnv
	admin sarama.ClusterAdmin

	// Ressources propres au test courant.
	topic    string
	groupID  string
	producer sarama.SyncProducer
	pub      *Publisher
}

// Point d'entrée `go test`.
func TestKafkaIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("tests d'intégration ignorés avec -short")
	}
	suite.Run(t, new(KafkaIntegrationSuite))
}

// ---------- Cycle de vie ----------

func (s *KafkaIntegrationSuite) SetupSuite() {
	env, err := StartKafka(context.Background(), KafkaConfigFromEnv())
	s.Require().NoError(err)
	s.env = env

	admin, err := sarama.NewClusterAdmin(env.Brokers, s.saramaConfig())
	s.Require().NoError(err)
	s.admin = admin
}

func (s *KafkaIntegrationSuite) TearDownSuite() {
	if s.admin != nil {
		s.NoError(s.admin.Close())
	}
	if s.env != nil {
		s.NoError(s.env.Terminate())
	}
}

var topicSeq atomic.Int64

func (s *KafkaIntegrationSuite) SetupTest() {
	// Nom de topic unique (s.T().Name() contient des "/", interdits dans un nom de topic).
	s.topic = fmt.Sprintf("orders-it-%d-%d", time.Now().UnixNano(), topicSeq.Add(1))
	s.groupID = "group-" + s.topic

	s.Require().NoError(s.admin.CreateTopic(s.topic, &sarama.TopicDetail{
		NumPartitions:     3,
		ReplicationFactor: 1,
	}, false))
	s.waitForTopic(s.topic)

	producer, err := sarama.NewSyncProducer(s.env.Brokers, s.saramaConfig())
	s.Require().NoError(err)
	s.producer = producer
	s.pub = NewPublisher(producer, s.topic)
}

func (s *KafkaIntegrationSuite) TearDownTest() {
	if s.producer != nil {
		s.NoError(s.producer.Close())
	}
	// Nettoyage best effort : n'échoue pas le test si la suppression échoue.
	_ = s.admin.DeleteTopic(s.topic)
}

// ---------- Helpers ----------

func (s *KafkaIntegrationSuite) saramaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.ClientID = "saramatest-it"
	cfg.Version = sarama.V3_5_0_0 // confluent-local 7.5 = Kafka 3.5
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewHashPartitioner // même clé => même partition
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Offsets.AutoCommit.Interval = 100 * time.Millisecond
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	return cfg
}

// waitForTopic attend que les métadonnées du topic soient propagées,
// sinon le premier SendMessage peut échouer avec UnknownTopicOrPartition.
func (s *KafkaIntegrationSuite) waitForTopic(topic string) {
	s.Require().Eventually(func() bool {
		md, err := s.admin.DescribeTopics([]string{topic})
		return err == nil && len(md) == 1 && md[0].Err == sarama.ErrNoError
	}, 10*time.Second, 100*time.Millisecond, "topic %s non disponible", topic)
}

// chanProcessor transmet chaque événement reçu sur un canal ; fail simule une panne.
type chanProcessor struct {
	events chan OrderEvent
	fail   error
}

func newChanProcessor(fail error) *chanProcessor {
	return &chanProcessor{events: make(chan OrderEvent, 100), fail: fail}
}

func (p *chanProcessor) Process(_ context.Context, e OrderEvent) error {
	p.events <- e
	return p.fail
}

// startConsumer lance un consumer group en arrière-plan et renvoie une fonction
// d'arrêt qui attend la fin propre (commit des offsets marqués compris).
func (s *KafkaIntegrationSuite) startConsumer(proc OrderProcessor) (stop func()) {
	group, err := sarama.NewConsumerGroup(s.env.Brokers, s.groupID, s.saramaConfig())
	s.Require().NoError(err)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	t := s.T()
	go func() {
		defer close(done)
		// Pas de s.Require() ici : FailNow ne doit pas être appelé hors de la goroutine du test.
		if err := RunConsumer(ctx, group, []string{s.topic}, NewOrderHandler(proc, nil)); err != nil &&
			!errors.Is(err, sarama.ErrClosedConsumerGroup) {
			t.Logf("consumer stopped: %v", err)
		}
	}()

	var stopped atomic.Bool
	stop = func() {
		if !stopped.CompareAndSwap(false, true) {
			return
		}
		cancel()
		s.NoError(group.Close())
		<-done
	}
	s.T().Cleanup(stop) // filet de sécurité si le test échoue avant stop()
	return stop
}

func (s *KafkaIntegrationSuite) receive(ch <-chan OrderEvent) OrderEvent {
	select {
	case e := <-ch:
		return e
	case <-time.After(30 * time.Second):
		s.FailNow("timeout : aucun message reçu")
		return OrderEvent{}
	}
}

// ---------- Tests ----------

func (s *KafkaIntegrationSuite) TestPublishThenConsume() {
	want := OrderEvent{ID: "order-1", Amount: 12.5}
	s.Require().NoError(s.pub.PublishOrder(context.Background(), want))

	proc := newChanProcessor(nil)
	defer s.startConsumer(proc)()

	s.Equal(want, s.receive(proc.events))
}

// Les messages de même clé arrivent dans l'ordre d'envoi (même partition).
func (s *KafkaIntegrationSuite) TestOrderingPerKey() {
	for i := 1; i <= 5; i++ {
		s.Require().NoError(s.pub.PublishOrder(context.Background(), OrderEvent{ID: "same-key", Amount: float64(i)}))
	}

	proc := newChanProcessor(nil)
	defer s.startConsumer(proc)()

	for i := 1; i <= 5; i++ {
		s.Equal(float64(i), s.receive(proc.events).Amount, "message %d hors ordre", i)
	}
}

// Comportement "at-least-once" : un message dont le traitement échoue n'est pas
// commité et doit être relu par le consommateur suivant du même groupe.
func (s *KafkaIntegrationSuite) TestFailedMessageIsRedelivered() {
	s.Require().NoError(s.pub.PublishOrder(context.Background(), OrderEvent{ID: "retry-me"}))

	failing := newChanProcessor(errors.New("db down"))
	stop := s.startConsumer(failing)
	s.Equal("retry-me", s.receive(failing.events).ID)
	stop()

	ok := newChanProcessor(nil)
	defer s.startConsumer(ok)()
	s.Equal("retry-me", s.receive(ok.events).ID)
}

// Un message traité avec succès est commité : un nouveau consommateur du groupe ne le relit pas.
func (s *KafkaIntegrationSuite) TestCommittedMessageIsNotRedelivered() {
	s.Require().NoError(s.pub.PublishOrder(context.Background(), OrderEvent{ID: "once"}))

	first := newChanProcessor(nil)
	stop := s.startConsumer(first)
	s.Equal("once", s.receive(first.events).ID)
	stop() // Close() commite les offsets marqués

	// On publie un marqueur : le 2e consommateur doit le recevoir EN PREMIER.
	s.Require().NoError(s.pub.PublishOrder(context.Background(), OrderEvent{ID: "marker"}))
	second := newChanProcessor(nil)
	defer s.startConsumer(second)()

	deadline := time.After(30 * time.Second)
	for {
		select {
		case e := <-second.events:
			s.Require().NotEqual("once", e.ID, "message déjà commité relu")
			if e.ID == "marker" {
				return
			}
		case <-deadline:
			s.FailNow("timeout : marqueur non reçu")
		}
	}
}
