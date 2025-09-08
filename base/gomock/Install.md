Pour installer **GoMock** et son outil de génération **mockgen** sur macOS, voici les étapes à suivre :

---

### ✅ Prérequis

Assure-toi que Go est déjà installé :

```bash
go version
```

---

### 1. 📦 Installer `GoMock` (librairie)

Aucune installation spéciale n’est requise pour la bibliothèque elle-même : elle s’installe automatiquement avec `go get` ou lors du `go mod tidy`.

Tu peux l’ajouter à ton projet avec :

```bash
go get go.uber.org/mock
```

---

### 2. 🔧 Installer `mockgen` (outil CLI)

#### Option 1 : via `go install` (recommandé)

```bash
go install go.uber.org/mock/mockgen@latest
```

Cela installe l’exécutable `mockgen` dans ton `$GOPATH/bin` ou `$HOME/go/bin`.

#### Vérifie l’installation :

```bash
mockgen --version
```

Si tu obtiens une erreur de commande introuvable, ajoute `$HOME/go/bin` à ton `PATH` :

```bash
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
source ~/.zshrc
```

(ou `~/.bash_profile` si tu utilises bash)

---

### 3. 📄 Générer un mock avec `mockgen`

Voici un exemple de commande :

```bash
mockgen -source=your_interface.go -destination=your_interface_mock.go -package=yourpackage
```

---
