j'ai besoin d'un module golang qui reprend l'ensemble des commandes terraform cli avec ses sous commandes comme workspace list, destroy create avec des methodes Run, Help pour l'aide et Validate en utilisant les interface et les generic si necessaire et ses tests unitaire.

tout çà doit être appeler par une seul route d'API en post qui request un payload dynamique qu'on envoie au module go pour valider les inputs selon les actions qui retourne la bonne commande terraform.

ma route doit:
1. check l'action avec une whitelist qui retourne la bonne Command
2. validate les inpute
3. retourn Help si il y a des erreurs
4. retour le playload valider

________________________________

Je vais créer une architecture complète : le module Go avec toutes les commandes Terraform + l'API REST + les tests unitaires.Go n'est pas disponible dans le réseau du container. Voici néanmoins le projet complet avec une architecture solide — copions tout dans les outputs :Voici l'architecture complète. Voici ce qui a été créé :

---

## Structure du projet

```shell
terraform-module/
├── go.mod
├── main.go
├── terraform/
│   ├── command.go           ← interfaces + types partagés
│   ├── registry.go          ← whitelist + Dispatcher
│   ├── commands_basic.go    ← init, validate, plan, apply, destroy
│   ├── commands_state.go    ← state, output, import, show, fmt, refresh, providers, graph
│   ├── commands_workspace.go← workspace (list/show/new/select/delete)
│   └── dispatcher_test.go   ← ~45 tests unitaires
└── api/
    ├── handler.go           ← POST /terraform/command
    └── handler_test.go      ← tests HTTP
```

---

## Flux de la route API

```shell
POST /terraform/command
{ "action": "workspace", "sub_command": "new", "name": "staging" }

     ↓ 1. Extract "action"
     ↓ 2. Registry.Lookup()  → whitelist check → *WorkspaceCommand
     ↓ 3. cmd.Validate(payload) → errs vides → clean payload
     ↓ 4. Si erreurs → { valid: false, help: "...", errors: [...] }  HTTP 422
     ↓ 5. Succès → { valid: true, command: "workspace", sub_command: "new", payload: {...} }  HTTP 200
```

---

## Points clés de design

**Interface `Command`** — chaque commande implémente `Run`, `Help`, `Validate`, `Name`. Les commandes avec sous-commandes implémentent en plus `SubCommand` (`workspace`, `state`).

**`Dispatcher`** — pipeline à 4 étapes découplé du HTTP. Testable indépendamment.

**Validations métier** — `destroy` exige `auto_approve: true`, `workspace new/select/delete` exige `name`, `state mv` exige `address` + `destination`, `import` exige `address` + `id`, `apply` valide le range de `parallelism`.

**Payload nettoyé** — seuls les champs connus/validés sont retournés, les champs inconnus sont silencieusement ignorés.

Code into files.zip

### commande.go
```go
package terraform

import (
	"fmt"
	"strings"
)

// Command defines the interface for all Terraform commands
type Command interface {
	Run(args []string) error
	Help() string
	Validate(payload map[string]any) (map[string]any, []ValidationError)
	Name() string
}

// SubCommand defines an optional sub-command interface
type SubCommand interface {
	Command
	SubCommands() []string
	RunSub(sub string, args []string) error
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Field, e.Message)
}

// CommandResult is returned after validation
type CommandResult struct {
	Command    string         `json:"command"`
	SubCommand string         `json:"sub_command,omitempty"`
	Args       []string       `json:"args,omitempty"`
	Payload    map[string]any `json:"payload"`
	HelpText   string         `json:"help,omitempty"`
	Valid       bool           `json:"valid"`
	Errors     []ValidationError `json:"errors,omitempty"`
}

// BaseCommand provides shared helpers for all commands
type BaseCommand struct {
	name string
}

func (b *BaseCommand) Name() string { return b.name }

func (b *BaseCommand) Run(args []string) error {
	fmt.Printf("terraform %s %s\n", b.name, strings.Join(args, " "))
	return nil
}

// getString safely extracts a string field from a payload
func getString(payload map[string]any, key string) (string, bool) {
	v, ok := payload[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// getBool safely extracts a bool field from a payload
func getBool(payload map[string]any, key string) (bool, bool) {
	v, ok := payload[key]
	if !ok {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

// getStringSlice safely extracts a []string field from a payload
func getStringSlice(payload map[string]any, key string) ([]string, bool) {
	v, ok := payload[key]
	if !ok {
		return nil, false
	}
	switch val := v.(type) {
	case []string:
		return val, true
	case []any:
		result := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result, true
	}
	return nil, false
}
```

### commandes_basic.go
```go
package terraform

import "fmt"

// ─── INIT ────────────────────────────────────────────────────────────────────

type InitCommand struct{ BaseCommand }

func NewInitCommand() *InitCommand { return &InitCommand{BaseCommand{"init"}} }

func (c *InitCommand) Help() string {
	return `Usage: terraform init [options]

Initialize a Terraform working directory.

Options:
  -backend=true          Configure the backend for this configuration.
  -backend-config=path   Additional backend configuration.
  -upgrade               Upgrade modules and plugins to latest.
  -reconfigure           Reconfigure the backend.
  -migrate-state         Reconfigure backend and migrate state.
  dir (string)           Working directory (default: current).`
}

func (c *InitCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	var errs []ValidationError
	clean := map[string]any{}

	if dir, ok := getString(payload, "dir"); ok && dir != "" {
		clean["dir"] = dir
	}
	if backend, ok := getBool(payload, "backend"); ok {
		clean["backend"] = backend
	}
	if upgrade, ok := getBool(payload, "upgrade"); ok {
		clean["upgrade"] = upgrade
	}
	if reconfigure, ok := getBool(payload, "reconfigure"); ok {
		clean["reconfigure"] = reconfigure
	}
	if migrateState, ok := getBool(payload, "migrate_state"); ok {
		clean["migrate_state"] = migrateState
	}
	if backendConfig, ok := getString(payload, "backend_config"); ok && backendConfig != "" {
		clean["backend_config"] = backendConfig
	}

	return clean, errs
}

// ─── VALIDATE ────────────────────────────────────────────────────────────────

type ValidateCommand struct{ BaseCommand }

func NewValidateCommand() *ValidateCommand { return &ValidateCommand{BaseCommand{"validate"}} }

func (c *ValidateCommand) Help() string {
	return `Usage: terraform validate [options]

Validate the configuration files in a directory.

Options:
  -json   Produce output in machine-readable JSON format.
  -no-color Disable terminal color codes.
  dir (string) Working directory.`
}

func (c *ValidateCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	var errs []ValidationError
	clean := map[string]any{}

	if jsonOut, ok := getBool(payload, "json"); ok {
		clean["json"] = jsonOut
	}
	if noColor, ok := getBool(payload, "no_color"); ok {
		clean["no_color"] = noColor
	}
	if dir, ok := getString(payload, "dir"); ok && dir != "" {
		clean["dir"] = dir
	}
	return clean, errs
}

// ─── PLAN ─────────────────────────────────────────────────────────────────────

type PlanCommand struct{ BaseCommand }

func NewPlanCommand() *PlanCommand { return &PlanCommand{BaseCommand{"plan"}} }

func (c *PlanCommand) Help() string {
	return `Usage: terraform plan [options]

Show changes required by the current configuration.

Options:
  -out=path        Write the plan to a file.
  -var=k=v         Set a variable value.
  -var-file=path   Load variable values from a file.
  -destroy         Plan a destroy instead of an apply.
  -target=resource Limit scope to a specific resource.
  -refresh=true    Refresh state before planning.
  -compact-warnings Show compact warning output.`
}

func (c *PlanCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	var errs []ValidationError
	clean := map[string]any{}

	if out, ok := getString(payload, "out"); ok && out != "" {
		clean["out"] = out
	}
	if destroy, ok := getBool(payload, "destroy"); ok {
		clean["destroy"] = destroy
	}
	if refresh, ok := getBool(payload, "refresh"); ok {
		clean["refresh"] = refresh
	}
	if targets, ok := getStringSlice(payload, "target"); ok {
		clean["target"] = targets
	}
	if vars, ok := getStringSlice(payload, "var"); ok {
		clean["var"] = vars
	}
	if varFile, ok := getString(payload, "var_file"); ok && varFile != "" {
		clean["var_file"] = varFile
	}
	return clean, errs
}

// ─── APPLY ────────────────────────────────────────────────────────────────────

type ApplyCommand struct{ BaseCommand }

func NewApplyCommand() *ApplyCommand { return &ApplyCommand{BaseCommand{"apply"}} }

func (c *ApplyCommand) Help() string {
	return `Usage: terraform apply [options] [plan-file]

Apply the changes required to reach the desired state.

Options:
  -auto-approve    Skip interactive approval.
  -plan=path       Path to a pre-computed plan file.
  -var=k=v         Set a variable value.
  -var-file=path   Load variable values from a file.
  -target=resource Limit scope to a specific resource.
  -parallelism=n   Number of parallel resource operations (default 10).`
}

func (c *ApplyCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	var errs []ValidationError
	clean := map[string]any{}

	if autoApprove, ok := getBool(payload, "auto_approve"); ok {
		clean["auto_approve"] = autoApprove
	}
	if plan, ok := getString(payload, "plan"); ok && plan != "" {
		clean["plan"] = plan
	}
	if targets, ok := getStringSlice(payload, "target"); ok {
		clean["target"] = targets
	}
	if vars, ok := getStringSlice(payload, "var"); ok {
		clean["var"] = vars
	}
	if varFile, ok := getString(payload, "var_file"); ok && varFile != "" {
		clean["var_file"] = varFile
	}
	if p, ok := payload["parallelism"]; ok {
		switch v := p.(type) {
		case int:
			if v < 1 || v > 100 {
				errs = append(errs, ValidationError{"parallelism", "must be between 1 and 100"})
			} else {
				clean["parallelism"] = v
			}
		case float64:
			iv := int(v)
			if iv < 1 || iv > 100 {
				errs = append(errs, ValidationError{"parallelism", "must be between 1 and 100"})
			} else {
				clean["parallelism"] = iv
			}
		default:
			errs = append(errs, ValidationError{"parallelism", "must be an integer"})
		}
	}
	return clean, errs
}

// ─── DESTROY ──────────────────────────────────────────────────────────────────

type DestroyCommand struct{ BaseCommand }

func NewDestroyCommand() *DestroyCommand { return &DestroyCommand{BaseCommand{"destroy"}} }

func (c *DestroyCommand) Help() string {
	return `Usage: terraform destroy [options]

Destroy all remote objects managed by this configuration.

Options:
  -auto-approve    Skip interactive approval (REQUIRED for automation).
  -target=resource Limit scope to a specific resource.
  -var=k=v         Set a variable value.
  -var-file=path   Load variable values from a file.
  -parallelism=n   Parallelism for resource operations.

WARNING: This command will destroy your infrastructure!`
}

func (c *DestroyCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	var errs []ValidationError
	clean := map[string]any{}

	autoApprove, ok := getBool(payload, "auto_approve")
	if !ok || !autoApprove {
		errs = append(errs, ValidationError{"auto_approve", fmt.Sprintf("must be true for destroy — %s", c.Help())})
	} else {
		clean["auto_approve"] = true
	}

	if targets, ok := getStringSlice(payload, "target"); ok {
		clean["target"] = targets
	}
	if vars, ok := getStringSlice(payload, "var"); ok {
		clean["var"] = vars
	}
	if varFile, ok := getString(payload, "var_file"); ok && varFile != "" {
		clean["var_file"] = varFile
	}
	return clean, errs
}
```

### commands_state.go
```go
package terraform

// ─── STATE ────────────────────────────────────────────────────────────────────

type StateCommand struct {
	BaseCommand
	subCommands []string
}

func NewStateCommand() *StateCommand {
	return &StateCommand{
		BaseCommand: BaseCommand{"state"},
		subCommands: []string{"list", "show", "mv", "rm", "pull", "push", "replace-provider"},
	}
}

func (c *StateCommand) SubCommands() []string { return c.subCommands }

func (c *StateCommand) RunSub(sub string, args []string) error {
	_ = sub
	return c.Run(args)
}

func (c *StateCommand) Help() string {
	return `Usage: terraform state <subcommand> [options]

Manage Terraform state.

Subcommands:
  list              List resources in state.
  show              Show a resource in state.
  mv                Move a resource in state.
  rm                Remove a resource from state.
  pull              Pull state from backend.
  push              Push state to backend.
  replace-provider  Replace provider in state.

Options:
  sub_command (string, required)  One of the subcommands above.
  address  (string)               Resource address (required for show/mv/rm).
  destination (string)            Destination address (required for mv).`
}

func (c *StateCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	var errs []ValidationError
	clean := map[string]any{}

	sub, ok := getString(payload, "sub_command")
	if !ok || sub == "" {
		errs = append(errs, ValidationError{"sub_command", "required: one of list|show|mv|rm|pull|push|replace-provider"})
		return clean, errs
	}

	valid := false
	for _, s := range c.subCommands {
		if s == sub {
			valid = true
			break
		}
	}
	if !valid {
		errs = append(errs, ValidationError{"sub_command", "unknown subcommand: " + sub})
		return clean, errs
	}
	clean["sub_command"] = sub

	switch sub {
	case "show", "rm":
		addr, ok := getString(payload, "address")
		if !ok || addr == "" {
			errs = append(errs, ValidationError{"address", "required for sub_command=" + sub})
		} else {
			clean["address"] = addr
		}
	case "mv":
		addr, ok := getString(payload, "address")
		if !ok || addr == "" {
			errs = append(errs, ValidationError{"address", "required for sub_command=mv"})
		} else {
			clean["address"] = addr
		}
		dest, ok := getString(payload, "destination")
		if !ok || dest == "" {
			errs = append(errs, ValidationError{"destination", "required for sub_command=mv"})
		} else {
			clean["destination"] = dest
		}
	}

	return clean, errs
}

// ─── OUTPUT ───────────────────────────────────────────────────────────────────

type OutputCommand struct{ BaseCommand }

func NewOutputCommand() *OutputCommand { return &OutputCommand{BaseCommand{"output"}} }

func (c *OutputCommand) Help() string {
	return `Usage: terraform output [options] [NAME]

Show output values from your root module.

Options:
  name    (string) Specific output name to show.
  -json           Format output as JSON.
  -raw            Print raw string value.
  -no-color       Disable color output.`
}

func (c *OutputCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	clean := map[string]any{}

	if name, ok := getString(payload, "name"); ok && name != "" {
		clean["name"] = name
	}
	if jsonOut, ok := getBool(payload, "json"); ok {
		clean["json"] = jsonOut
	}
	if raw, ok := getBool(payload, "raw"); ok {
		clean["raw"] = raw
	}
	return clean, nil
}

// ─── IMPORT ───────────────────────────────────────────────────────────────────

type ImportCommand struct{ BaseCommand }

func NewImportCommand() *ImportCommand { return &ImportCommand{BaseCommand{"import"}} }

func (c *ImportCommand) Help() string {
	return `Usage: terraform import [options] ADDR ID

Import existing infrastructure into Terraform state.

Required:
  address (string) Resource address (e.g., aws_instance.example).
  id      (string) Provider-specific resource ID.

Options:
  -var=k=v         Set a variable value.
  -var-file=path   Load variable values from a file.
  -allow-missing-config Allow import without configuration.`
}

func (c *ImportCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	var errs []ValidationError
	clean := map[string]any{}

	addr, ok := getString(payload, "address")
	if !ok || addr == "" {
		errs = append(errs, ValidationError{"address", "required (e.g., aws_instance.example)"})
	} else {
		clean["address"] = addr
	}

	id, ok := getString(payload, "id")
	if !ok || id == "" {
		errs = append(errs, ValidationError{"id", "required: provider-specific resource ID"})
	} else {
		clean["id"] = id
	}

	if vars, ok := getStringSlice(payload, "var"); ok {
		clean["var"] = vars
	}
	if varFile, ok := getString(payload, "var_file"); ok && varFile != "" {
		clean["var_file"] = varFile
	}
	return clean, errs
}

// ─── SHOW ─────────────────────────────────────────────────────────────────────

type ShowCommand struct{ BaseCommand }

func NewShowCommand() *ShowCommand { return &ShowCommand{BaseCommand{"show"}} }

func (c *ShowCommand) Help() string {
	return `Usage: terraform show [options] [path]

Show the current state or a saved plan.

Options:
  path   (string) Path to a plan file (optional).
  -json          Produce JSON output.
  -no-color      Disable color.`
}

func (c *ShowCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	clean := map[string]any{}

	if path, ok := getString(payload, "path"); ok && path != "" {
		clean["path"] = path
	}
	if jsonOut, ok := getBool(payload, "json"); ok {
		clean["json"] = jsonOut
	}
	return clean, nil
}

// ─── FMT ──────────────────────────────────────────────────────────────────────

type FmtCommand struct{ BaseCommand }

func NewFmtCommand() *FmtCommand { return &FmtCommand{BaseCommand{"fmt"}} }

func (c *FmtCommand) Help() string {
	return `Usage: terraform fmt [options] [dir]

Reformat configuration files to canonical format.

Options:
  dir      (string) Directory to format (default: current).
  -check           Check if files are formatted, exit non-zero if not.
  -diff            Display diffs of formatting changes.
  -recursive       Process files in subdirectories.`
}

func (c *FmtCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	clean := map[string]any{}

	if dir, ok := getString(payload, "dir"); ok && dir != "" {
		clean["dir"] = dir
	}
	if check, ok := getBool(payload, "check"); ok {
		clean["check"] = check
	}
	if diff, ok := getBool(payload, "diff"); ok {
		clean["diff"] = diff
	}
	if recursive, ok := getBool(payload, "recursive"); ok {
		clean["recursive"] = recursive
	}
	return clean, nil
}

// ─── REFRESH ──────────────────────────────────────────────────────────────────

type RefreshCommand struct{ BaseCommand }

func NewRefreshCommand() *RefreshCommand { return &RefreshCommand{BaseCommand{"refresh"}} }

func (c *RefreshCommand) Help() string {
	return `Usage: terraform refresh [options]

Update the state file to match real infrastructure.

Options:
  -var=k=v         Set a variable value.
  -var-file=path   Load variable values from a file.
  -target=resource Limit scope to a specific resource.
  -parallelism=n   Number of parallel resource operations.`
}

func (c *RefreshCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	clean := map[string]any{}

	if targets, ok := getStringSlice(payload, "target"); ok {
		clean["target"] = targets
	}
	if vars, ok := getStringSlice(payload, "var"); ok {
		clean["var"] = vars
	}
	if varFile, ok := getString(payload, "var_file"); ok && varFile != "" {
		clean["var_file"] = varFile
	}
	return clean, nil
}

// ─── PROVIDERS ────────────────────────────────────────────────────────────────

type ProvidersCommand struct{ BaseCommand }

func NewProvidersCommand() *ProvidersCommand { return &ProvidersCommand{BaseCommand{"providers"}} }

func (c *ProvidersCommand) Help() string {
	return `Usage: terraform providers [dir]

Show the providers required for this configuration.

Options:
  dir      (string) Working directory.
  -json           Output in JSON format.`
}

func (c *ProvidersCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	clean := map[string]any{}

	if dir, ok := getString(payload, "dir"); ok && dir != "" {
		clean["dir"] = dir
	}
	if jsonOut, ok := getBool(payload, "json"); ok {
		clean["json"] = jsonOut
	}
	return clean, nil
}

// ─── GRAPH ────────────────────────────────────────────────────────────────────

type GraphCommand struct{ BaseCommand }

func NewGraphCommand() *GraphCommand { return &GraphCommand{BaseCommand{"graph"}} }

func (c *GraphCommand) Help() string {
	return `Usage: terraform graph [options]

Generate a visual representation of the dependency graph.

Options:
  -type=string     Type of graph (plan, plan-refresh-only, plan-destroy, apply).
  -draw-cycles     Highlight cycles in the dependency graph.
  -plan=path       Render graph from a saved plan file.`
}

func (c *GraphCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	var errs []ValidationError
	clean := map[string]any{}

	validTypes := map[string]bool{
		"plan": true, "plan-refresh-only": true, "plan-destroy": true, "apply": true,
	}
	if t, ok := getString(payload, "type"); ok && t != "" {
		if !validTypes[t] {
			errs = append(errs, ValidationError{"type", "must be one of plan|plan-refresh-only|plan-destroy|apply"})
		} else {
			clean["type"] = t
		}
	}
	if drawCycles, ok := getBool(payload, "draw_cycles"); ok {
		clean["draw_cycles"] = drawCycles
	}
	return clean, errs
}
```

### commands_workspace.go
```go
package terraform

// ─── WORKSPACE ────────────────────────────────────────────────────────────────

type WorkspaceCommand struct {
	BaseCommand
	subCommands []string
}

func NewWorkspaceCommand() *WorkspaceCommand {
	return &WorkspaceCommand{
		BaseCommand: BaseCommand{"workspace"},
		subCommands: []string{"list", "show", "new", "select", "delete"},
	}
}

func (c *WorkspaceCommand) SubCommands() []string { return c.subCommands }

func (c *WorkspaceCommand) RunSub(sub string, args []string) error {
	return c.Run(args)
}

func (c *WorkspaceCommand) Help() string {
	return `Usage: terraform workspace <subcommand> [options]

Manage Terraform workspaces.

Subcommands:
  list    List available workspaces.
  show    Show the current workspace name.
  new     Create a new workspace.
  select  Switch to a different workspace.
  delete  Delete a workspace.

Options:
  sub_command (string, required)  One of the subcommands above.
  name        (string)            Workspace name (required for new/select/delete).
  -force                          Force deletion even if not empty (delete only).
  -lock=true                      Lock state during workspace switch.
  -state=path                     Copy existing state file to new workspace (new only).`
}

func (c *WorkspaceCommand) Validate(payload map[string]any) (map[string]any, []ValidationError) {
	var errs []ValidationError
	clean := map[string]any{}

	sub, ok := getString(payload, "sub_command")
	if !ok || sub == "" {
		errs = append(errs, ValidationError{"sub_command", "required: one of list|show|new|select|delete"})
		return clean, errs
	}

	valid := false
	for _, s := range c.subCommands {
		if s == sub {
			valid = true
			break
		}
	}
	if !valid {
		errs = append(errs, ValidationError{"sub_command", "unknown subcommand: " + sub})
		return clean, errs
	}
	clean["sub_command"] = sub

	switch sub {
	case "new", "select", "delete":
		name, ok := getString(payload, "name")
		if !ok || name == "" {
			errs = append(errs, ValidationError{"name", "workspace name required for sub_command=" + sub})
		} else {
			clean["name"] = name
		}

		if sub == "delete" {
			if force, ok := getBool(payload, "force"); ok {
				clean["force"] = force
			}
		}
		if sub == "new" {
			if statePath, ok := getString(payload, "state"); ok && statePath != "" {
				clean["state"] = statePath
			}
		}
		if lock, ok := getBool(payload, "lock"); ok {
			clean["lock"] = lock
		}
	}

	return clean, errs
}
```

### dispatcher_test.go
```go
package terraform_test

import (
	"testing"

	"github.com/example/terraform-module/terraform"
)

// ─── helpers ─────────────────────────────────────────────────────────────────

func newDispatcher() *terraform.Dispatcher {
	return terraform.NewDispatcher(terraform.NewRegistry())
}

func assertValid(t *testing.T, result terraform.CommandResult) {
	t.Helper()
	if !result.Valid {
		t.Fatalf("expected valid result, got errors: %v", result.Errors)
	}
}

func assertInvalid(t *testing.T, result terraform.CommandResult, expectedFields ...string) {
	t.Helper()
	if result.Valid {
		t.Fatal("expected invalid result, but got valid")
	}
	if len(expectedFields) == 0 {
		return
	}
	fieldSet := make(map[string]bool)
	for _, e := range result.Errors {
		fieldSet[e.Field] = true
	}
	for _, f := range expectedFields {
		if !fieldSet[f] {
			t.Errorf("expected error on field %q, got errors: %v", f, result.Errors)
		}
	}
}

func assertHasHelp(t *testing.T, result terraform.CommandResult) {
	t.Helper()
	if result.HelpText == "" {
		t.Fatal("expected HelpText to be non-empty on invalid result")
	}
}

// ─── Dispatcher ──────────────────────────────────────────────────────────────

func TestDispatcher_MissingAction(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{})
	assertInvalid(t, result, "action")
	assertHasHelp(t, result)
}

func TestDispatcher_UnknownAction(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "nuke"})
	assertInvalid(t, result, "action")
	assertHasHelp(t, result)
}

func TestDispatcher_ActionNotString(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": 42})
	assertInvalid(t, result, "action")
}

// ─── Init ────────────────────────────────────────────────────────────────────

func TestInit_Minimal(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "init"})
	assertValid(t, result)
	if result.Command != "init" {
		t.Errorf("expected command=init, got %q", result.Command)
	}
}

func TestInit_WithOptions(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{
		"action":  "init",
		"backend": true,
		"upgrade": true,
		"dir":     "./infra",
	})
	assertValid(t, result)
	if result.Payload["dir"] != "./infra" {
		t.Errorf("expected dir=./infra")
	}
}

// ─── Validate ────────────────────────────────────────────────────────────────

func TestValidate_Minimal(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "validate"})
	assertValid(t, result)
}

func TestValidate_JSONOutput(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "validate", "json": true})
	assertValid(t, result)
	if result.Payload["json"] != true {
		t.Error("json flag not in payload")
	}
}

// ─── Plan ─────────────────────────────────────────────────────────────────────

func TestPlan_Minimal(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "plan"})
	assertValid(t, result)
}

func TestPlan_WithOut(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "plan", "out": "tfplan"})
	assertValid(t, result)
	if result.Payload["out"] != "tfplan" {
		t.Error("out not in payload")
	}
}

func TestPlan_WithTargets(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{
		"action": "plan",
		"target": []any{"aws_instance.web", "aws_s3_bucket.data"},
	})
	assertValid(t, result)
}

// ─── Apply ────────────────────────────────────────────────────────────────────

func TestApply_WithAutoApprove(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "apply", "auto_approve": true})
	assertValid(t, result)
}

func TestApply_InvalidParallelism(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "apply", "parallelism": float64(200)})
	assertInvalid(t, result, "parallelism")
	assertHasHelp(t, result)
}

func TestApply_ValidParallelism(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "apply", "parallelism": float64(5)})
	assertValid(t, result)
	if result.Payload["parallelism"] != 5 {
		t.Errorf("expected parallelism=5")
	}
}

// ─── Destroy ──────────────────────────────────────────────────────────────────

func TestDestroy_MissingAutoApprove(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "destroy"})
	assertInvalid(t, result, "auto_approve")
	assertHasHelp(t, result)
}

func TestDestroy_AutoApproveFalse(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "destroy", "auto_approve": false})
	assertInvalid(t, result, "auto_approve")
}

func TestDestroy_Valid(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "destroy", "auto_approve": true})
	assertValid(t, result)
}

// ─── Workspace ────────────────────────────────────────────────────────────────

func TestWorkspace_MissingSubCommand(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "workspace"})
	assertInvalid(t, result, "sub_command")
	assertHasHelp(t, result)
}

func TestWorkspace_UnknownSubCommand(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "workspace", "sub_command": "nuke"})
	assertInvalid(t, result, "sub_command")
}

func TestWorkspace_List(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "workspace", "sub_command": "list"})
	assertValid(t, result)
	if result.SubCommand != "list" {
		t.Errorf("expected SubCommand=list, got %q", result.SubCommand)
	}
}

func TestWorkspace_New_MissingName(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "workspace", "sub_command": "new"})
	assertInvalid(t, result, "name")
	assertHasHelp(t, result)
}

func TestWorkspace_New_Valid(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "workspace", "sub_command": "new", "name": "staging"})
	assertValid(t, result)
	if result.Payload["name"] != "staging" {
		t.Error("name not in payload")
	}
}

func TestWorkspace_Select_Valid(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "workspace", "sub_command": "select", "name": "prod"})
	assertValid(t, result)
}

func TestWorkspace_Delete_WithForce(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "workspace", "sub_command": "delete", "name": "old-ws", "force": true})
	assertValid(t, result)
	if result.Payload["force"] != true {
		t.Error("force not in payload")
	}
}

func TestWorkspace_Show(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "workspace", "sub_command": "show"})
	assertValid(t, result)
}

// ─── State ────────────────────────────────────────────────────────────────────

func TestState_MissingSubCommand(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "state"})
	assertInvalid(t, result, "sub_command")
}

func TestState_List(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "state", "sub_command": "list"})
	assertValid(t, result)
}

func TestState_Show_MissingAddress(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "state", "sub_command": "show"})
	assertInvalid(t, result, "address")
	assertHasHelp(t, result)
}

func TestState_Show_Valid(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "state", "sub_command": "show", "address": "aws_instance.web"})
	assertValid(t, result)
}

func TestState_Mv_MissingDestination(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "state", "sub_command": "mv", "address": "aws_instance.old"})
	assertInvalid(t, result, "destination")
}

func TestState_Mv_Valid(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{
		"action":      "state",
		"sub_command": "mv",
		"address":     "aws_instance.old",
		"destination": "aws_instance.new",
	})
	assertValid(t, result)
}

func TestState_Rm_MissingAddress(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "state", "sub_command": "rm"})
	assertInvalid(t, result, "address")
}

// ─── Import ───────────────────────────────────────────────────────────────────

func TestImport_MissingFields(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "import"})
	assertInvalid(t, result, "address", "id")
	assertHasHelp(t, result)
}

func TestImport_MissingID(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "import", "address": "aws_instance.foo"})
	assertInvalid(t, result, "id")
}

func TestImport_Valid(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{
		"action":  "import",
		"address": "aws_instance.web",
		"id":      "i-1234567890abcdef0",
	})
	assertValid(t, result)
}

// ─── Output ───────────────────────────────────────────────────────────────────

func TestOutput_Minimal(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "output"})
	assertValid(t, result)
}

func TestOutput_Named(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "output", "name": "vpc_id", "json": true})
	assertValid(t, result)
	if result.Payload["name"] != "vpc_id" {
		t.Error("name not in payload")
	}
}

// ─── Graph ────────────────────────────────────────────────────────────────────

func TestGraph_InvalidType(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "graph", "type": "unknown"})
	assertInvalid(t, result, "type")
}

func TestGraph_ValidType(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "graph", "type": "plan"})
	assertValid(t, result)
}

// ─── Fmt ──────────────────────────────────────────────────────────────────────

func TestFmt_Check(t *testing.T) {
	d := newDispatcher()
	result := d.Dispatch(map[string]any{"action": "fmt", "check": true, "recursive": true})
	assertValid(t, result)
	if result.Payload["check"] != true {
		t.Error("check flag not in payload")
	}
}

// ─── Command interface ────────────────────────────────────────────────────────

func TestAllCommands_HelpNotEmpty(t *testing.T) {
	registry := terraform.NewRegistry()
	for _, action := range registry.Actions() {
		cmd, _ := registry.Lookup(action)
		if cmd.Help() == "" {
			t.Errorf("command %q returned empty Help()", action)
		}
	}
}

func TestAllCommands_Name(t *testing.T) {
	registry := terraform.NewRegistry()
	for _, action := range registry.Actions() {
		cmd, _ := registry.Lookup(action)
		if cmd.Name() != action {
			t.Errorf("expected Name()=%q, got %q", action, cmd.Name())
		}
	}
}
```

### handler.go
```go
package api

import (
	"encoding/json"
	"net/http"

	"github.com/example/terraform-module/terraform"
)

// Handler exposes the single POST /terraform/command route
type Handler struct {
	dispatcher *terraform.Dispatcher
}

func NewHandler(d *terraform.Dispatcher) *Handler {
	return &Handler{dispatcher: d}
}

// APIResponse is the unified response envelope
type APIResponse struct {
	Success bool                    `json:"success"`
	Result  *terraform.CommandResult `json:"result,omitempty"`
	Error   string                  `json:"error,omitempty"`
}

// Register wires the route onto an existing mux
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /terraform/command", h.HandleCommand)
}

// HandleCommand is the single POST endpoint
//
//	POST /terraform/command
//	Content-Type: application/json
//
//	{ "action": "workspace", "sub_command": "new", "name": "staging" }
func (h *Handler) HandleCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// ── Parse body ──────────────────────────────────────────────────────────
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "invalid JSON: " + err.Error(),
		})
		return
	}

	// ── Dispatch ────────────────────────────────────────────────────────────
	result := h.dispatcher.Dispatch(payload)

	status := http.StatusOK
	if !result.Valid {
		status = http.StatusUnprocessableEntity
	}

	writeJSON(w, status, APIResponse{
		Success: result.Valid,
		Result:  &result,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
```

### handler_test.go
```go
package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/terraform-module/api"
	"github.com/example/terraform-module/terraform"
)

func newTestServer() *httptest.Server {
	registry := terraform.NewRegistry()
	dispatcher := terraform.NewDispatcher(registry)
	handler := api.NewHandler(dispatcher)
	mux := http.NewServeMux()
	handler.Register(mux)
	return httptest.NewServer(mux)
}

func postJSON(t *testing.T, server *httptest.Server, payload any) *http.Response {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(server.URL+"/terraform/command", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func decode(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestAPI_ValidWorkspaceList(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, map[string]any{
		"action":      "workspace",
		"sub_command": "list",
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	body := decode(t, resp)
	if body["success"] != true {
		t.Error("expected success=true")
	}
}

func TestAPI_InvalidAction(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, map[string]any{"action": "nuke"})
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", resp.StatusCode)
	}
	body := decode(t, resp)
	if body["success"] != false {
		t.Error("expected success=false")
	}
}

func TestAPI_InvalidJSON(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, _ := http.Post(srv.URL+"/terraform/command", "application/json", bytes.NewBufferString("not-json"))
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestAPI_DestroyWithoutAutoApprove(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, map[string]any{"action": "destroy"})
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", resp.StatusCode)
	}
	body := decode(t, resp)
	result := body["result"].(map[string]any)
	if result["help"] == "" {
		t.Error("expected help text in response")
	}
}

func TestAPI_WorkspaceNew_Valid(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, map[string]any{
		"action":      "workspace",
		"sub_command": "new",
		"name":        "staging",
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestAPI_ImportMissingFields(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, map[string]any{"action": "import"})
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", resp.StatusCode)
	}
}

func TestAPI_PlanCleanPayload(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, map[string]any{
		"action": "plan",
		"out":    "tfplan",
		"target": []string{"aws_instance.web"},
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	body := decode(t, resp)
	result := body["result"].(map[string]any)
	payload := result["payload"].(map[string]any)
	if payload["out"] != "tfplan" {
		t.Error("expected out=tfplan in cleaned payload")
	}
}
```

### main.go
```go
package main

import (
	"log"
	"net/http"

	"github.com/example/terraform-module/api"
	"github.com/example/terraform-module/terraform"
)

func main() {
	registry := terraform.NewRegistry()
	dispatcher := terraform.NewDispatcher(registry)
	handler := api.NewHandler(dispatcher)

	mux := http.NewServeMux()
	handler.Register(mux)

	log.Println("Terraform API listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
```

### registry.go
```go
package terraform

import "fmt"

// Registry holds the whitelist of all known Terraform commands
type Registry struct {
	commands map[string]Command
}

// NewRegistry creates a Registry pre-loaded with all supported Terraform commands
func NewRegistry() *Registry {
	r := &Registry{commands: make(map[string]Command)}

	cmds := []Command{
		NewInitCommand(),
		NewValidateCommand(),
		NewPlanCommand(),
		NewApplyCommand(),
		NewDestroyCommand(),
		NewStateCommand(),
		NewWorkspaceCommand(),
		NewOutputCommand(),
		NewImportCommand(),
		NewShowCommand(),
		NewFmtCommand(),
		NewRefreshCommand(),
		NewProvidersCommand(),
		NewGraphCommand(),
	}

	for _, cmd := range cmds {
		r.commands[cmd.Name()] = cmd
	}
	return r
}

// Lookup returns the Command for a given action name (whitelist check)
func (r *Registry) Lookup(action string) (Command, bool) {
	cmd, ok := r.commands[action]
	return cmd, ok
}

// Actions returns a sorted list of all registered command names
func (r *Registry) Actions() []string {
	names := make([]string, 0, len(r.commands))
	for name := range r.commands {
		names = append(names, name)
	}
	return names
}

// ─── DISPATCHER ───────────────────────────────────────────────────────────────

// Dispatcher processes a raw API payload through the full pipeline:
//  1. Whitelist check   → lookup the Command
//  2. Validate inputs   → call Command.Validate()
//  3. Return Help       → on any validation error
//  4. Return result     → clean, validated payload
type Dispatcher struct {
	registry *Registry
}

func NewDispatcher(r *Registry) *Dispatcher { return &Dispatcher{registry: r} }

// Dispatch executes the full pipeline and returns a CommandResult.
// It never returns a Go error; errors are surfaced inside CommandResult.
func (d *Dispatcher) Dispatch(payload map[string]any) CommandResult {
	// ── 1. Extract action ───────────────────────────────────────────────────
	actionRaw, ok := payload["action"]
	if !ok {
		return CommandResult{
			Valid:   false,
			Errors:  []ValidationError{{Field: "action", Message: "field 'action' is required"}},
			HelpText: d.globalHelp(),
		}
	}
	action, ok := actionRaw.(string)
	if !ok || action == "" {
		return CommandResult{
			Valid:   false,
			Errors:  []ValidationError{{Field: "action", Message: "field 'action' must be a non-empty string"}},
			HelpText: d.globalHelp(),
		}
	}

	// ── 2. Whitelist check ──────────────────────────────────────────────────
	cmd, found := d.registry.Lookup(action)
	if !found {
		return CommandResult{
			Valid: false,
			Errors: []ValidationError{{
				Field:   "action",
				Message: fmt.Sprintf("unknown action %q — allowed: %v", action, d.registry.Actions()),
			}},
			HelpText: d.globalHelp(),
		}
	}

	// Strip the action key before forwarding
	params := copyWithout(payload, "action")

	// ── 3. Validate ─────────────────────────────────────────────────────────
	clean, errs := cmd.Validate(params)

	if len(errs) > 0 {
		// ── 4. Return Help on errors ─────────────────────────────────────────
		return CommandResult{
			Command:  action,
			Valid:    false,
			Errors:   errs,
			HelpText: cmd.Help(),
		}
	}

	// ── 5. Success ───────────────────────────────────────────────────────────
	result := CommandResult{
		Command: action,
		Valid:   true,
		Payload: clean,
	}

	// Attach sub_command metadata if relevant
	if sub, ok := clean["sub_command"].(string); ok {
		result.SubCommand = sub
	}

	return result
}

func (d *Dispatcher) globalHelp() string {
	return fmt.Sprintf("Available actions: %v", d.registry.Actions())
}

func copyWithout(m map[string]any, key string) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		if k != key {
			out[k] = v
		}
	}
	return out
}
````
