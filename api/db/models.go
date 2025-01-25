// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Board struct {
	BoardID     pgtype.Int8
	Name        string
	Description pgtype.Text
	UserID      pgtype.Int8
	CreatedDate pgtype.Timestamptz
}

type ColumnItem struct {
	ColumnItemID pgtype.Int8
	ColumnID     pgtype.Int8
	ItemID       pgtype.Int8
	UserID       pgtype.Int8
	CreatedDate  pgtype.Timestamptz
	Position     pgtype.Int4
}

type Item struct {
	ItemID      pgtype.Int8
	Title       string
	StatusID    pgtype.Int8
	CreatedDate pgtype.Timestamptz
}

type Kbcolumn struct {
	ColumnID    pgtype.Int8
	Name        string
	BoardID     pgtype.Int8
	UserID      pgtype.Int8
	Position    int32
	CreatedDate pgtype.Timestamptz
}

type Review struct {
	ReviewID    pgtype.Int8
	Content     string
	UserID      pgtype.Int8
	ItemID      pgtype.Int8
	CreatedDate pgtype.Timestamptz
}

type Status struct {
	StatusID    pgtype.Int8
	UserID      pgtype.Int8
	Label       pgtype.Text
	CreatedDate pgtype.Timestamptz
}

type User struct {
	UserID         pgtype.Int8
	Username       string
	HashedPassword string
	Email          string
	FullName       pgtype.Text
	Bio            pgtype.Text
	CreatedDate    pgtype.Timestamptz
}
