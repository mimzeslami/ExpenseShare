// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"database/sql"
	"time"
)

type Currencies struct {
	ID           int64     `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Symbol       string    `json:"symbol"`
	ExchangeRate float64   `json:"exchange_rate"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ExpenseShares struct {
	ID         int64     `json:"id"`
	ExpenseID  int64     `json:"expense_id"`
	UserID     int64     `json:"user_id"`
	Share      string    `json:"share"`
	PaidStatus bool      `json:"paid_status"`
	CreatedAt  time.Time `json:"created_at"`
}

type Expenses struct {
	ID          int64     `json:"id"`
	GroupID     int64     `json:"group_id"`
	PaidByID    int64     `json:"paid_by_id"`
	Amount      string    `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupCategories struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupMembers struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Groups struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	CategoryID  int64     `json:"category_id"`
	ImagePath   string    `json:"image_path"`
	CreatedByID int64     `json:"created_by_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Invitations struct {
	ID         int64        `json:"id"`
	InviterID  int64        `json:"inviter_id"`
	InviteeID  int64        `json:"invitee_id"`
	GroupID    int64        `json:"group_id"`
	Status     string       `json:"status"`
	Code       string       `json:"code"`
	CreatedAt  time.Time    `json:"created_at"`
	AcceptedAt sql.NullTime `json:"accepted_at"`
	RejectedAt sql.NullTime `json:"rejected_at"`
}

type Notifications struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type Users struct {
	ID           int64     `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Phone        string    `json:"phone"`
	ImagePath    string    `json:"image_path"`
	TimeZone     string    `json:"time_zone"`
	CreatedAt    time.Time `json:"created_at"`
}
