// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: group_category.sql

package db

import (
	"context"
)

const createGroupCategory = `-- name: CreateGroupCategory :one
INSERT INTO group_categories (
  name
) VALUES (
  $1
) RETURNING id, name, created_at
`

func (q *Queries) CreateGroupCategory(ctx context.Context, name string) (GroupCategories, error) {
	row := q.db.QueryRowContext(ctx, createGroupCategory, name)
	var i GroupCategories
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const deleteGroupCategory = `-- name: DeleteGroupCategory :exec
DELETE FROM group_categories
WHERE id = $1
`

func (q *Queries) DeleteGroupCategory(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteGroupCategory, id)
	return err
}

const getGroupCategory = `-- name: GetGroupCategory :one
SELECT id, name, created_at FROM group_categories
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetGroupCategory(ctx context.Context, id int64) (GroupCategories, error) {
	row := q.db.QueryRowContext(ctx, getGroupCategory, id)
	var i GroupCategories
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const listGroupCategories = `-- name: ListGroupCategories :many
SELECT id, name, created_at FROM group_categories
LIMIT $1 OFFSET $2
`

type ListGroupCategoriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListGroupCategories(ctx context.Context, arg ListGroupCategoriesParams) ([]GroupCategories, error) {
	rows, err := q.db.QueryContext(ctx, listGroupCategories, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GroupCategories{}
	for rows.Next() {
		var i GroupCategories
		if err := rows.Scan(&i.ID, &i.Name, &i.CreatedAt); err != nil {
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

const updateGroupCategory = `-- name: UpdateGroupCategory :one
UPDATE group_categories SET
  name = $2
WHERE id = $1 RETURNING id, name, created_at
`

type UpdateGroupCategoryParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) UpdateGroupCategory(ctx context.Context, arg UpdateGroupCategoryParams) (GroupCategories, error) {
	row := q.db.QueryRowContext(ctx, updateGroupCategory, arg.ID, arg.Name)
	var i GroupCategories
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}
