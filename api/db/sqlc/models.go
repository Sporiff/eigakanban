// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Item struct {
	ItemID      pgtype.Int8        `json:"item_id"`
	Uuid        pgtype.UUID        `json:"uuid"`
	Title       string             `json:"title"`
	CreatedDate pgtype.Timestamptz `json:"created_date"`
}

type List struct {
	ListID      pgtype.Int8        `json:"list_id"`
	Uuid        pgtype.UUID        `json:"uuid"`
	Name        string             `json:"name"`
	UserID      int64              `json:"user_id"`
	CreatedDate pgtype.Timestamptz `json:"created_date"`
}

type ListItem struct {
	ListItemID  pgtype.Int8        `json:"list_item_id"`
	Uuid        pgtype.UUID        `json:"uuid"`
	ListID      int64              `json:"list_id"`
	ItemID      int64              `json:"item_id"`
	Position    int32              `json:"position"`
	PrevItemID  pgtype.Int8        `json:"prev_item_id"`
	NextItemID  pgtype.Int8        `json:"next_item_id"`
	StatusID    pgtype.Int8        `json:"status_id"`
	CreatedDate pgtype.Timestamptz `json:"created_date"`
}

type ListStatus struct {
	ListStatusID pgtype.Int8        `json:"list_status_id"`
	Uuid         pgtype.UUID        `json:"uuid"`
	ListID       int64              `json:"list_id"`
	StatusID     int64              `json:"status_id"`
	CreatedDate  pgtype.Timestamptz `json:"created_date"`
}

type RefreshToken struct {
	TokenID   pgtype.Int8        `json:"token_id"`
	UserID    int64              `json:"user_id"`
	Token     string             `json:"token"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type Review struct {
	ReviewID    pgtype.Int8        `json:"review_id"`
	Uuid        pgtype.UUID        `json:"uuid"`
	Content     string             `json:"content"`
	UserID      int64              `json:"user_id"`
	ItemID      int64              `json:"item_id"`
	CreatedDate pgtype.Timestamptz `json:"created_date"`
}

type Status struct {
	StatusID    pgtype.Int8        `json:"status_id"`
	Uuid        pgtype.UUID        `json:"uuid"`
	UserID      pgtype.Int8        `json:"user_id"`
	Label       pgtype.Text        `json:"label"`
	CreatedDate pgtype.Timestamptz `json:"created_date"`
}

type User struct {
	UserID         pgtype.Int8        `json:"user_id"`
	Uuid           pgtype.UUID        `json:"uuid"`
	Username       string             `json:"username"`
	HashedPassword string             `json:"hashed_password"`
	Email          string             `json:"email"`
	FullName       pgtype.Text        `json:"full_name"`
	Bio            pgtype.Text        `json:"bio"`
	Superuser      bool               `json:"superuser"`
	CreatedDate    pgtype.Timestamptz `json:"created_date"`
}
