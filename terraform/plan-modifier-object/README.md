En Terraform, un `ObjectPlanModifier` avec `UseStateForUnknown` est utilisé pour gérer les valeurs d'attributs complexes (objets) lors de la planification, particulièrement quand certaines valeurs ne peuvent pas être connues avant l'application.

## Contexte et utilité

`UseStateForUnknown` est particulièrement utile dans ces situations :
- Quand un attribut objet contient des valeurs calculées par le provider
- Pour éviter des "diffs" non nécessaires entre le plan et l'état
- Quand certaines propriétés de l'objet dépendent de ressources non encore créées

## Implémentation basique

```go
func (r *MyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "config": schema.SingleNestedAttribute{
                Attributes: map[string]schema.Attribute{
                    "computed_field": schema.StringAttribute{
                        Computed: true,
                    },
                    "user_field": schema.StringAttribute{
                        Optional: true,
                    },
                },
                PlanModifiers: []planmodifier.Object{
                    objectplanmodifier.UseStateForUnknown(),
                },
            },
        },
    }
}
```

## Comportement

Avec `UseStateForUnknown()`, Terraform :
1. **Lors du plan** : Si la nouvelle valeur contient des `unknown` et que l'état actuel a une valeur, il utilise la valeur de l'état
2. **Évite les changements forcés** : Réduit les recomputations inutiles d'attributs dépendants
3. **Préserve la stabilité** : Maintient la cohérence entre les cycles plan/apply

## Exemple pratique

```go
// Cas d'usage typique : configuration avec des valeurs calculées
"database_config": schema.SingleNestedAttribute{
    Attributes: map[string]schema.Attribute{
        "endpoint": schema.StringAttribute{
            Computed: true, // Calculé par le provider
        },
        "port": schema.Int64Attribute{
            Optional: true,
            Computed: true,
        },
        "ssl_config": schema.SingleNestedAttribute{
            Computed: true,
            Attributes: map[string]schema.Attribute{
                "ca_cert": schema.StringAttribute{Computed: true},
                "enabled": schema.BoolAttribute{Computed: true},
            },
        },
    },
    PlanModifiers: []planmodifier.Object{
        objectplanmodifier.UseStateForUnknown(),
    },
},
```

## Points d'attention

- **Ne pas overutiliser** : Seulement pour les objets avec des valeurs vraiment inconnues au moment du plan
- **Combinaison avec RequiresReplace** : Peut être combiné avec d'autres plan modifiers selon les besoins
- **Impact sur les dépendances** : Peut affecter comment Terraform résout les dépendances entre ressources

Cette approche est essentielle pour créer des providers Terraform robustes qui gèrent correctement les valeurs calculées dans des structures complexes.

#######################################################################################################################################

Je vais vous créer un exemple complet d'un provider Terraform custom avec `ObjectPlanModifier` et `UseStateForUnknown`, incluant les tests.Voici un exemple complet d'un provider Terraform custom avec `ObjectPlanModifier` et `UseStateForUnknown`. L'exemple simule une ressource "application" avec une configuration complexe contenant des valeurs calculées.

## Points clés de l'implémentation :

### 1. **Structure de l'objet avec UseStateForUnknown**
```go
"config": schema.SingleNestedAttribute{
    // ...
    PlanModifiers: []planmodifier.Object{
        objectplanmodifier.UseStateForUnknown(),
    },
},
```

### 2. **Valeurs calculées dans l'objet**
- `endpoint` : URL générée basée sur l'environnement
- `port` : Port assigné dynamiquement 
- `ssl` : Configuration SSL complète générée automatiquement

### 3. **Tests spécialisés**
- **Test principal** : Vérifie la création, mise à jour et import
- **Test UseStateForUnknown** : Vérifie que les valeurs calculées restent stables entre les applications avec `PlanOnly: true`

## Avantages de cette approche :

1. **Stabilité des plans** : Les valeurs calculées ne changent pas inutilement
2. **Performance** : Évite les recomputations non nécessaires
3. **Expérience utilisateur** : Plans plus prévisibles

## Pour tester l'exemple :

```bash
# Lancer les tests d'acceptance
make testacc

# Ou directement
TF_ACC=1 go test ./internal/provider -v -timeout 120m
```

Le `UseStateForUnknown` est particulièrement utile ici car les valeurs comme `endpoint`, `port` et `ssl` sont calculées côté provider mais doivent rester stables tant que l'utilisateur ne change pas la configuration de base (comme `environment`).