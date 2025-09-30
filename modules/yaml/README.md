En Go, tu peux tester une structure YAML en utilisant le package standard [`gopkg.in/yaml.v3`](https://pkg.go.dev/gopkg.in/yaml.v3) (ou `yaml.v2`).
Le principe est :

1. Définir une **struct** qui reflète la structure YAML attendue.
2. Charger le YAML (string ou fichier).
3. Le parser avec `yaml.Unmarshal`.
4. Écrire un test avec `testing` pour valider le résultat.

Voici un exemple minimal :

```go
package config

import (
	"testing"

	"gopkg.in/yaml.v3"
	"github.com/stretchr/testify/assert"
)

// Exemple de structure correspondant au YAML
type AppConfig struct {
	Name    string   `yaml:"name"`
	Port    int      `yaml:"port"`
	Enabled bool     `yaml:"enabled"`
	Tags    []string `yaml:"tags"`
}

func TestAppConfigYAML(t *testing.T) {
	yamlData := `
name: "myapp"
port: 8080
enabled: true
tags:
  - "go"
  - "yaml"
`

	var cfg AppConfig
	err := yaml.Unmarshal([]byte(yamlData), &cfg)
	assert.NoError(t, err)

	assert.Equal(t, "myapp", cfg.Name)
	assert.Equal(t, 8080, cfg.Port)
	assert.True(t, cfg.Enabled)
	assert.Equal(t, []string{"go", "yaml"}, cfg.Tags)
}
```

👉 Ce test charge un YAML en mémoire et vérifie que le mapping dans ta structure Go est correct.

---

Si tu veux tester avec un fichier YAML, tu peux remplacer :

```go
data, err := os.ReadFile("config.yaml")
err = yaml.Unmarshal(data, &cfg)
```

---

Veux-tu que je te montre aussi **comment tester une structure YAML invalide** (par exemple, un champ manquant ou de mauvais type) pour couvrir les cas d’erreurs ?

Parfait 👍
En Go, tu peux tester les cas **invalide** ou **incomplet** en fournissant un YAML volontairement incorrect et en vérifiant que :

* soit le parsing (`yaml.Unmarshal`) renvoie une erreur,
* soit les valeurs dans la struct ne correspondent pas à ce que tu attends.

Voici un exemple complet avec **plusieurs cas de test** :

```go
package config

import (
	"testing"

	"gopkg.in/yaml.v3"
	"github.com/stretchr/testify/assert"
)

// Notre structure cible
type AppConfig struct {
	Name    string   `yaml:"name"`
	Port    int      `yaml:"port"`
	Enabled bool     `yaml:"enabled"`
	Tags    []string `yaml:"tags"`
}

func TestAppConfigYAML_Invalid(t *testing.T) {
	tests := []struct {
		name      string
		yamlInput string
		wantError bool
	}{
		{
			name: "Champ manquant (pas de port)",
			yamlInput: `
name: "myapp"
enabled: true
tags: ["go"]
`,
			wantError: false, // pas d’erreur, mais Port restera 0
		},
		{
			name: "Type invalide (port string au lieu d’un entier)",
			yamlInput: `
name: "badapp"
port: "huitmille"
enabled: true
`,
			wantError: true, // ici yaml.Unmarshal échouera
		},
		{
			name: "YAML mal formé",
			yamlInput: `
name: "oops
port: 8080
`,
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var cfg AppConfig
			err := yaml.Unmarshal([]byte(tc.yamlInput), &cfg)

			if tc.wantError {
				assert.Error(t, err, "on attendait une erreur")
			} else {
				assert.NoError(t, err, "on ne devait pas avoir d'erreur")
			}
		})
	}
}
```

🔎 Ce que fait ce test :

* **Cas 1** → champ `port` absent → pas d’erreur YAML, mais `cfg.Port` vaudra `0` (valeur par défaut d’un int).
* **Cas 2** → `port` est une string au lieu d’un int → `yaml.Unmarshal` renvoie une erreur.
* **Cas 3** → YAML mal indenté/fermé → `yaml.Unmarshal` renvoie une erreur de parsing.

---

Veux-tu que je t’ajoute aussi un exemple de **validation métier** (ex: "le port doit être > 0") avec un package comme `go-playground/validator`, en plus du test de parsing pur YAML ?

Super 👍
Tu peux en effet séparer deux étapes :

1. **Parsing YAML → struct Go** avec `yaml.Unmarshal`.
2. **Validation métier** sur la struct avec [`go-playground/validator`](https://pkg.go.dev/github.com/go-playground/validator/v10).

Voici un exemple complet qui combine les deux :

```go
package config

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Name    string   `yaml:"name" validate:"required"`
	Port    int      `yaml:"port" validate:"required,min=1,max=65535"`
	Enabled bool     `yaml:"enabled"`
	Tags    []string `yaml:"tags" validate:"dive,required"`
}

var validate = validator.New()

func TestAppConfigYAML_WithValidation(t *testing.T) {
	tests := []struct {
		name      string
		yamlInput string
		wantError bool
	}{
		{
			name: "Config valide",
			yamlInput: `
name: "myapp"
port: 8080
enabled: true
tags: ["go","yaml"]
`,
			wantError: false,
		},
		{
			name: "Port manquant",
			yamlInput: `
name: "app-without-port"
enabled: true
tags: ["test"]
`,
			wantError: true, // fail car "port" est required
		},
		{
			name: "Port invalide (0)",
			yamlInput: `
name: "bad-port"
port: 0
enabled: true
tags: ["test"]
`,
			wantError: true, // fail car min=1
		},
		{
			name: "Tag vide",
			yamlInput: `
name: "bad-tags"
port: 9000
tags: [""]
`,
			wantError: true, // fail car un tag est vide
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var cfg AppConfig
			err := yaml.Unmarshal([]byte(tc.yamlInput), &cfg)
			assert.NoError(t, err, "le YAML doit être valide syntaxiquement")

			// Validation métier
			err = validate.Struct(cfg)
			if tc.wantError {
				assert.Error(t, err, "on attendait une erreur de validation")
			} else {
				assert.NoError(t, err, "on ne devait pas avoir d'erreur de validation")
			}
		})
	}
}
```

---

✅ Ce que tu gagnes avec `validator` :

* Vérifier automatiquement si un champ est **obligatoire** (`required`),
* Vérifier des **contraintes numériques** (`min`, `max`),
* Vérifier des **tableaux** avec `dive` (ici : chaque `tag` doit être non vide).

---

Veux-tu que je te montre aussi comment brancher ça avec la lecture d’un **fichier `config.yaml` réel** (et non une string dans le test) pour reproduire des conditions "production" ?

Parfait 🚀
Je vais te montrer comment :

1. Lire une configuration YAML depuis un **fichier réel** (`config.yaml`).
2. L’unmarshal dans ta struct Go.
3. La valider avec `go-playground/validator`.
4. Ajouter des **tests unitaires** qui utilisent des fichiers YAML valides et invalides.

---

### Exemple d’implémentation

**📂 arborescence :**

```
project/
 ├── config/
 │    ├── config.go
 │    ├── config.yaml
 │    ├── config_invalid.yaml
 │    └── config_test.go
```

---

**config/config.go** – code de chargement et validation :

```go
package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Name    string   `yaml:"name" validate:"required"`
	Port    int      `yaml:"port" validate:"required,min=1,max=65535"`
	Enabled bool     `yaml:"enabled"`
	Tags    []string `yaml:"tags" validate:"dive,required"`
}

var validate = validator.New()

// LoadConfig lit et valide un fichier YAML
func LoadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
```

---

**config/config.yaml** – version valide :

```yaml
name: "myapp"
port: 8080
enabled: true
tags:
  - "go"
  - "yaml"
```

---

**config/config\_invalid.yaml** – version invalide (port absent) :

```yaml
name: "badapp"
enabled: false
tags:
  - ""
```

---

**config/config\_test.go** – tests avec fichiers :

```go
package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Valid(t *testing.T) {
	cfg, err := LoadConfig("config.yaml")
	assert.NoError(t, err, "le fichier config.yaml devrait être valide")
	assert.Equal(t, "myapp", cfg.Name)
	assert.Equal(t, 8080, cfg.Port)
	assert.True(t, cfg.Enabled)
	assert.Equal(t, []string{"go", "yaml"}, cfg.Tags)
}

func TestLoadConfig_Invalid(t *testing.T) {
	cfg, err := LoadConfig("config_invalid.yaml")
	assert.Error(t, err, "le fichier config_invalid.yaml devrait être invalide")
	assert.Nil(t, cfg)
}
```

---

✅ Résultat :

* `go test ./config -v` passera sur le fichier **valide**,
* échouera volontairement sur le fichier **invalide** (grâce aux règles `validator`).

---

Veux-tu que je t’ajoute aussi un **test avec un fichier YAML inexistant** (pour couvrir le cas d’erreur `os.ReadFile`) ?
