package services

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/types"
	"context"
	"errors"
)

type ListItemsService struct {
	q *queries.Queries
}

func NewListItemsService(q *queries.Queries) *ListItemsService {
	return &ListItemsService{q: q}
}

// GetAllListItems fetches all list items from the database and returns them as a paginated list
func (s *ListItemsService) GetAllListItems(ctx context.Context, pagination *types.Pagination) ([]queries.GetAllListItemsRow, *types.Pagination, error) {
	total, err := s.q.GetAllListItemsCount(ctx)
	if err != nil {
		return nil, pagination, errors.New("failed to query all list items count")
	}

	pagination.Total = total

	listItems, err := s.q.GetAllListItems(ctx, queries.GetAllListItemsParams{
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	})
	if err != nil {
		return nil, pagination, errors.New("failed to query all list items")
	}

	return listItems, pagination, nil
}

// GetListItemsForList returns all items in a list
func (s *ListItemsService) GetListItemsForList(ctx context.Context, pagination *types.Pagination, uuid string) ([]queries.GetListItemsByListUuidRow, *types.Pagination, error) {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return nil, pagination, errors.New("invalid uuid")
	}

	total, err := s.q.GetListItemsCountForList(ctx, pgUuid)
	if err != nil {
		return nil, pagination, errors.New("failed to query list items count")
	}

	pagination.Total = total

	listItems, err := s.q.GetListItemsByListUuid(ctx, queries.GetListItemsByListUuidParams{
		ListUuid: pgUuid,
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	})
	if err != nil {
		return nil, pagination, errors.New("failed to query list items")
	}

	return listItems, pagination, nil
}
