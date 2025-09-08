En Go (Golang), un **buffer** est une zone de mémoire temporaire utilisée pour stocker des données avant qu'elles ne soient traitées. Il est souvent utilisé pour améliorer les performances en évitant des opérations fréquentes de lecture/écriture directement sur une ressource coûteuse (comme un fichier, un réseau, etc.).

---

### Types de buffers courants en Go

1. **`bytes.Buffer`**

   * Situé dans le package `bytes`
   * Sert à construire ou manipuler dynamiquement des chaînes de caractères ou des tableaux de bytes.

   Exemple :

   ```go
   var buf bytes.Buffer
   buf.WriteString("Hello ")
   buf.WriteString("World!")
   fmt.Println(buf.String()) // Output: Hello World!
   ```

2. **`bufio` (Buffered I/O)**

   * Le package `bufio` fournit des buffers pour la lecture (`bufio.Reader`) ou l’écriture (`bufio.Writer`) au-dessus d’un flux comme un fichier ou une connexion réseau.
   * Il permet de lire ou écrire plus efficacement en limitant les appels système.

   Exemple d’écriture avec `bufio.Writer` :

   ```go
   file, _ := os.Create("example.txt")
   writer := bufio.NewWriter(file)
   writer.WriteString("Hello buffered world!")
   writer.Flush() // Nécessaire pour écrire dans le fichier
   ```

   Exemple de lecture avec `bufio.Reader` :

   ```go
   file, _ := os.Open("example.txt")
   reader := bufio.NewReader(file)
   line, _ := reader.ReadString('\n')
   fmt.Println(line)
   ```

---

### Pourquoi utiliser un buffer ?

* **Performance :** réduit les appels directs à des ressources lentes (disque, réseau…)
* **Flexibilité :** permet de manipuler des données en mémoire avant de les envoyer ou les sauvegarder
* **Utilitaire :** permet de construire ou parser facilement des données complexes

Souhaites-tu un exemple complet avec un buffer appliqué à une communication réseau ?


Voici un exemple simple et complet d'utilisation de **buffers avec une communication réseau en Go**, en utilisant `bufio` pour **gérer les lectures/écritures sur une connexion TCP**.

---

## 🔌 Exemple : Serveur + Client TCP avec `bufio`

### 🖥️ Serveur TCP (avec buffer)

```go
// server.go
package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	// Écoute sur le port 9000
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Println("Serveur en écoute sur le port 9000...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erreur de connexion :", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connexion fermée.")
			return
		}
		fmt.Print("Message reçu :", message)

		response := "Reçu : " + message
		writer.WriteString(response)
		writer.Flush() // Nécessaire pour envoyer le message
	}
}
```

---

### 🧑‍💻 Client TCP (avec buffer)

```go
// client.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)
	serverWriter := bufio.NewWriter(conn)

	for {
		fmt.Print("Entrez un message : ")
		text, _ := reader.ReadString('\n')

		// Écrire dans la connexion
		serverWriter.WriteString(text)
		serverWriter.Flush()

		// Lire la réponse
		response, _ := serverReader.ReadString('\n')
		fmt.Println("Réponse du serveur :", response)
	}
}
```

---

### 🧪 Résultat

* Le client envoie une ligne de texte.
* Le serveur la lit via un buffer (`bufio.Reader`), la traite et répond.
* La réponse est envoyée via un buffer (`bufio.Writer`) et reçue par le client.

---

Souhaites-tu qu’on ajoute un protocole simple ou un format (JSON, par exemple) pour les messages échangés ?
