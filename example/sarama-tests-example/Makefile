GO            ?= go
PKG           ?= ./...
KAFKA_IMAGE   ?= confluentinc/confluent-local:7.5.0
KAFKA_NAME    ?= saramatest-kafka
KAFKA_PORT    ?= 9092
COVER_FILE    ?= coverage.out
RUN           ?=

.DEFAULT_GOAL := help

.PHONY: help
help: ## Affiche cette aide
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2}'

## ---------- Qualité ----------

.PHONY: tidy
tidy: ## go mod tidy
	$(GO) mod tidy

.PHONY: vet
vet: ## go vet (y compris les fichiers taggés integration)
	$(GO) vet $(PKG)
	$(GO) vet -tags=integration $(PKG)

## ---------- Tests ----------

.PHONY: test
test: test-unit ## Alias de test-unit

.PHONY: test-unit
test-unit: ## Tests unitaires (sans Docker)
	$(GO) test -race -count=1 $(if $(RUN),-run '$(RUN)') $(PKG)

.PHONY: test-integration
test-integration: ## Tests d'intégration : Kafka démarré par testcontainers (Docker requis)
	$(GO) test -tags=integration -race -count=1 -timeout 10m -v \
		-run '$(if $(RUN),$(RUN),TestKafkaIntegration)' $(PKG)

.PHONY: test-integration-logs
test-integration-logs: ## Idem, avec les logs du conteneur Kafka
	KAFKA_TEST_LOGS=1 $(MAKE) test-integration

.PHONY: test-integration-local
test-integration-local: ## Tests d'intégration contre le Kafka lancé par `make kafka-up` (plus rapide en boucle de dev)
	KAFKA_BROKERS=localhost:$(KAFKA_PORT) $(MAKE) test-integration

.PHONY: test-all
test-all: vet test-unit test-integration ## vet + unitaires + intégration

.PHONY: cover
cover: ## Couverture unitaires + intégration, rapport HTML dans coverage.html
	$(GO) test -tags=integration -race -count=1 -timeout 10m \
		-covermode=atomic -coverprofile=$(COVER_FILE) $(PKG)
	$(GO) tool cover -func=$(COVER_FILE) | tail -1
	$(GO) tool cover -html=$(COVER_FILE) -o coverage.html

## ---------- Kafka local (optionnel) ----------

.PHONY: kafka-up
kafka-up: ## Démarre un Kafka local persistant sur localhost:$(KAFKA_PORT)
	@docker inspect $(KAFKA_NAME) >/dev/null 2>&1 || \
		docker run -d --name $(KAFKA_NAME) -p $(KAFKA_PORT):9092 \
			-e KAFKA_AUTO_CREATE_TOPICS_ENABLE=false \
			-e KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS=0 \
			$(KAFKA_IMAGE)
	@echo "Attente de Kafka..."
	@until docker exec $(KAFKA_NAME) kafka-topics --bootstrap-server localhost:9092 --list >/dev/null 2>&1; do sleep 1; done
	@echo "Kafka prêt sur localhost:$(KAFKA_PORT)"

.PHONY: kafka-down
kafka-down: ## Arrête et supprime le Kafka local
	-docker rm -f $(KAFKA_NAME)

.PHONY: kafka-topics
kafka-topics: ## Liste les topics du Kafka local
	docker exec $(KAFKA_NAME) kafka-topics --bootstrap-server localhost:9092 --list

.PHONY: clean
clean: ## Supprime les artefacts de couverture
	rm -f $(COVER_FILE) coverage.html
