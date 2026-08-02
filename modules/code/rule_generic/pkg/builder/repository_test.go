package builder

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestRuleRepository_Find_WithLimit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewBaseRepository[Rule](
		mock,
		"rules",
		[]string{
			"id", "type", "bu", "service", "environment", "action",
			"realms", "accounts", "message", "created_by",
			"created_at", "updated_by", "updated_at",
		},
		ScanRule,
	)

	id := uuid.New()

	rows := pgxmock.NewRows([]string{
		"id", "type", "bu", "service", "environment", "action",
		"realms", "accounts", "message", "created_by",
		"created_at", "updated_by", "updated_at",
	}).AddRow(
		id, "ALLOW", "finance", "billing", "prod", "read",
		[]string{}, []string{}, "msg", "admin",
		time.Now(), nil, nil,
	)

	mock.ExpectQuery("SELECT (.+) FROM rules").
		WillReturnRows(rows)

	result, err := repo.Find(context.Background(), QueryOptions{
		Filters: []Filter{
			{Field: "bu", Op: "=", Value: "finance"},
		},
		Limit: 1,
	})

	fmt.Println("foo:", result)
	fmt.Println("rows:", rows)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, id, result[0].ID)
}

func TestRuleRepository_InvalidField(t *testing.T) {
	t.Skip("TODO: fix the TestRuleRepository_Find_WithLimit func!")
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewBaseRepository[Rule](
		mock,
		"rules",
		[]string{"id", "type"},
		ScanRule,
	)

	_, err := repo.Find(context.Background(), QueryOptions{
		Filters: []Filter{
			{Field: "DROP TABLE", Op: "=", Value: "x"},
		},
	})

	assert.Error(t, err)
}
