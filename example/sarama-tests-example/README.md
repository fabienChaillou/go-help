Je vais te préparer un petit projet d'exemple compilable qui montre les deux niveaux : des tests unitaires avec le package `sarama/mocks`, et des tests d'intégration contre un vrai Kafka via testcontainers.

No mock for `ConsumerGroup` in sarama — that's an important point; I'll write lightweight fakes. Now the code.

Unit tests pass. Now the integration test with testcontainers.

Everything compiles and `go vet` is clean. I'll try starting a Docker daemon here to actually run the integration tests.

The Kafka image can't be pulled here (Docker Hub is blocked from my environment), so the integration tests compile but couldn't run. I'll package the files.

Je t'ai préparé un petit projet d'exemple dans `sarama-tests-example.zip` : un producer, un consumer group, et les deux niveaux de tests. Les tests unitaires passent (y compris avec `-race`). Les tests d'intégration compilent et `go vet` ne signale rien, mais je n'ai pas pu les exécuter : mon environnement ne peut pas télécharger l'image Kafka depuis Docker Hub. Lance-les une fois chez toi avant de t'appuyer dessus.

**L'idée principale :** ton code doit dépendre des **interfaces** de Sarama (`sarama.SyncProducer`, `sarama.ConsumerGroupHandler`, `sarama.ConsumerGroupSession`…) et pas de structures concrètes. C'est ce qui permet de remplacer Kafka dans les tests. Garde aussi la logique métier derrière ta propre interface (`OrderProcessor` dans l'exemple), séparée du code qui lit Kafka.

**Tests unitaires (sans Kafka)**
- **Producer :** Sarama fournit `github.com/IBM/sarama/mocks`. On déclare à l'avance ce qu'on attend, par exemple avec `ExpectSendMessageWithMessageCheckerFunctionAndSucceed` (qui permet de vérifier le topic, la clé, les headers et le contenu) ou `ExpectSendMessageAndFail(err)`. Le mock fait échouer le test si une attente n'est pas utilisée, ou si le code envoie un message alors qu'on n'en attendait aucun.
- **Consumer group :** le package `mocks` ne fournit rien pour `ConsumerGroup`. J'ai donc écrit de petits faux objets pour `ConsumerGroupSession` et `ConsumerGroupClaim` : un canal de messages qu'on remplit à l'avance, et la liste des offsets acquittés. Ça permet de tester `ConsumeClaim` directement : message valide acquitté, payload invalide ignoré mais acquitté, erreur métier non acquittée, arrêt quand le contexte est annulé.

**Tests d'intégration (vrai Kafka)**
- Ils sont derrière le build tag `integration` : un simple `go test` ne les lance pas.
- `TestMain` démarre **un seul** Kafka sans ZooKeeper pour tout le package, avec le module Kafka de testcontainers (image `confluentinc/confluent-local:7.5.0`). Si la variable `KAFKA_BROKERS` est définie, il utilise ce Kafka à la place, ce qui est pratique en CI.
- Chaque test crée son propre topic et son propre groupe, pour que les tests ne se gênent pas entre eux.
- Deux scénarios : un message publié puis bien reçu par le consumer, et un message qui échoue au traitement puis est relu par le consommateur suivant du même groupe (garantie « au moins une fois »).

```bash
go mod tidy
go test -race ./...                       # unitaires
go test -tags=integration -race -v ./...  # intégration (Docker requis)
```

Pour les versions : la dernière version de Sarama (v1.61) demande Go 1.26. J'ai testé avec Sarama v1.46.3 et testcontainers v0.38.0 sous Go 1.24. Garde la version de Sarama que ton appli utilise déjà.

---

Je l'ajoute aussi.

J'ai ajouté la config testcontainers, la suite testify et le Makefile ; tout est dans le nouveau `sarama-tests-example.zip`. Les tests d'intégration compilent, `go vet` ne signale rien et les tests unitaires passent. En revanche, je n'ai toujours pas pu lancer les tests d'intégration contre un vrai Kafka : mon environnement ne peut pas télécharger l'image depuis Docker Hub. Lance-les une fois chez toi avec `make test-integration` pour confirmer.

**Config testcontainers** (`kafka_container_test.go`)
- `KafkaConfigFromEnv()` lit trois variables :
  - `KAFKA_BROKERS` : utilise un Kafka existant au lieu de démarrer un conteneur (utile en CI).
  - `KAFKA_IMAGE` : choisit une autre image.
  - `KAFKA_TEST_LOGS=1` : affiche les logs du conteneur.
- `StartKafka()` démarre Kafka avec quelques réglages pour les tests :
  - l'auto-création des topics est désactivée, ce qui fait échouer un test si un nom de topic est mal tapé ;
  - le délai avant le premier rebalance est à 0, ce qui fait gagner plusieurs secondes par test de consumer ;
  - le conteneur est supprimé si le démarrage échoue en cours de route.
- J'ai aussi mis en commentaire les variables propres à testcontainers : `DOCKER_HOST` pour Podman ou Colima, `TESTCONTAINERS_RYUK_DISABLED` et le préfixe pour un miroir de registry.

**Suite testify** (`integration_test.go`)
- `SetupSuite` / `TearDownSuite` démarrent et arrêtent un seul Kafka pour toute la suite.
- `SetupTest` / `TearDownTest` donnent à chaque test son propre topic (3 partitions), son groupe et son producer, puis suppriment le topic à la fin.
- La suite contient quatre tests :
  - un message publié est bien consommé ;
  - les messages de même clé arrivent dans l'ordre ;
  - un message en échec est relu ;
  - un message commité n'est pas relu.
- Avec `-short`, la suite est ignorée.
- Correction d'un bug de ma version précédente : le nom du test contient des `/`, qui sont interdits dans un nom de topic. Les noms de topic sont maintenant générés autrement.

**Makefile** (`make help` liste toutes les commandes)
- `test-unit` : tests unitaires, sans Docker.
- `test-integration` : Kafka démarré par testcontainers ; `test-integration-logs` fait pareil avec les logs du conteneur.
- `kafka-up`, puis `test-integration-local`, puis `kafka-down` : garde un Kafka lancé pendant que tu développes, ce qui évite de redémarrer un conteneur à chaque lancement.
- `RUN=...` lance un seul test, par exemple `make test-integration RUN='TestKafkaIntegration/TestOrderingPerKey'`.
- `cover` produit le rapport de couverture HTML, et `test-all` enchaîne vet, unitaires et intégration.
