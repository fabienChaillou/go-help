L’argument `-gcflags` dans Go (Golang) est utilisé avec la commande `go build`, `go run`, ou `go test` pour passer des options spécifiques au **compilateur Go (le "gc", ou Go compiler)**. Cet argument est principalement utilisé pour le **debug**, le **profiling**, ou le **tuning des performances** du code compilé.

---

### Syntaxe de base

```bash
go build -gcflags="options"
```

Ou pour des fichiers spécifiques :

```bash
go build -gcflags="pattern=options"
```

---

### Exemples courants

#### 1. **Désactiver l'inlining**

Empêche le compilateur d’optimiser certaines fonctions en les « inlinant » (copier le code à l’endroit de l’appel au lieu d’un appel de fonction).

```bash
go build -gcflags="all=-l"
```

* `-l` : désactive l'inlining
* `all=` : applique à tous les packages

#### 2. **Désactiver l'optimisation globale**

Utile pour le debug, surtout quand on veut que le debugger (comme `delve`) ait une vision fidèle du code.

```bash
go build -gcflags="all=-N -l"
```

* `-N` : désactive les optimisations
* `-l` : désactive l’inlining

👉 Très utilisé avec Delve pour déboguer :

```bash
dlv debug -- -gcflags="all=-N -l"
```

#### 3. **Options ciblées pour un package**

Pour appliquer une option de compilation uniquement à un package donné :

```bash
go build -gcflags="mypkg=-N -l"
```

---

### Cas d'usage typiques

| Objectif                         | Option `-gcflags` utilisée                             |
| -------------------------------- | ------------------------------------------------------ |
| Debug avec Delve                 | `-gcflags="all=-N -l"`                                 |
| Analyser les performances        | `-gcflags="all=-m"` (affiche les décisions d'inlining) |
| Empêcher certaines optimisations | `-gcflags="all=-N"` ou `-gcflags="all=-N -l"`          |

---

### À noter

* `-gcflags` est destiné à des **usages avancés** (debug, test bas niveau, développement du compilateur, etc.)
* Mal utilisé, il peut rendre le binaire plus lent ou incorrect pour la production.

Souhaites-tu un exemple concret avec compilation et effet visible d’un `-gcflags` ?

Parfait ! Voici un exemple concret pour illustrer l'effet de `-gcflags="all=-m"` (affiche les décisions d’inlining du compilateur).

---

## 🔧 Exemple : Observer les décisions d'inlining

### 1. Crée un fichier Go

`main.go`

```go
package main

import "fmt"

func sayHello() {
	fmt.Println("Hello, world!")
}

func main() {
	sayHello()
}
```

---

### 2. Compile avec `-gcflags="all=-m"`

```bash
go build -gcflags="all=-m" main.go
```

### 3. Résultat attendu

Tu verras une sortie comme celle-ci :

```txt
# command-line-arguments
./main.go:6:6: can inline sayHello
./main.go:10:2: inlining call to sayHello
```

🧠 **Explication** :

* `can inline sayHello` : le compilateur estime que la fonction est simple et peut être "inlinée"
* `inlining call to sayHello` : la fonction a effectivement été inlinée dans `main()`, donc le code est copié au lieu d’un appel de fonction

---

### 4. Désactiver l'inlining pour observer la différence

```bash
go build -gcflags="all=-l -m" main.go
```

Résultat :

```txt
# command-line-arguments
./main.go:6:6: cannot inline sayHello: inlining disabled
```

✅ **Tu viens de voir l’effet réel du flag `-l` et comment `-m` t’informe des décisions du compilateur.**

---

Souhaites-tu un exemple avec `-N` ou un débogage avec Delve également ?

Super ! Passons maintenant à un **exemple avec `-N`** (désactive les optimisations) et **`delve` (le debugger Go)** pour voir comment cela facilite le débogage.

---

## 🎯 Objectif : Déboguer un programme Go avec Delve sans optimisations

### 1. Code source : `main.go`

```go
package main

import "fmt"

func compute(x int) int {
	y := x * 2
	return y + 3
}

func main() {
	result := compute(5)
	fmt.Println("Result:", result)
}
```

---

### 2. Lancer le débogueur avec `-gcflags="all=-N -l"`

> Cette commande désactive l'inlining (`-l`) et les optimisations (`-N`) afin que le debugger fonctionne correctement.

```bash
dlv debug -- -gcflags="all=-N -l"
```

> Cela lance `dlv` (Delve) en mode interactif sur ton code compilé avec les bonnes options.

---

### 3. Dans Delve, tu peux faire :

```bash
(dlv) break main.compute
Breakpoint 1 set at 0x... for main.compute
(dlv) run
(dlv) print x
5
(dlv) next
(dlv) print y
10
(dlv) continue
```

💡 Sans `-N -l`, `compute()` pourrait être **inlinée**, rendant impossible la mise en pause dedans ou l'inspection de ses variables.

---

### ⚠️ Résumé de l'utilité de `-gcflags="all=-N -l"` avec Delve

| Option   | Effet                                                                                           |
| -------- | ----------------------------------------------------------------------------------------------- |
| `-N`     | Désactive les optimisations (ex: suppression de variables inutilisées, fusion de boucles, etc.) |
| `-l`     | Empêche l'inlining des fonctions pour que tu puisses mettre des points d'arrêt à l’intérieur    |
| Ensemble | Donne une vision fidèle du code, idéale pour le **pas à pas et l’inspection de variables**      |

---

Souhaites-tu que je t'aide à installer et utiliser `dlv` si ce n'est pas encore fait sur ta machine ?
