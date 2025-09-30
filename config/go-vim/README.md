Tu veux probablement parler de **l’intégration de Go avec Vim/Neovim**, souvent appelée *vim-go* (le plugin principal pour Go dans Vim), ou bien de l’utilisation de Go avec une config Vim adaptée.
👉 Il n’existe pas de projet officiel appelé **go-vim**, mais plutôt **vim-go** (par [fatih/vim-go](https://github.com/fatih/vim-go)).

Voici comment installer et utiliser **vim-go** :

---

## 1. Pré-requis

* Avoir **Go** installé (>= 1.18 recommandé)

  ```bash
  go version
  ```
* Avoir un gestionnaire de plugins Vim/Neovim (ex. `vim-plug`, `packer.nvim`, `dein.vim`, etc.).

---

## 2. Installation avec vim-plug (exemple)

Ajoute ceci à ton `~/.vimrc` ou `~/.config/nvim/init.vim` :

```vim
call plug#begin('~/.vim/plugged')

Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }

call plug#end()
```

Ensuite dans Vim/Neovim :

```vim
:PlugInstall
```

Cela installera **vim-go** et téléchargera les outils Go nécessaires.

---

## 3. Mise à jour des outils Go

Après installation, lance :

```vim
:GoInstallBinaries
```

Cela va installer les outils utiles (formatters, linters, gopls, etc.).

---

## 4. Fonctionnalités principales

* **Compilation & test**

  ```vim
  :GoBuild
  :GoRun
  :GoTest
  ```
* **Navigation dans le code**

  * `gd` → aller à la définition
  * `gr` → références
* **Auto-format & import**
  (par défaut avec `gofmt` ou `goimports`)

  ```vim
  :GoFmt
  :GoImports
  ```
* **Lint**

  ```vim
  :GoLint
  ```
* **Documentation rapide**

  ```vim
  :GoDoc
  ```

---

## 5. Utiliser LSP (optionnel mais recommandé)

Aujourd’hui, beaucoup préfèrent utiliser `gopls` (LSP officiel de Go) avec un client LSP dans Vim/Neovim (`coc.nvim`, `nvim-lspconfig`, etc.) plutôt que vim-go seul.

Exemple (Neovim + `nvim-lspconfig`) :

```lua
require('lspconfig').gopls.setup{}
```

---

👉 Question : tu veux que je t’explique **la config simple avec juste vim-go**, ou plutôt **une config moderne avec LSP (gopls + autocomplétion + treesitter)** pour travailler confortablement en Go ?

Parfait 👍 tu as raison : **vim-go** était incontournable il y a quelques années, mais aujourd’hui le workflow moderne en Go sous Vim/Neovim repose plutôt sur :

* **`gopls`** → serveur LSP officiel pour Go (autocomplétion, go-to-def, refactoring, diagnostics).
* **Treesitter** → coloration et parsing syntaxique avancé.
* **Mason** ou **nvim-lsp-installer** → gestion facile des serveurs LSP et outils.
* **nvim-cmp** → moteur d’autocomplétion moderne.
* (optionnel) **null-ls** → intégration de linters/formatters.

---

## 🔧 Exemple de config moderne (Neovim, en Lua)

### 1. Installer les plugins (avec `lazy.nvim` ou `packer.nvim`)

Exemple avec **packer.nvim** (`~/.config/nvim/init.lua`) :

```lua
require('packer').startup(function(use)
  -- Gestionnaire de plugins
  use 'wbthomason/packer.nvim'

  -- LSP
  use 'neovim/nvim-lspconfig'
  use 'williamboman/mason.nvim'
  use 'williamboman/mason-lspconfig.nvim'

  -- Autocomplétion
  use 'hrsh7th/nvim-cmp'
  use 'hrsh7th/cmp-nvim-lsp'
  use 'hrsh7th/cmp-buffer'
  use 'hrsh7th/cmp-path'
  use 'L3MON4D3/LuaSnip'

  -- Treesitter
  use { 'nvim-treesitter/nvim-treesitter', run = ':TSUpdate' }

  -- (optionnel) Formatter/Linter
  use 'jose-elias-alvarez/null-ls.nvim'
end)
```

Ensuite, fais :

```bash
:PackerSync
```

---

### 2. Config LSP + Go (`gopls`)

Dans `~/.config/nvim/init.lua` ou `lua/config/lsp.lua` :

```lua
-- Mason (installe gopls automatiquement)
require("mason").setup()
require("mason-lspconfig").setup({
  ensure_installed = { "gopls" }
})

-- Setup LSP
local lspconfig = require("lspconfig")

lspconfig.gopls.setup {
  capabilities = require('cmp_nvim_lsp').default_capabilities(),
}

-- Keymaps pratiques
vim.keymap.set("n", "gd", vim.lsp.buf.definition, {})
vim.keymap.set("n", "gr", vim.lsp.buf.references, {})
vim.keymap.set("n", "K", vim.lsp.buf.hover, {})
vim.keymap.set("n", "<leader>rn", vim.lsp.buf.rename, {})
```

---

### 3. Config autocomplétion (nvim-cmp)

Toujours dans `init.lua` :

```lua
local cmp = require'cmp'

cmp.setup({
  snippet = {
    expand = function(args)
      require('luasnip').lsp_expand(args.body)
    end,
  },
  mapping = cmp.mapping.preset.insert({
    ['<C-Space>'] = cmp.mapping.complete(),
    ['<CR>'] = cmp.mapping.confirm({ select = true }),
  }),
  sources = cmp.config.sources({
    { name = 'nvim_lsp' },
    { name = 'buffer' },
    { name = 'path' },
  }),
})
```

---

### 4. Treesitter pour Go

```lua
require'nvim-treesitter.configs'.setup {
  ensure_installed = { "go", "gomod", "lua", "json" },
  highlight = { enable = true },
}
```

---

### 5. Formatter & Linter (optionnel)

Avec **null-ls**, tu peux utiliser `goimports` ou `golangci-lint` :

```lua
local null_ls = require("null-ls")

null_ls.setup({
  sources = {
    null_ls.builtins.formatting.goimports,
    null_ls.builtins.diagnostics.golangci_lint,
  },
})
```

Ensuite, dans Neovim :

```vim
:MasonInstall gopls goimports golangci-lint
```

---

## 🚀 Résultat

Avec cette config moderne tu auras :

* Complétion intelligente (LSP + snippets).
* Navigation (`gd`, `gr`, `K`, renommage).
* Highlight et parsing via Treesitter.
* Formatage auto avec `goimports`.
* Diagnostics & lint intégrés.

---

👉 Veux-tu que je te prépare un **`init.lua` complet prêt à copier-coller** pour démarrer direct, ou tu préfères une explication par modules (LSP, cmp, treesitter séparés) ?

Super 🚀 Je vais te préparer un `init.lua` **complet et prêt à l’emploi** pour Neovim, avec tout ce qu’il faut pour coder confortablement en Go (LSP, autocomplétion, Treesitter, formatage).

---

## 📂 Arborescence

Tu auras quelque chose comme :

```
~/.config/nvim/
└── init.lua
```

---

## 📜 `init.lua` complet

```lua
-- ================================
--   Neovim config pour Go moderne
-- ================================

-- === 1. Bootstrap packer si besoin ===
local fn = vim.fn
local install_path = fn.stdpath('data')..'/site/pack/packer/start/packer.nvim'
if fn.empty(fn.glob(install_path)) > 0 then
  fn.system({'git', 'clone', '--depth', '1',
    'https://github.com/wbthomason/packer.nvim', install_path})
  vim.cmd [[packadd packer.nvim]]
end

-- === 2. Plugins ===
require('packer').startup(function(use)
  use 'wbthomason/packer.nvim'

  -- LSP
  use 'neovim/nvim-lspconfig'
  use 'williamboman/mason.nvim'
  use 'williamboman/mason-lspconfig.nvim'

  -- Autocomplétion
  use 'hrsh7th/nvim-cmp'
  use 'hrsh7th/cmp-nvim-lsp'
  use 'hrsh7th/cmp-buffer'
  use 'hrsh7th/cmp-path'
  use 'L3MON4D3/LuaSnip'

  -- Treesitter
  use { 'nvim-treesitter/nvim-treesitter', run = ':TSUpdate' }

  -- Format/Lint
  use 'jose-elias-alvarez/null-ls.nvim'
end)

-- === 3. Mason setup (gestion LSP/formatters) ===
require("mason").setup()
require("mason-lspconfig").setup({
  ensure_installed = { "gopls" },
})

-- === 4. LSP config (Go via gopls) ===
local lspconfig = require("lspconfig")
local capabilities = require('cmp_nvim_lsp').default_capabilities()

lspconfig.gopls.setup {
  capabilities = capabilities,
}

-- Keymaps pratiques LSP
local opts = { noremap=true, silent=true }
vim.keymap.set("n", "gd", vim.lsp.buf.definition, opts)
vim.keymap.set("n", "gr", vim.lsp.buf.references, opts)
vim.keymap.set("n", "K", vim.lsp.buf.hover, opts)
vim.keymap.set("n", "<leader>rn", vim.lsp.buf.rename, opts)
vim.keymap.set("n", "<leader>ca", vim.lsp.buf.code_action, opts)

-- === 5. nvim-cmp (autocomplétion) ===
local cmp = require'cmp'
cmp.setup({
  snippet = {
    expand = function(args)
      require('luasnip').lsp_expand(args.body)
    end,
  },
  mapping = cmp.mapping.preset.insert({
    ['<C-Space>'] = cmp.mapping.complete(),
    ['<CR>'] = cmp.mapping.confirm({ select = true }),
  }),
  sources = cmp.config.sources({
    { name = 'nvim_lsp' },
    { name = 'buffer' },
    { name = 'path' },
  }),
})

-- === 6. Treesitter ===
require'nvim-treesitter.configs'.setup {
  ensure_installed = { "go", "gomod", "lua", "json" },
  highlight = { enable = true },
}

-- === 7. null-ls (formatter/linter) ===
local null_ls = require("null-ls")
null_ls.setup({
  sources = {
    null_ls.builtins.formatting.goimports,
    null_ls.builtins.diagnostics.golangci_lint,
  },
})

-- === 8. Options pratiques ===
vim.o.number = true
vim.o.relativenumber = true
vim.o.expandtab = true
vim.o.shiftwidth = 2
vim.o.tabstop = 2
vim.o.termguicolors = true
vim.o.updatetime = 250

print("✅ Neovim Go config chargée !")
```

---

## 🚀 Étapes pour démarrer

1. Crée `~/.config/nvim/init.lua` avec ce contenu.
2. Ouvre Neovim et fais :

   ```vim
   :PackerSync
   ```

   (ça installera les plugins).
3. Installe les outils Go via Mason :

   ```vim
   :MasonInstall gopls goimports golangci-lint
   ```
4. Ouvre un fichier `.go` → tu auras :

   * Autocomplétion LSP
   * Navigation (`gd`, `gr`, `K`)
   * Format auto avec `goimports` (`:lua vim.lsp.buf.format()`)
   * Diagnostics intégrés

---

👉 Veux-tu que je t’ajoute aussi la **sauvegarde auto avec formatage Go** (genre `:w` déclenche `goimports` automatiquement) ?

Parfait 👌 on va ajouter le **formatage automatique à l’enregistrement** dans ton `init.lua`.
Ça évite d’oublier de lancer `goimports` et ça garde ton code Go toujours propre.

---

## 📜 Ajout à ton `init.lua`

Ajoute ce bloc **après la config `null-ls`** :

```lua
-- === 9. Format auto à la sauvegarde ===
vim.api.nvim_create_autocmd("BufWritePre", {
  pattern = "*.go",
  callback = function()
    vim.lsp.buf.format({ async = false })
  end,
})
```

---

## 🔧 Résultat

* À chaque `:w`, Neovim lance `gopls`/`goimports` pour reformater ton fichier.
* Ça gère aussi les imports manquants (grâce à `goimports`).
* Si `golangci-lint` est installé via Mason, tu auras aussi les diagnostics lint automatiquement.

---

👉 Veux-tu que je t’intègre ça directement dans le `init.lua` **complet final** (avec la section 9 incluse), comme ça tu n’auras qu’à copier-coller un seul fichier ?

Parfait 🙌 voici la version **complète et finale de `init.lua`**, avec le formatage automatique Go intégré. Tu peux copier-coller directement ce fichier dans `~/.config/nvim/init.lua`.

---

## 📜 `init.lua`

```lua
-- ================================
--   Neovim config pour Go moderne
-- ================================

-- === 1. Bootstrap packer si besoin ===
local fn = vim.fn
local install_path = fn.stdpath('data')..'/site/pack/packer/start/packer.nvim'
if fn.empty(fn.glob(install_path)) > 0 then
  fn.system({'git', 'clone', '--depth', '1',
    'https://github.com/wbthomason/packer.nvim', install_path})
  vim.cmd [[packadd packer.nvim]]
end

-- === 2. Plugins ===
require('packer').startup(function(use)
  use 'wbthomason/packer.nvim'

  -- LSP
  use 'neovim/nvim-lspconfig'
  use 'williamboman/mason.nvim'
  use 'williamboman/mason-lspconfig.nvim'

  -- Autocomplétion
  use 'hrsh7th/nvim-cmp'
  use 'hrsh7th/cmp-nvim-lsp'
  use 'hrsh7th/cmp-buffer'
  use 'hrsh7th/cmp-path'
  use 'L3MON4D3/LuaSnip'

  -- Treesitter
  use { 'nvim-treesitter/nvim-treesitter', run = ':TSUpdate' }

  -- Format/Lint
  use 'jose-elias-alvarez/null-ls.nvim'
end)

-- === 3. Mason setup (gestion LSP/formatters) ===
require("mason").setup()
require("mason-lspconfig").setup({
  ensure_installed = { "gopls" },
})

-- === 4. LSP config (Go via gopls) ===
local lspconfig = require("lspconfig")
local capabilities = require('cmp_nvim_lsp').default_capabilities()

lspconfig.gopls.setup {
  capabilities = capabilities,
}

-- Keymaps pratiques LSP
local opts = { noremap=true, silent=true }
vim.keymap.set("n", "gd", vim.lsp.buf.definition, opts)
vim.keymap.set("n", "gr", vim.lsp.buf.references, opts)
vim.keymap.set("n", "K", vim.lsp.buf.hover, opts)
vim.keymap.set("n", "<leader>rn", vim.lsp.buf.rename, opts)
vim.keymap.set("n", "<leader>ca", vim.lsp.buf.code_action, opts)

-- === 5. nvim-cmp (autocomplétion) ===
local cmp = require'cmp'
cmp.setup({
  snippet = {
    expand = function(args)
      require('luasnip').lsp_expand(args.body)
    end,
  },
  mapping = cmp.mapping.preset.insert({
    ['<C-Space>'] = cmp.mapping.complete(),
    ['<CR>'] = cmp.mapping.confirm({ select = true }),
  }),
  sources = cmp.config.sources({
    { name = 'nvim_lsp' },
    { name = 'buffer' },
    { name = 'path' },
  }),
})

-- === 6. Treesitter ===
require'nvim-treesitter.configs'.setup {
  ensure_installed = { "go", "gomod", "lua", "json" },
  highlight = { enable = true },
}

-- === 7. null-ls (formatter/linter) ===
local null_ls = require("null-ls")
null_ls.setup({
  sources = {
    null_ls.builtins.formatting.goimports,
    null_ls.builtins.diagnostics.golangci_lint,
  },
})

-- === 8. Format auto à la sauvegarde ===
vim.api.nvim_create_autocmd("BufWritePre", {
  pattern = "*.go",
  callback = function()
    vim.lsp.buf.format({ async = false })
  end,
})

-- === 9. Options pratiques ===
vim.o.number = true
vim.o.relativenumber = true
vim.o.expandtab = true
vim.o.shiftwidth = 2
vim.o.tabstop = 2
vim.o.termguicolors = true
vim.o.updatetime = 250

print("✅ Neovim Go config chargée !")
```

---

## 🚀 Étapes pour installer

1. Mets ce fichier dans `~/.config/nvim/init.lua`.
2. Ouvre Neovim et fais :

   ```vim
   :PackerSync
   ```
3. Installe les outils nécessaires avec Mason :

   ```vim
   :MasonInstall gopls goimports golangci-lint
   ```
4. Ouvre un fichier Go, modifie du code → `:w` déclenche **goimports + formatage auto** 🎉

---

👉 Veux-tu que je t’ajoute aussi une **config pour lancer `:GoTest` directement dans Neovim avec une touche rapide** (sans passer par vim-go) ?

Top 👍 ajoutons une petite intégration pour **exécuter les tests Go directement dans Neovim**, sans installer `vim-go`.
On va utiliser un simple **autocommande + keymap** qui lance `go test` dans une fenêtre flottante (via `:terminal`).

---

## 📜 Ajout à ton `init.lua`

Colle ce bloc **après la section 8 (format auto)** et avant les options pratiques :

```lua
-- === 9. Lancer les tests Go ===
function RunGoTest()
  local cmd = "go test ./..."
  -- Ouvre un terminal flottant
  vim.cmd("botright split | resize 15 | terminal " .. cmd)
end

-- Keymap pratique pour lancer les tests
vim.keymap.set("n", "<leader>t", RunGoTest, { noremap = true, silent = true })
```

---

## 🎮 Utilisation

* Ouvre un projet Go.
* Tape `<leader>t` (par défaut `\t` dans Vim/Neovim, sauf si tu as redéfini `mapleader`).
* Ça lance `go test ./...` dans une **fenêtre terminal en bas**.

Tu peux adapter la commande, par exemple :

* `go test ./mypkg` → seulement un package.
* `go test -v ./...` → mode verbeux.
* `go test -run TestNom ./...` → lancer un test spécifique.

---

👉 Veux-tu que je prépare une version **un peu plus avancée**, avec plusieurs raccourcis :

* `<leader>ta` → tous les tests (`go test ./...`)
* `<leader>tf` → tests du fichier courant
* `<leader>tn` → test sous le curseur (fonction `Test...`) ?
