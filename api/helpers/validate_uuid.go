package helpers

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

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
