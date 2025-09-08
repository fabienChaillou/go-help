Si tu cherches un exemple simple d’utilisation de **flags** (ou **"flaog"** si c'était une faute de frappe) en **Go**, voici un petit exemple utilisant le package standard [`flag`](https://pkg.go.dev/flag) pour parser des arguments en ligne de commande :

### Exemple basique d'utilisation de `flag` en Go

```go
package main

import (
	"flag"
	"fmt"
)

func main() {
	// Déclaration des flags
	name := flag.String("name", "World", "Votre nom")
	age := flag.Int("age", 30, "Votre âge")
	debug := flag.Bool("debug", false, "Mode debug")

	// Parse les flags
	flag.Parse()

	// Utilisation des valeurs
	fmt.Printf("Bonjour %s !\n", *name)
	fmt.Printf("Vous avez %d ans.\n", *age)

	if *debug {
		fmt.Println("Mode debug activé.")
	}
}
```

### Exécution

Tu peux compiler ce fichier (nommé par exemple `main.go`) et l’exécuter comme ceci :

```bash
go run main.go -name=Alice -age=25 -debug
```

### Résultat :

```
Bonjour Alice !
Vous avez 25 ans.
Mode debug activé.
```

Souhaites-tu un exemple plus avancé (comme avec des sous-commandes ou des flags obligatoires) ?
