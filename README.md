## question to claud ia

je voudrais un module golang qui gère les commandes terraform cli avec une validation des inputes.
le commande seront init, plan apply et destroy
ainsi que workspace et ses sous commande

la fonction commandes doit être une interface avec 3 functions, Run, Help (pour l'aide) et Validate (pour valider les var - argument) et la commande et sous commande.

J'aimerai avoir un fonctionnement comme lorsqu'on laisse le binaire terraform.

tout ceci sera appelé à partir d'une route d'api en post


1. si pas d'action, retourne l'aide avec une erreur.

2. si action sans argument retourne l'aide avec une erreur

3. si action non exist, retour l'aide avec une erreur

4. si action valide avec de mauvais argument retourn l'aide avec les erreurs

5. si action valide avec argument valide, retourn le payload et genere un fichier tf

6. si actions valide avec une sub command (comme pour worspace) sans argument retourn l'aide avec les erreurs

7. si actions valide avec une sub command (comme pour worspace) avec de mauvais arguments retourn l'aide avec les erreurs

8. si actions valide avec une sub command (comme pour worspace) avec de bon arguments retourn l'aide avec les erreurs.

9. enfin, j'ai besoin d'ajouter un whiteliste de filtre pour limité le nombre d'appel à certaine commande par profile dans le context de go

______________________

Now let me write the four core commands:Now let's verify everything compiles cleanly:Go isn't available in the sandbox. Let me do a careful manual review to verify correctness:I noticed a potential issue — the `"POST /terraform"` mux pattern requires Go 1.22+. Our module declares `go 1.21`. Let me fix this:Now let me also verify the one unused function and clean it up, and do a final pass on the workspace handler's `helpFor` which is an unexported method used only internally:Everything is consistent. Let me also remove the unused `failResult` variadic function to keep things clean:Voici le module complet — 12 fichiers Go, zéro dépendance externe. Voici une synthèse de l'architecture :

---

**Structure du module**

```shell
terraform-commander/
├── go.mod
├── main.go
├── terraform/
│   ├── command.go      ← Interface Command + types partagés
│   ├── registry.go     ← Dispatch + scénarios 1-3
│   ├── generator.go    ← Génération .tf / .tfvars
│   ├── init.go
│   ├── plan.go
│   ├── apply.go
│   ├── destroy.go
│   └── workspace.go    ← new | select | list | delete | show
├── whitelist/
│   └── filter.go       ← Whitelist + compteurs par profil via context
└── api/
    └── handler.go      ← POST /terraform, HTTP 400/422/429/200
```

---

**Interface `Command`**

```go
type Command interface {
    Run(ctx context.Context, req *Request) (*CommandResult, error)
    Help() string
    Validate(req *Request) []ValidationError
}
```

---

**Les 9 scénarios** sont gérés en deux niveaux — `registry.Dispatch` pour les cas globaux (1, 3), puis chaque `Command.Run` → `Validate` pour les cas fins (2, 4-8). Le handler HTTP retourne `400` pour les actions manquantes/inconnues, `422` pour les mauvais arguments, et `429` pour les violations whitelist.

---

**Whitelist / rate-limit par profil (scénario 9)**

Le filtre est injecté dans le `context.Context` via `whitelist.WithFilter(ctx, filter)`. Il gère à la fois une liste blanche de profils autorisés *et* des limites d'appels par commande/sous-commande :

```go
filter := whitelist.New(
    []whitelist.Rule{
        {Command: "apply",     MaxCalls: 5},
        {Command: "destroy",   MaxCalls: 2},
        {Command: "workspace", Subcommand: "new", MaxCalls: 10},
    },
    "alice", "bob", "ci-bot",   // profils autorisés
)
```

---

**Exemple de requête → réponse**

```json
POST /terraform
{ "action": "workspace", "subcommand": "new", "args": {"name": "prod"}, "profile": "alice" }

// 200 OK
{
  "success": true,
  "payload": { "action": "workspace", "subcommand": "new", "args": {...}, "profile": "alice" },
  "generated_file": { "name": "workspace.tf", "content": "locals { workspace_name = \"prod\" }\n" }
}
```
___________
## Add tests unites

All helpers are defined exactly once — `hasError` and `hasErrorContaining` in `init_test.go`, `containsStr` in `plan_test.go`. Since they're all in `package terraform` the other test files can use them freely. Let me do a final check to confirm every test file references only functions that exist:**90 tests unitaires** répartis en 6 fichiers, tous dans `package terraform` (même package que le code testé, donc accès direct aux types internes).

---

**Couverture par fichier**

| Fichier | Tests | Ce qui est couvert |
|---|---|---|
| `init_test.go` | 10 | args vides, nil, flags inconnus, mélange, tous flags valides, génération `backend.tf` |
| `plan_test.go` | 13 | `-var` non-map, `-var-file` vide/non-string, flags inconnus, erreurs multiples, contenu du `.tfvars` |
| `apply_test.go` | 12 | même structure que plan + `-auto-approve`, champ `Help` vide au succès |
| `destroy_test.go` | 11 | idem + vérification du contenu du fichier généré avec les vars |
| `workspace_test.go` | 28 | tous les sous-commandes (`new/select/list/delete/show`), nom avec espace, `-force`, `-state`, args superflus sur `list`/`show` |
| `registry_test.go` | 16 | scénarios 1→8, action en majuscules (case-insensitive), `Register` custom, `GlobalHelp` |

**Convention de nommage :** `Test<Commande><Méthode>_<Condition>_<Résultat>` — lisible comme une phrase (`TestPlanValidate_VarNotAMap_ReturnsError`).

Pour exécuter les tests une fois Go disponible :
```bash
cd terraform-commander
go test ./terraform/... -v
```
