Oui 👍 je comprends ton besoin : tu veux une fonction qui valide que le champ `DataSourceDefault` de ton `AppConfig` correspond bien à une des clés `Name` de la liste `Tags []Datasource`.

Je te propose une méthode de validation simple :

```go
package config

import (
	"fmt"
)

type AppConfig struct {
	Name              string       `yaml:"name" validate:"required"`
	Port              int          `yaml:"port" validate:"required,min=1,max=65535"`
	Enabled           bool         `yaml:"enabled"`
	DataSourceDefault string       `yaml:"datasource_default"`
	Tags              []Datasource `yaml:"data_sources"`
}

type Datasource struct {
	Name string `yaml:"name"`
	URI  string `yaml:"uri"`
}

// Validate vérifie que DataSourceDefault est bien présent dans Tags.Name
func (c *AppConfig) Validate() error {
	if c.DataSourceDefault == "" {
		return fmt.Errorf("datasource_default ne peut pas être vide")
	}

	for _, ds := range c.Tags {
		if ds.Name == c.DataSourceDefault {
			return nil // trouvé ✅
		}
	}

	return fmt.Errorf("datasource_default %q n'existe pas dans data_sources", c.DataSourceDefault)
}
```

👉 Remarques :

* J’ai corrigé ton tag YAML `yaml:datasource_default` → `yaml:"datasource_default"` (sinon ça ne sera pas bien parsé).
* Dans `Datasource`, tu avais deux fois `yaml:"name"`, j’ai corrigé en `yaml:"uri"` pour le champ `URI`.
* La méthode `Validate()` retourne une erreur si la clé par défaut n’est pas trouvée.

Exemple d’usage :

```go
cfg := AppConfig{
    Name: "MyApp",
    Port: 8080,
    Enabled: true,
    DataSourceDefault: "db1",
    Tags: []Datasource{
        {Name: "db1", URI: "postgres://..."},
        {Name: "db2", URI: "sqlite://..."},
    },
}

if err := cfg.Validate(); err != nil {
    panic(err) // ou log.Fatal
}
```

Veux-tu que j’intègre directement cette validation avec **go-playground/validator.v10** (puisque tu sembles déjà l’utiliser avec `validate:"..."`), ou tu préfères garder cette méthode séparée ?

Oui bien sûr 🙂

Actuellement je l’ai faite retourner une `error` car c’est l’usage courant en Go pour la validation, mais tu peux très bien la simplifier pour qu’elle retourne juste un `bool`.

Voici la variante :

```go
// Validate vérifie que DataSourceDefault existe dans Tags.Name
func (c *AppConfig) Validate() bool {
	if c.DataSourceDefault == "" {
		return false
	}

	for _, ds := range c.Tags {
		if ds.Name == c.DataSourceDefault {
			return true // trouvé ✅
		}
	}
	return false
}
```

👉 Dans ce cas, `true` = valide, `false` = invalide.

⚖️ Mais attention :

* Si tu veux donner plus d’infos à l’appelant (ex : "la datasource par défaut n’existe pas"), un `error` est plus adapté.
* Si tu veux juste savoir rapidement si c’est valide ou pas, un `bool` suffit.

Tu veux que je t’écrive une version qui **supporte les deux** (par exemple `ValidateBool()` et `ValidateError()`) pour que tu puisses choisir selon ton contexte ?

Parfait 👍 ajoutons un petit test unitaire Go pour la méthode `Validate()`.
Je vais te montrer les deux versions : si ta méthode retourne `bool`, et si elle retourne `error`.

---

### ✅ Version avec `bool`

```go
package config

import "testing"

func TestAppConfig_Validate(t *testing.T) {
	tests := []struct {
		name     string
		config   AppConfig
		expected bool
	}{
		{
			name: "datasource par défaut valide",
			config: AppConfig{
				DataSourceDefault: "db1",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
					{Name: "db2", URI: "sqlite://..."},
				},
			},
			expected: true,
		},
		{
			name: "datasource par défaut inexistante",
			config: AppConfig{
				DataSourceDefault: "db3",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
					{Name: "db2", URI: "sqlite://..."},
				},
			},
			expected: false,
		},
		{
			name: "datasource par défaut vide",
			config: AppConfig{
				DataSourceDefault: "",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.Validate(); got != tt.expected {
				t.Errorf("Validate() = %v, attendu %v", got, tt.expected)
			}
		})
	}
}
```

---

### ✅ Version avec `error`

Si tu gardes la version `Validate() error`, le test change légèrement :

```go
package config

import "testing"

func TestAppConfig_Validate(t *testing.T) {
	tests := []struct {
		name       string
		config     AppConfig
		shouldPass bool
	}{
		{
			name: "datasource par défaut valide",
			config: AppConfig{
				DataSourceDefault: "db1",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
					{Name: "db2", URI: "sqlite://..."},
				},
			},
			shouldPass: true,
		},
		{
			name: "datasource par défaut inexistante",
			config: AppConfig{
				DataSourceDefault: "db3",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
					{Name: "db2", URI: "sqlite://..."},
				},
			},
			shouldPass: false,
		},
		{
			name: "datasource par défaut vide",
			config: AppConfig{
				DataSourceDefault: "",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
				},
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.shouldPass && err != nil {
				t.Errorf("Validate() a échoué: %v", err)
			}
			if !tt.shouldPass && err == nil {
				t.Errorf("Validate() aurait dû échouer")
			}
		})
	}
}
```

---

👉 Tu veux que je parte sur la version **bool** (plus simple) ou **error** (plus informative) dans ton projet ?

Parfait 👍 voici un exemple de test unitaire pour la version `Validate() error` :

```go
package config

import "testing"

func TestAppConfig_Validate(t *testing.T) {
	tests := []struct {
		name       string
		config     AppConfig
		shouldPass bool
	}{
		{
			name: "datasource par défaut valide",
			config: AppConfig{
				DataSourceDefault: "db1",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
					{Name: "db2", URI: "sqlite://..."},
				},
			},
			shouldPass: true,
		},
		{
			name: "datasource par défaut inexistante",
			config: AppConfig{
				DataSourceDefault: "db3",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
					{Name: "db2", URI: "sqlite://..."},
				},
			},
			shouldPass: false,
		},
		{
			name: "datasource par défaut vide",
			config: AppConfig{
				DataSourceDefault: "",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
				},
			},
			shouldPass: false,
		},
		{
			name: "aucune datasource définie",
			config: AppConfig{
				DataSourceDefault: "db1",
				Tags:              []Datasource{},
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.shouldPass && err != nil {
				t.Errorf("Validate() a échoué alors qu'il devait réussir : %v", err)
			}
			if !tt.shouldPass && err == nil {
				t.Errorf("Validate() aurait dû échouer mais a réussi")
			}
		})
	}
}
```

👉 Ce test couvre plusieurs cas :

* valeur par défaut correcte ✅
* valeur par défaut inexistante ❌
* valeur vide ❌
* aucune datasource définie ❌

Veux-tu que je t’ajoute aussi une **version table-driven avec `t.Helper()`** pour centraliser la logique de comparaison (plus lisible si tu ajoutes d’autres validations) ?

Super 👍 voilà une version **table-driven** améliorée, avec un helper pour réduire la duplication :

```go
package config

import "testing"

// helper pour vérifier le résultat attendu
func assertValidation(t *testing.T, cfg AppConfig, shouldPass bool) {
	t.Helper()

	err := cfg.Validate()
	if shouldPass && err != nil {
		t.Errorf("Validate() a échoué alors qu'il devait réussir : %v", err)
	}
	if !shouldPass && err == nil {
		t.Errorf("Validate() aurait dû échouer mais a réussi")
	}
}

func TestAppConfig_Validate(t *testing.T) {
	tests := []struct {
		name       string
		config     AppConfig
		shouldPass bool
	}{
		{
			name: "datasource par défaut valide",
			config: AppConfig{
				DataSourceDefault: "db1",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
					{Name: "db2", URI: "sqlite://..."},
				},
			},
			shouldPass: true,
		},
		{
			name: "datasource par défaut inexistante",
			config: AppConfig{
				DataSourceDefault: "db3",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
					{Name: "db2", URI: "sqlite://..."},
				},
			},
			shouldPass: false,
		},
		{
			name: "datasource par défaut vide",
			config: AppConfig{
				DataSourceDefault: "",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
				},
			},
			shouldPass: false,
		},
		{
			name: "aucune datasource définie",
			config: AppConfig{
				DataSourceDefault: "db1",
				Tags:              []Datasource{},
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValidation(t, tt.config, tt.shouldPass)
		})
	}
}
```

✅ Avantages :

* `assertValidation` centralise la logique → le code des tests est plus lisible.
* Ajout d’un `t.Helper()` pour que les erreurs pointent vers le bon endroit dans le test.
* Facile à étendre si tu ajoutes d’autres règles dans `Validate()`.

Veux-tu que je te propose aussi un test équivalent mais en utilisant **`stretchr/testify`** (`require` / `assert`) pour simplifier encore plus la lisibilité ?

Parfait 🙌 si tu utilises [`stretchr/testify`](https://github.com/stretchr/testify), tu peux simplifier encore davantage tes tests.

Voici la version avec **`require`** et **`assert`** :

```go
package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppConfig_Validate(t *testing.T) {
	tests := []struct {
		name       string
		config     AppConfig
		shouldPass bool
	}{
		{
			name: "datasource par défaut valide",
			config: AppConfig{
				DataSourceDefault: "db1",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
					{Name: "db2", URI: "sqlite://..."},
				},
			},
			shouldPass: true,
		},
		{
			name: "datasource par défaut inexistante",
			config: AppConfig{
				DataSourceDefault: "db3",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
					{Name: "db2", URI: "sqlite://..."},
				},
			},
			shouldPass: false,
		},
		{
			name: "datasource par défaut vide",
			config: AppConfig{
				DataSourceDefault: "",
				Tags: []Datasource{
					{Name: "db1", URI: "postgres://..."},
				},
			},
			shouldPass: false,
		},
		{
			name: "aucune datasource définie",
			config: AppConfig{
				DataSourceDefault: "db1",
				Tags:              []Datasource{},
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.shouldPass {
				require.NoError(t, err, "la validation aurait dû réussir")
			} else {
				require.Error(t, err, "la validation aurait dû échouer")
			}
		})
	}
}
```

👉 Ici :

* `require.NoError` arrête le test immédiatement si l’erreur n’est pas conforme.
* `require.Error` fait pareil mais dans le cas inverse.
* Tu peux remplacer par `assert.NoError` / `assert.Error` si tu veux que **tous les cas continuent même si un échoue**.

Veux-tu que je te propose aussi une **version mixte** (`require` pour bloquer les préconditions, `assert` pour vérifier les détails de l’erreur) afin de tester aussi le contenu exact du message d’erreur retourné par `Validate()` ?
