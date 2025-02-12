package helpers

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func MakePgString(s string) pgtype.Text {
	if len(s) == 0 {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func ValidateAndConvertUUID(uuidString string) (pgtype.UUID, error) {
	pgUuid := pgtype.UUID{}
	parsedUuid, err := uuid.Parse(uuidString)
	if err != nil {
		return pgUuid, err
	}
	copy(pgUuid.Bytes[:], parsedUuid[:])
	pgUuid.Valid = true
	return pgUuid, nil
}

// AssignIfNotNil assigns values to a field if its current value is nil
func AssignIfNotNil[T any](dest *T, src *T, defaultValue T) {
	if src != nil {
		*dest = *src
	} else {
		*dest = defaultValue
	}
}

// AssignPgtypeText assigns postgres text type values assigns values to a field if its current value is nil
func AssignPgtypeText(dest *pgtype.Text, src *string) {
	if src != nil {
		dest.String = *src
		dest.Valid = true
	} else {
		dest.Valid = false
	}
}
