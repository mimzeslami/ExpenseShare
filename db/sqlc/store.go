package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	CreateGroupTx(ctx context.Context, arg CreateGroupParams) (CreateGroupTxResults, error)
	DeleteGroupTx(ctx context.Context, arg DeleteGroupParams) error
}

type SQLStroe struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStroe{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStroe) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type CreateGroupTxResults struct {
	Group        Groups
	GroupMembers GroupMembers
}

func (store *SQLStroe) CreateGroupTx(ctx context.Context, arg CreateGroupParams) (CreateGroupTxResults, error) {
	result := CreateGroupTxResults{}
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Group, err = q.CreateGroup(ctx, arg)
		result.GroupMembers, err = q.CreateGroupMember(ctx, CreateGroupMemberParams{
			GroupID: result.Group.ID,
			UserID:  arg.CreatedByID,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}

func (store *SQLStroe) DeleteGroupTx(ctx context.Context, arg DeleteGroupParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		arg := DeleteGroupParams{
			ID:          arg.ID,
			CreatedByID: arg.CreatedByID,
		}
		err = q.DeleteGroupMembers(ctx, arg.ID)
		if err != nil {
			return err
		}
		err = q.DeleteGroup(ctx, arg)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}
