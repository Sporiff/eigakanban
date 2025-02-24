package services

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/types"
	"context"
	"net/http"
)

type StatusesService struct {
	q *queries.Queries
}

func NewStatusesService(q *queries.Queries) *StatusesService {
	return &StatusesService{q: q}
}

// AddStatus adds a new status to the database
func (s *StatusesService) AddStatus(ctx context.Context, status types.AddStatusRequest, uuid string) (*queries.AddStatusRow, error) {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return nil, err
	}

	result, err := s.q.AddStatus(ctx, queries.AddStatusParams{
		StatusLabel: helpers.MakePgString(status.StatusLabel),
		UserUuid:    *pgUuid,
	})
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error adding status")
	}

	return &result, err
}

// GetStatusesForUser retrieves all statuses for the authenticated user
func (s *StatusesService) GetStatusesForUser(ctx context.Context, uuid string, pagination *types.Pagination) (*types.PaginatedStatusesResponse, error) {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return nil, err
	}

	total, err := s.q.GetStatusesCountForUser(ctx, *pgUuid)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error fetching status count")
	}

	pagination.Total = total

	if total == 0 {
		response := &types.PaginatedStatusesResponse{
			Pagination: *pagination,
			Statuses:   []types.StatusesResponse{},
		}
		return response, nil
	}

	items, err := s.q.GetStatusesForUser(ctx, queries.GetStatusesForUserParams{
		UserUuid: *pgUuid,
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	})
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error fetching statuses")
	}

	statuses := make([]types.StatusesResponse, len(items))

	for i, item := range items {
		statuses[i] = types.StatusesResponse{
			UUID:  item.Uuid.String(),
			Label: item.Label.String,
		}
	}

	response := types.PaginatedStatusesResponse{
		Pagination: *pagination,
		Statuses:   statuses,
	}

	return &response, nil
}
