// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: group.sql

package db

import (
	"context"
)

const createGroup = `-- name: CreateGroup :one

INSERT INTO groups (
  name,
  category_id,
  created_by_id,
  image_path
) VALUES (
  $1, $2, $3, $4
) RETURNING id, name, category_id, image_path, created_by_id, created_at
`

type CreateGroupParams struct {
	Name        string `json:"name"`
	CategoryID  int64  `json:"category_id"`
	CreatedByID int64  `json:"created_by_id"`
	ImagePath   string `json:"image_path"`
}

// Groups
// Create a group
func (q *Queries) CreateGroup(ctx context.Context, arg CreateGroupParams) (Groups, error) {
	row := q.db.QueryRowContext(ctx, createGroup,
		arg.Name,
		arg.CategoryID,
		arg.CreatedByID,
		arg.ImagePath,
	)
	var i Groups
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CategoryID,
		&i.ImagePath,
		&i.CreatedByID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteGroup = `-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1
`

// Delete a group by ID
func (q *Queries) DeleteGroup(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteGroup, id)
	return err
}

const getGroupByID = `-- name: GetGroupByID :one
SELECT id, name, category_id, image_path, created_by_id, created_at FROM groups
WHERE id = $1 LIMIT 1
`

// Get a group by ID
func (q *Queries) GetGroupByID(ctx context.Context, id int64) (Groups, error) {
	row := q.db.QueryRowContext(ctx, getGroupByID, id)
	var i Groups
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CategoryID,
		&i.ImagePath,
		&i.CreatedByID,
		&i.CreatedAt,
	)
	return i, err
}

const listGroups = `-- name: ListGroups :many
SELECT id, name, category_id, image_path, created_by_id, created_at FROM groups
LIMIT $1 OFFSET $2
`

type ListGroupsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

// List groups with pagination
func (q *Queries) ListGroups(ctx context.Context, arg ListGroupsParams) ([]Groups, error) {
	rows, err := q.db.QueryContext(ctx, listGroups, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Groups{}
	for rows.Next() {
		var i Groups
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CategoryID,
			&i.ImagePath,
			&i.CreatedByID,
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

const updateGroup = `-- name: UpdateGroup :one
UPDATE groups SET
  name = $2,
  category_id = $3,
  created_by_id = $4
WHERE id = $1 RETURNING id, name, category_id, image_path, created_by_id, created_at
`

type UpdateGroupParams struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	CategoryID  int64  `json:"category_id"`
	CreatedByID int64  `json:"created_by_id"`
}

// Update a group by ID
func (q *Queries) UpdateGroup(ctx context.Context, arg UpdateGroupParams) (Groups, error) {
	row := q.db.QueryRowContext(ctx, updateGroup,
		arg.ID,
		arg.Name,
		arg.CategoryID,
		arg.CreatedByID,
	)
	var i Groups
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CategoryID,
		&i.ImagePath,
		&i.CreatedByID,
		&i.CreatedAt,
	)
	return i, err
}