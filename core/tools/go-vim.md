Bien sûr ! Pour t’aider efficacement, voici les différentes étapes pour configurer et utiliser **Go avec Vim** (souvent appelé **go-vim**), afin d'avoir une bonne expérience de développement Go dans Vim :

---

## 🔧 Prérequis

1. **Vim** (version 8+ recommandée, ou Neovim).
2. **Go** installé (`go version` pour vérifier).
3. **Plugin manager pour Vim** – comme [vim-plug](https://github.com/junegunn/vim-plug), [Vundle](https://github.com/VundleVim/Vundle.vim) ou autre.

---

## 🚀 Étapes pour configurer go-vim avec `vim-plug`

1. **Installer `vim-plug`** (si pas déjà fait) :

   ```bash
   curl -fLo ~/.vim/autoload/plug.vim --create-dirs \
       https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
   ```

2. **Modifier ton fichier `~/.vimrc`** :

   ```vim
   call plug#begin('~/.vim/plugged')

   " Go support
   Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }

   call plug#end()
   ```

3. **Installer les plugins dans Vim** :
   Lance Vim, puis tape :

   ```vim
   :PlugInstall
   ```

4. **Installer les outils Go recommandés par vim-go** :
   Une fois les plugins installés, exécute :

   ```vim
   :GoInstallBinaries
   ```

---

## ✨ Fonctionnalités apportées par `vim-go`

* `:GoBuild`, `:GoRun`, `:GoTest`
* Autocomplétion (si couplé avec un moteur comme coc.nvim ou nvim-cmp)
* Formatage automatique (`gofmt`, `goimports`)
* Affichage des erreurs (`:GoLint`)
* Navigation rapide (`:GoDef`, `:GoDoc`, `:GoInfo`)
* Test avec `:GoTestFunc`, etc.

---

## ✅ Conseils utiles

* Active le **formatage à la sauvegarde** :

  ```vim
  let g:go_fmt_command = "goimports"
  autocmd BufWritePre *.go :silent! lua vim.lsp.buf.format()
  ```
* Pour une expérience plus moderne, tu peux aussi essayer **Neovim + LSP (gopls)** au lieu de `vim-go`.

---

Voici un **aide-mémoire (cheat sheet)** pour utiliser **Go avec Vim (`vim-go`)** :
Tu peux l’enregistrer dans un fichier `go-vim-cheatsheet.md` ou l'imprimer.

---

# 🧠 Go + Vim (vim-go) – Aide-mémoire

## ⚙️ Commandes de base

| Commande Vim  | Action                                                 |
| ------------- | ------------------------------------------------------ |
| `:GoRun`      | Exécute le fichier Go courant                          |
| `:GoBuild`    | Compile le fichier Go courant                          |
| `:GoTest`     | Exécute tous les tests du package                      |
| `:GoTestFunc` | Exécute uniquement le test sous le curseur             |
| `:GoVet`      | Analyse statique de code                               |
| `:GoLint`     | Linting avec `golint`                                  |
| `:GoFmt`      | Format avec `gofmt`                                    |
| `:GoImports`  | Format avec `goimports` (importe/supprime les imports) |
| `:GoModTidy`  | Nettoie et corrige les dépendances du `go.mod`         |

---

## 🧭 Navigation

| Commande Vim  | Action                                              |
| ------------- | --------------------------------------------------- |
| `:GoDef`      | Aller à la définition sous le curseur               |
| `:GoDecls`    | Liste des déclarations dans le fichier courant      |
| `:GoDoc`      | Documentation de l'élément sous le curseur          |
| `:GoDescribe` | Description de l’objet sous le curseur              |
| `:GoInfo`     | Affiche le type et info sur l'objet sous le curseur |
| `K`           | Affiche la documentation (équivalent à `:GoDoc`)    |

---

## 🔧 Débogage et diagnostics

| Commande Vim         | Action                                       |
| -------------------- | -------------------------------------------- |
| `:GoErrCheck`        | Vérifie les erreurs non traitées             |
| `:GoMetaLinter`      | Lance plusieurs linters en parallèle         |
| `:GoInstallBinaries` | Installe tous les outils requis par `vim-go` |

---

## 📦 Modules & Dépendances

| Commande Vim | Action                                                |
| ------------ | ----------------------------------------------------- |
| `:GoModInit` | Initialise un module Go (`go.mod`)                    |
| `:GoModTidy` | Supprime les imports inutilisés, ajoute les manquants |
| `:GoModWhy`  | Pourquoi une dépendance est là                        |

---

## 🧪 Tests

| Commande Vim  | Action                          |
| ------------- | ------------------------------- |
| `:GoTest`     | Lance tous les tests du package |
| `:GoTestFunc` | Lance le test sous le curseur   |
| `:GoCoverage` | Affiche la couverture de tests  |

---

## 🔁 Commandes utiles à ajouter à `.vimrc`

```vim
let g:go_fmt_command = "goimports"      " Utilise goimports à la place de gofmt
autocmd FileType go nmap <Leader>b :GoBuild<CR>
autocmd FileType go nmap <Leader>r :GoRun<CR>
autocmd FileType go nmap <Leader>t :GoTest<CR>
autocmd BufWritePre *.go :GoImports    " Formatte automatiquement à la sauvegarde
```

---

Souhaites-tu ce mémo au format PDF ou Markdown ?


## cheat sheet
# 🧠 Go + Vim (`vim-go`) – Aide-mémoire

## ⚙️ Commandes de base

| Commande Vim        | Action                                                    |
|---------------------|-----------------------------------------------------------|
| `:GoRun`            | Exécute le fichier Go courant                             |
| `:GoBuild`          | Compile le fichier Go courant                             |
| `:GoTest`           | Exécute tous les tests du package                         |
| `:GoTestFunc`       | Exécute uniquement le test sous le curseur                |
| `:GoVet`            | Analyse statique de code                                  |
| `:GoLint`           | Linting avec `golint`                                     |
| `:GoFmt`            | Format avec `gofmt`                                       |
| `:GoImports`        | Format avec `goimports` (importe/supprime les imports)    |
| `:GoModTidy`        | Nettoie et corrige les dépendances du `go.mod`            |

---

## 🧭 Navigation

| Commande Vim        | Action                                                    |
|---------------------|-----------------------------------------------------------|
| `:GoDef`            | Aller à la définition sous le curseur                     |
| `:GoDecls`          | Liste des déclarations dans le fichier courant            |
| `:GoDoc`            | Documentation de l'élément sous le curseur                |
| `:GoDescribe`       | Description de l’objet sous le curseur                    |
| `:GoInfo`           | Affiche le type et info sur l'objet sous le curseur       |
| `K`                 | Affiche la documentation (équivalent à `:GoDoc`)          |

---

## 🔧 Débogage et diagnostics

| Commande Vim        | Action                                                    |
|---------------------|-----------------------------------------------------------|
| `:GoErrCheck`       | Vérifie les erreurs non traitées                          |
| `:GoMetaLinter`     | Lance plusieurs linters en parallèle                      |
| `:GoInstallBinaries`| Installe tous les outils requis par `vim-go`              |

---

## 📦 Modules & Dépendances

| Commande Vim        | Action                                                    |
|---------------------|-----------------------------------------------------------|
| `:GoModInit`        | Initialise un module Go (`go.mod`)                        |
| `:GoModTidy`        | Supprime les imports inutilisés, ajoute les manquants     |
| `:GoModWhy`         | Explique pourquoi une dépendance est utilisée             |

---

## 🧪 Tests

| Commande Vim        | Action                                                    |
|---------------------|-----------------------------------------------------------|
| `:GoTest`           | Lance tous les tests du package                           |
| `:GoTestFunc`       | Lance le test sous le curseur                             |
| `:GoCoverage`       | Affiche la couverture de tests                            |

---

## 🔁 Extrait utile pour `.vimrc`

```vim
let g:go_fmt_command = "goimports"
autocmd FileType go nmap <Leader>b :GoBuild<CR>
autocmd FileType go nmap <Leader>r :GoRun<CR>
autocmd FileType go nmap <Leader>t :GoTest<CR>
autocmd BufWritePre *.go :GoImports
```

