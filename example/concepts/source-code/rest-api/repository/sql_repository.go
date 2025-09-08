package repository

import (
	"fmt"
	"reflect"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type SQLRepository[T any] struct {
	DB    *sqlx.DB
	Table string
}

func (r *SQLRepository[T]) GetByID(id int) (*T, error) {
	var t T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", r.Table)
	err := r.DB.Get(&t, query, id)
	return &t, err
}

func (r *SQLRepository[T]) ListPaginated(limit, offset int) ([]T, error) {
	var items []T
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id LIMIT $1 OFFSET $2", r.Table)
	err := r.DB.Select(&items, query, limit, offset)
	return items, err
}

func (r *SQLRepository[T]) Create(entity *T) error {
	val := reflect.Indirect(reflect.ValueOf(entity))
	typ := val.Type()

	values := map[string]interface{}{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" || dbTag == "id" {
			continue
		}
		values[dbTag] = val.Field(i).Interface()
	}

	query, args, err := sq.Insert(r.Table).SetMap(values).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(query, args...)
	return err
}

func (r *SQLRepository[T]) Update(id int, entity *T) error {
	val := reflect.Indirect(reflect.ValueOf(entity))
	typ := val.Type()
	values := map[string]interface{}{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" || dbTag == "id" {
			continue
		}
		values[dbTag] = val.Field(i).Interface()
	}
	query, args, err := sq.Update(r.Table).SetMap(values).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(query, args...)
	return err
}

func (r *SQLRepository[T]) Delete(id int) error {
	query, args, err := sq.Delete(r.Table).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(query, args...)
	return err
}
