package rule

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestFindBestRule(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &RuleRepository{db: mock}

	now := time.Now()

	rows := pgxmock.NewRows([]string{
		"i_d",
		"type",
		"b_u",
		"service",
		"environment",
		"action",
		"realms",
		"accounts",
		"message",
		"created_by",
		"created_at",
		"updated_by",
		"updated_at",
	}).AddRow(
		"rule-1",
		"ALLOW",
		"BU1",
		"SVC1",
		"prod",
		"read",
		[]string{"realm1"},
		[]string{"acc1"},
		"ok",
		"admin",
		now,
		"admin",
		now,
	)

	mock.ExpectQuery(`SELECT \* FROM rules`).
		WithArgs(
			"BU1",
			"SVC1",
			"prod",
			"default",
			"read",
			"write",
		).
		WillReturnRows(rows)

	rule, err := repo.FindBestRule(
		context.Background(),
		"BU1",
		"SVC1",
		"prod",
		"default",
		"read",
		"write",
	)

	assert.NoError(t, err)
	assert.Equal(t, "rule-1", rule.ID)
	assert.Equal(t, "ALLOW", rule.Type)

	assert.NoError(t, mock.ExpectationsWereMet())
}
