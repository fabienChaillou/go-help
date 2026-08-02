package rule

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type DBTX interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Rule struct {
	ID          string
	Type        string
	BU          string
	Service     string
	Environment string
	Action      string
	Realms      []string
	Accounts    []string
	Message     string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   string
	UpdatedAt   time.Time
}

type RuleRepository struct {
	db DBTX // interface pour mocker
}

func (r *RuleRepository) FindBestRule(
	ctx context.Context,
	bu, service,
	env1, env2,
	action1, action2 string,
) (*Rule, error) {

	query := `
SELECT * FROM rules
WHERE bu = $1
AND service = $2
AND (environment = $3 OR environment = $4)
AND (action = $5 OR action = $6)
ORDER BY CASE WHEN type='ALLOW' THEN 0 ELSE 1 END,
created_at DESC
LIMIT 1
`

	row := r.db.QueryRow(ctx, query,
		bu, service,
		env1, env2,
		action1, action2,
	)

	var rule Rule
	err := row.Scan(
		&rule.ID,
		&rule.Type,
		&rule.BU,
		&rule.Service,
		&rule.Environment,
		&rule.Action,
		&rule.Realms,
		&rule.Accounts,
		&rule.Message,
		&rule.CreatedBy,
		&rule.CreatedAt,
		&rule.UpdatedBy,
		&rule.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &rule, nil
}
