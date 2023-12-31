// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: expense.sql

package db

import (
	"context"
)

const createExpense = `-- name: CreateExpense :one

INSERT INTO expenses (
  group_id,
  paid_by_id,
  amount,
  description
) VALUES (
  $1, $2, $3, $4
) RETURNING id, group_id, paid_by_id, amount, description, created_at
`

type CreateExpenseParams struct {
	GroupID     int64  `json:"group_id"`
	PaidByID    int64  `json:"paid_by_id"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
}

// expenses.sql
// Create an expense
func (q *Queries) CreateExpense(ctx context.Context, arg CreateExpenseParams) (Expenses, error) {
	row := q.db.QueryRowContext(ctx, createExpense,
		arg.GroupID,
		arg.PaidByID,
		arg.Amount,
		arg.Description,
	)
	var i Expenses
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.PaidByID,
		&i.Amount,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}

const deleteExpense = `-- name: DeleteExpense :exec
DELETE FROM expenses
WHERE id = $1
`

// Delete an expense by ID
func (q *Queries) DeleteExpense(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteExpense, id)
	return err
}

const getExpenseByID = `-- name: GetExpenseByID :one
SELECT id, group_id, paid_by_id, amount, description, created_at FROM expenses
WHERE id = $1 LIMIT 1
`

// Get an expense by ID
func (q *Queries) GetExpenseByID(ctx context.Context, id int64) (Expenses, error) {
	row := q.db.QueryRowContext(ctx, getExpenseByID, id)
	var i Expenses
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.PaidByID,
		&i.Amount,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}

const listExpenses = `-- name: ListExpenses :many
SELECT id, group_id, paid_by_id, amount, description, created_at FROM expenses
WHERE group_id = $1
LIMIT $2 OFFSET $3
`

type ListExpensesParams struct {
	GroupID int64 `json:"group_id"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
}

// List expenses for a group with pagination
func (q *Queries) ListExpenses(ctx context.Context, arg ListExpensesParams) ([]Expenses, error) {
	rows, err := q.db.QueryContext(ctx, listExpenses, arg.GroupID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Expenses{}
	for rows.Next() {
		var i Expenses
		if err := rows.Scan(
			&i.ID,
			&i.GroupID,
			&i.PaidByID,
			&i.Amount,
			&i.Description,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateExpense = `-- name: UpdateExpense :one
UPDATE expenses SET
  amount = $2,
  description = $3
WHERE id = $1 RETURNING id, group_id, paid_by_id, amount, description, created_at
`

type UpdateExpenseParams struct {
	ID          int64  `json:"id"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
}

// Update an expense by ID
func (q *Queries) UpdateExpense(ctx context.Context, arg UpdateExpenseParams) (Expenses, error) {
	row := q.db.QueryRowContext(ctx, updateExpense, arg.ID, arg.Amount, arg.Description)
	var i Expenses
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.PaidByID,
		&i.Amount,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}
