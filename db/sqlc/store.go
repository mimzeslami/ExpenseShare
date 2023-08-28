package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mimzeslami/expense_share/util"
)

type Store interface {
	Querier
	CreateGroupTx(ctx context.Context, arg CreateGroupParams) (CreateGroupTxResults, error)
	DeleteGroupTx(ctx context.Context, arg DeleteGroupParams) error
	AddUserToGroupTx(ctx context.Context, arg AddUserToGroupParams) (AddUserToGroupResults, error)
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

type AddUserToGroupParams struct {
	GroupID      int64  `json:"group_id" binding:"required"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	TimeZone     string `json:"time_zone"`
	GroupOwnerID int64  `json:"group_owner_id" binding:"required"`
}

type AddUserToGroupResults struct {
	GroupMembers GroupMembers `json:"group_members"`
	User         Users        `json:"user"`
	Invitations  Invitations  `json:"invitations"`
}

func (store *SQLStroe) AddUserToGroupTx(ctx context.Context, arg AddUserToGroupParams) (AddUserToGroupResults, error) {
	result := AddUserToGroupResults{}
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		hashedPassword, err := util.HashPassword(util.RandomString(6))
		if err != nil {
			return err
		}

		user, err := q.GetUserByPhone(ctx, arg.Phone)

		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if user == (Users{}) {
			result.User, err = q.CreateUser(ctx, CreateUserParams{
				FirstName:    arg.FirstName,
				LastName:     arg.LastName,
				Email:        arg.Email,
				Phone:        arg.Phone,
				TimeZone:     arg.TimeZone,
				PasswordHash: hashedPassword,
			})

			if err != nil {
				return err
			}
		} else {
			result.User = user
		}

		result.Invitations, err = q.CreateInvitation(ctx, CreateInvitationParams{
			InviterID: arg.GroupOwnerID,
			InviteeID: result.User.ID,
			GroupID:   arg.GroupID,
			Status:    "Pending",
			Code:      util.RandomString(6),
		})

		if err != nil {
			return err
		}

		result.GroupMembers, err = q.CreateGroupMember(ctx, CreateGroupMemberParams{
			GroupID: arg.GroupID,
			UserID:  result.User.ID,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}
