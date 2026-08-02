package builder

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type ActionType string

const (
	ActionTypeAll   EnvType = "ALL"
	ActionTypeAllow EnvType = "DENY"
)

type RuleType string

const (
	RuleTypeAll   RuleType = "ALL"
	RuleTypeAllow RuleType = "ALLOW"
)

type EnvType string

const (
	EnvTypeAll   EnvType = "ALL"
	EnvTypeAllow EnvType = "PROD"
)

type Rule struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Type        RuleType   `json:"type" db:"type"`
	BU          string     `json:"bu" db:"bu"`
	Service     string     `json:"service" db:"service"`
	Environment EnvType    `json:"environment" db:"environment"`
	Action      ActionType `json:"action" db:"action"`

	Realms   []string `json:"realms,omitempty" db:"realms"`
	Accounts []string `json:"accounts,omitempty" db:"accounts"`

	Message   string    `json:"message,omitempty" db:"message"`
	CreatedBy string    `json:"createdBy,omitempty" db:"created_by"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`

	UpdatedBy *string    `json:"updatedBy,omitempty" db:"updated_by"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

func ScanRule(rows pgx.Rows) (Rule, error) {
	var r Rule
	err := rows.Scan(
		&r.ID,
		&r.Type,
		&r.BU,
		&r.Service,
		&r.Environment,
		&r.Action,
		&r.Realms,
		&r.Accounts,
		&r.Message,
		&r.CreatedBy,
		&r.CreatedAt,
		&r.UpdatedBy,
		&r.UpdatedAt,
	)
	return r, err
}

// ruleRepo := NewBaseRepository[Rule](
// 	db,
// 	"rules",
// 	[]string{
// 		"id",
// 		"type",
// 		"bu",
// 		"service",
// 		"environment",
// 		"action",
// 		"realms",
// 		"accounts",
// 		"message",
// 		"created_by",
// 		"created_at",
// 		"updated_by",
// 		"updated_at",
// 	},
// 	ScanRule,
// )
