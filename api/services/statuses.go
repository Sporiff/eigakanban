package services

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/types"
	"context"
	"errors"
)

type StatusesService struct {
	q *queries.Queries
}

func NewStatusesService(q *queries.Queries) *StatusesService {
	return &StatusesService{q: q}
}

// AddStatus adds a new status to the database
func (s *StatusesService) AddStatus(ctx context.Context, status types.AddStatusRequest, uuid string) (queries.AddStatusRow, error) {
	emptyRow := queries.AddStatusRow{}

	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return emptyRow, errors.New("error validating uuid: " + err.Error())
	}

	statusRow, err := s.q.AddStatus(ctx, queries.AddStatusParams{
		StatusLabel: helpers.MakePgString(status.StatusLabel),
		UserUuid:    pgUuid,
	})
	if err != nil {
		return emptyRow, errors.New("error adding status: " + err.Error())
	}

	return statusRow, err
}
