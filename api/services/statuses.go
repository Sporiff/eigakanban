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
		return emptyRow, err
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

// GetStatusesForUser retrieves all statuses for the authenticated user
func (s *StatusesService) GetStatusesForUser(ctx context.Context, uuid string, pagination *types.Pagination) ([]queries.GetStatusesForUserRow, *types.Pagination, error) {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return nil, nil, err
	}

	total, err := s.q.GetStatusesCountForUser(ctx, pgUuid)
	if err != nil {
		return nil, nil, err
	}

	pagination.Total = total

	rows, err := s.q.GetStatusesForUser(ctx, queries.GetStatusesForUserParams{
		UserUuid: pgUuid,
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	})
	if err != nil {
		return nil, nil, err
	}

	return rows, pagination, nil
}
