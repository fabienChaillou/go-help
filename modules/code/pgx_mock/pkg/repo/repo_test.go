package repo

import (
	"context"
	"fmt"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	rows := pgxmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Alice")

	mock.ExpectQuery("SELECT id, name FROM users WHERE id=\\$1").
		WithArgs(1).
		WillReturnRows(rows)

	repo := &UserRepository{db: mock}

	user, err := repo.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, "Alice", user.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByName(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	rows := pgxmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Bob")

	mock.ExpectQuery("SELECT id, name FROM users WHERE name=\\$1").
		WithArgs("Bob").
		WillReturnRows(rows)

	repo := &UserRepository{db: mock}

	result, err := repo.GetByName(context.Background(), "Bob")

	assert.NoError(t, err)
	assert.Len(t, result, 1)

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
