package helpers

import (
	"codeberg.org/sporiff/eigakanban/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

func MakePgString(s string) pgtype.Text {
	if len(s) == 0 {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func ValidateAndConvertUUID(uuidString string) (*pgtype.UUID, error) {
	parsedUuid, err := uuid.Parse(uuidString)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error validating uuid")
	}

	result := pgtype.UUID{
		Bytes: parsedUuid,
		Valid: true,
	}

	return &result, nil
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
