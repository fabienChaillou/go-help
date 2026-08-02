package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// func (r *RuleRepository) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
// 	return r.db.Query(ctx, sql, args...)
// }

type Rule struct {
	ID          int64
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
	UpdateBy    string
}

type RuleRepository struct {
	db Queryer
	// db *pgxpool.Pool
}

func NewRuleRepository(db Queryer) *RuleRepository {
	return &RuleRepository{db: db}
}

func (r *RuleRepository) ListByBU(ctx context.Context, bu string) ([]Rule, error) {

	builder := sq.
		Select(
			"id",
			"type",
			"bu",
			"service",
			"environment",
			"action",
			"realms",
			"accounts",
			"message",
			"created_by",
			"created_at",
			"update_by",
		).
		From("rules").
		Where(sq.Eq{"bu": bu}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []Rule

	for rows.Next() {
		var rule Rule
		err := rows.Scan(
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
			&rule.UpdateBy,
		)
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	return rules, rows.Err()
}
