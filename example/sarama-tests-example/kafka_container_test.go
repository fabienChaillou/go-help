//go:build integration

package saramatest

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tckafka "github.com/testcontainers/testcontainers-go/modules/kafka"
)

// Configuration du Kafka de test, pilotable par variables d'environnement :
//
//	KAFKA_BROKERS     liste "host:port,host:port" d'un Kafka existant (CI, docker-compose).
//	                  Si définie, aucun conteneur n'est démarré.
//	KAFKA_IMAGE       image à utiliser (défaut : confluentinc/confluent-local:7.5.0).
//	KAFKA_TEST_LOGS   "1" pour afficher les logs du conteneur Kafka dans la sortie des tests.
//
// Variables propres à testcontainers utiles à connaître :
//
//	DOCKER_HOST                        socket Podman/Colima/Docker distant.
//	TESTCONTAINERS_RYUK_DISABLED=true  désactive le conteneur de nettoyage Ryuk
//	                                   (certains runners CI ne l'autorisent pas).
//	TESTCONTAINERS_HUB_IMAGE_NAME_PREFIX  préfixe de registry privé/miroir
//	                                   (ex. "registry.local/dockerhub/").
type KafkaConfig struct {
	ExternalBrokers []string
	Image           string
	ClusterID       string
	StreamLogs      bool
	StartupTimeout  time.Duration
}

const defaultKafkaImage = "confluentinc/confluent-local:7.5.0"

func KafkaConfigFromEnv() KafkaConfig {
	cfg := KafkaConfig{
		Image:          defaultKafkaImage,
		ClusterID:      "test-cluster",
		StreamLogs:     os.Getenv("KAFKA_TEST_LOGS") == "1",
		StartupTimeout: 2 * time.Minute,
	}
	if v := os.Getenv("KAFKA_BROKERS"); v != "" {
		cfg.ExternalBrokers = strings.Split(v, ",")
	}
	if v := os.Getenv("KAFKA_IMAGE"); v != "" {
		cfg.Image = v
	}
	return cfg
}

// KafkaEnv est l'environnement Kafka disponible pour les tests.
type KafkaEnv struct {
	Brokers   []string
	container *tckafka.KafkaContainer
}

// Terminate arrête le conteneur (no-op si on utilise un Kafka externe).
func (e *KafkaEnv) Terminate() error {
	if e.container == nil {
		return nil
	}
	return testcontainers.TerminateContainer(e.container)
}

// stdLogConsumer relaie les logs du conteneur vers la sortie standard.
type stdLogConsumer struct{}

func (stdLogConsumer) Accept(l testcontainers.Log) {
	log.Printf("[kafka] %s", strings.TrimRight(string(l.Content), "\n"))
}

func StartKafka(ctx context.Context, cfg KafkaConfig) (*KafkaEnv, error) {
	if len(cfg.ExternalBrokers) > 0 {
		return &KafkaEnv{Brokers: cfg.ExternalBrokers}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.StartupTimeout)
	defer cancel()

	opts := []testcontainers.ContainerCustomizer{
		tckafka.WithClusterID(cfg.ClusterID),
		testcontainers.WithEnv(map[string]string{
			// Les tests créent leurs topics explicitement : on désactive l'auto-création
			// pour détecter les fautes de frappe sur les noms de topic.
			"KAFKA_AUTO_CREATE_TOPICS_ENABLE": "false",
			// Rebalance immédiat au premier join : accélère nettement les tests de consumer group.
			"KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS": "0",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR": "1",
			"KAFKA_TRANSACTION_STATE_LOG_MIN_ISR":     "1",
		}),
		testcontainers.WithLabels(map[string]string{"app": "saramatest"}),
	}
	if cfg.StreamLogs {
		opts = append(opts, testcontainers.WithLogConsumers(stdLogConsumer{}))
	}

	kc, err := tckafka.Run(ctx, cfg.Image, opts...)
	if err != nil {
		// Run peut renvoyer un conteneur partiellement démarré : on le nettoie.
		if kc != nil {
			_ = testcontainers.TerminateContainer(kc)
		}
		return nil, fmt.Errorf("start kafka container (%s): %w", cfg.Image, err)
	}

	brokers, err := kc.Brokers(ctx)
	if err != nil {
		_ = testcontainers.TerminateContainer(kc)
		return nil, fmt.Errorf("get brokers: %w", err)
	}
	return &KafkaEnv{Brokers: brokers, container: kc}, nil
}
