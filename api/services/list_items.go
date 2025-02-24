package services

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/types"
	"context"
	"net/http"
)

type ListItemsService struct {
	q *queries.Queries
}

func NewListItemsService(q *queries.Queries) *ListItemsService {
	return &ListItemsService{q: q}
}

// GetAllListItems fetches all list items from the database and returns them as a paginated list
func (s *ListItemsService) GetAllListItems(ctx context.Context, pagination *types.Pagination) (*types.PaginatedListItemsResponse, error) {
	total, err := s.q.GetAllListItemsCount(ctx)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error retrieving list items")
	}

	items, err := s.q.GetAllListItems(ctx, queries.GetAllListItemsParams{
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	})
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error retrieving list items")
	}

	if total == 0 {
		response := &types.PaginatedListItemsResponse{
			Pagination: *pagination,
			ListItems:  []types.ListItemsResponse{},
		}
		return response, nil
	}

	pagination.Total = total

	listItemsResponse := make([]types.ListItemsResponse, len(items))

	for i, item := range items {
		listItemsResponse[i] = types.ListItemsResponse{
			UUID:     item.ListItemUuid.String(),
			ListUUID: item.ListUuid.String(),
			ItemUUID: item.ItemUuid.String(),
			Status:   item.Label.String,
			Position: item.Position,
		}
	}

	response := types.PaginatedListItemsResponse{
		Pagination: *pagination,
		ListItems:  listItemsResponse,
	}

	return &response, nil
}

// GetListItemsForList returns all items in a list
func (s *ListItemsService) GetListItemsForList(ctx context.Context, pagination *types.Pagination, uuid string) (*types.PaginatedListItemsResponse, error) {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return nil, err
	}

	total, err := s.q.GetListItemsCountForList(ctx, *pgUuid)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error fetching list items count")
	}

	items, err := s.q.GetListItemsByListUuid(ctx, queries.GetListItemsByListUuidParams{
		ListUuid: *pgUuid,
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	})
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error fetching list items")
	}

	pagination.Total = total

	if total == 0 {
		response := &types.PaginatedListItemsResponse{
			Pagination: *pagination,
			ListItems:  []types.ListItemsResponse{},
		}
		return response, nil
	}

	listItemsResponse := make([]types.ListItemsResponse, len(items))

	for i, item := range items {
		listItemsResponse[i] = types.ListItemsResponse{
			UUID:     item.ListItemUuid.String(),
			ListUUID: item.ListUuid.String(),
			ItemUUID: item.ItemUuid.String(),
			Status:   item.Label.String,
			Position: item.Position,
		}
	}

	response := types.PaginatedListItemsResponse{
		Pagination: *pagination,
		ListItems:  listItemsResponse,
	}

	return &response, nil
}
