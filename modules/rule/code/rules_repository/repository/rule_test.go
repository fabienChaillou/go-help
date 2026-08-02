package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"

	"rule-repository/repository"
)

func TestListByBU(t *testing.T) {

	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewRuleRepository(mock)

	// ⚠ On reconstruit EXACTEMENT la même requête Squirrel
	// builder := sq.
	// 	Select(
	// 		"id",
	// 		// "type",
	// 		"bu",
	// 		// "service",
	// 		// "environment",
	// 		// "action",
	// 		// "realms",
	// 		// "accounts",
	// 		// "message",
	// 		// "created_by",
	// 		// "created_at",
	// 		// "update_by",
	// 	).
	// 	From("rules").
	// 	Where(sq.Eq{"bu": "BU1"}).
	// 	// OrderBy("created_at DESC").
	// 	PlaceholderFormat(sq.Dollar)

	// query, args, err := builder.ToSql()
	// assert.NoError(t, err)

	// query := "SELECT id, bu FROM rules WHERE bu=\\$1"
	args := []any{"BU1"}

	now := time.Now()

	rows := pgxmock.NewRows([]string{
		"id", "type", "bu", "service", "environment", "action",
		"realms", "accounts", "message",
		"created_by", "created_at", "update_by",
	}).AddRow(
		int64(1),
		"security",
		"BU1",
		"svc1",
		"prod",
		"allow",
		[]string{"realm1"},
		[]string{"acc1"},
		"test message",
		"admin",
		now,
		"admin",
	)

	mock.ExpectQuery("SELECT (.+) FROM rules").
		WithArgs(args...).
		WillReturnRows(rows)

	result, err := repo.ListByBU(context.Background(), "BU1")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), result[0].ID)
	assert.Equal(t, "BU1", result[0].BU)
	assert.Equal(t, "security", result[0].Type)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryMatcherRegexp(t *testing.T) {
	type testCase struct {
		expected string
		actual   string
		err      error
	}

	cases := []testCase{
		{"?\\l", "SEL", fmt.Errorf("error parsing regexp: missing argument to repetition operator: `?`")},
		{"SELECT (.+) FROM users", "SELECT name, email FROM users WHERE id = ?", nil},
		{"Select (.+) FROM users", "SELECT name, email FROM users WHERE id = ?", fmt.Errorf(`could not match actual sql: "SELECT name, email FROM users WHERE id = ?" with expected regexp "Select (.+) FROM users"`)},
		{"SELECT (.+) FROM\nusers", "SELECT name, email\n FROM users\n WHERE id = ?", nil},
		{"SELECT (.+) FROM rules", "SELECT id, bu FROM rules WHERE bu = $1", nil},
	}

	for i, c := range cases {
		err := pgxmock.QueryMatcherRegexp.Match(c.expected, c.actual)
		if err == nil && c.err != nil {
			t.Errorf(`got no error, but expected "%v" at %d case`, c.err, i)
			continue
		}
		if err != nil && c.err == nil {
			t.Errorf(`got unexpected error "%v" at %d case`, err, i)
			continue
		}
		if err == nil {
			continue
		}
		if err.Error() != c.err.Error() {
			t.Errorf(`expected error "%v", but got "%v" at %d case`, c.err, err, i)
		}
	}
}
