En Go, le **buffer** fait généralement référence à un espace mémoire temporaire utilisé pour **stocker des données avant traitement ou transmission**. Le concept de buffer est souvent utilisé avec :

1. **Les channels bufferisés** (concurrence)
2. **Les buffers d'entrée/sortie** (fichiers, réseaux, etc.)

---

## 1. 🔹 Buffer dans les **channels**

Un **channel bufferisé** permet de **stocker temporairement plusieurs valeurs**, contrairement à un channel non bufferisé (qui bloque jusqu’à réception).

### ➤ Création d’un channel bufferisé :

```go
ch := make(chan int, 3) // channel de type int avec un buffer de 3
```

### ➤ Fonctionnement :

* Tu peux envoyer jusqu’à 3 valeurs dans le channel sans que personne ne les lise.
* Au 4e envoi, l'envoi **bloque** jusqu'à ce qu'une valeur soit lue.

### ➤ Exemple :

```go
func main() {
    ch := make(chan int, 2)

    ch <- 1
    ch <- 2
    // ch <- 3 // ceci bloquerait car le buffer est plein

    fmt.Println(<-ch) // lit 1
    fmt.Println(<-ch) // lit 2
}
```

---

## 2. 🔹 Buffer en **entrée/sortie (I/O)**

Go fournit des buffers pour optimiser les lectures/écritures sur les fichiers ou réseaux, via le package `bufio`.

### ➤ Exemple : écriture bufferisée dans un fichier

```go
import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    fichier, _ := os.Create("exemple.txt")
    writer := bufio.NewWriter(fichier)

    writer.WriteString("Hello, monde !\n") // écrit dans le buffer, pas directement dans le fichier
    writer.Flush()                        // vide le buffer dans le fichier

    fichier.Close()
}
```

🔸 `bufio.NewWriter` crée un buffer d’écriture.
🔸 `Flush()` est **obligatoire** pour forcer l’écriture dans le fichier.

---

## 🔹 Pourquoi utiliser des buffers ?

* ✅ **Meilleure performance** : moins d'accès directs au disque ou aux threads.
* ✅ **Moins de blocage** avec des channels bufferisés.
* ⚠️ Attention à bien **vider les buffers** (`Flush()`), sinon des données peuvent rester en mémoire.

---

Souhaites-tu un exemple comparatif entre channel bufferisé et non bufferisé dans un contexte de producteur/consommateur ?
