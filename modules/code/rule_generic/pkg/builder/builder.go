package builder

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type Filter struct {
	Field string
	Op    string
	Value any
}

type Sort struct {
	Field string
	Desc  bool
}

type QueryOptions struct {
	Filters []Filter
	OrderBy []Sort
	Limit   uint64
	Offset  uint64
}

type DBTX interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type BaseRepository[T any] struct {
	db            DBTX
	table         string
	columns       []string
	allowedFields map[string]struct{}
	scanFunc      func(pgx.Rows) (T, error)
}

func NewBaseRepository[T any](
	db DBTX,
	table string,
	columns []string,
	scanFunc func(pgx.Rows) (T, error),
) *BaseRepository[T] {

	allowed := make(map[string]struct{})
	for _, c := range columns {
		allowed[c] = struct{}{}
	}

	return &BaseRepository[T]{
		db:            db,
		table:         table,
		columns:       columns,
		allowedFields: allowed,
		scanFunc:      scanFunc,
	}
}

func (r *BaseRepository[T]) Find(
	ctx context.Context,
	opts QueryOptions,
) ([]T, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	qb := psql.
		Select(r.columns...).
		From(r.table)

	// WHERE
	for _, f := range opts.Filters {

		if _, ok := r.allowedFields[f.Field]; !ok {
			return nil, fmt.Errorf("field not allowed: %s", f.Field)
		}

		qb = qb.Where(
			squirrel.Expr(f.Field+" "+f.Op+" ?", f.Value),
		)
	}

	// ORDER BY
	for _, s := range opts.OrderBy {

		if _, ok := r.allowedFields[s.Field]; !ok {
			return nil, fmt.Errorf("field not allowed: %s", s.Field)
		}

		order := s.Field
		if s.Desc {
			order += " DESC"
		}

		qb = qb.OrderBy(order)
	}

	if opts.Limit > 0 {
		qb = qb.Limit(opts.Limit)
	}

	if opts.Offset > 0 {
		qb = qb.Offset(opts.Offset)
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	fmt.Println("query: ", sql)

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []T
	for rows.Next() {
		entity, err := r.scanFunc(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, entity)
	}

	return result, rows.Err()
}
