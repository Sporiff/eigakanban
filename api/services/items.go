package services

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/types"
	"context"
	"database/sql"
	"errors"
	"net/http"
)

type ItemsService struct {
	q *queries.Queries
}

func NewItemsService(q *queries.Queries) *ItemsService {
	return &ItemsService{q: q}
}

// GetAllItems retrieves all items from the database as a paginated list
func (s *ItemsService) GetAllItems(ctx context.Context, pagination *types.Pagination) (*types.PaginatedItemsResponse, error) {
	// Get the total number of items
	total, err := s.GetItemsCount(ctx)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error getting item count")
	}

	// Update the pagination with the total count
	pagination.Total = total

	if total == 0 {
		response := &types.PaginatedItemsResponse{
			Pagination: *pagination,
			Items:      []types.ItemsResponse{},
		}
		return response, nil
	}

	// Fetch the items for the current page
	items, err := s.q.GetAllItems(ctx, queries.GetAllItemsParams{
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	})
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error getting all items")
	}

	itemsResponse := make([]types.ItemsResponse, len(items))

	for i, item := range items {
		itemsResponse[i] = types.ItemsResponse{
			UUID:  item.Uuid.String(),
			Title: item.Title,
		}
	}

	response := types.PaginatedItemsResponse{
		Pagination: *pagination,
		Items:      itemsResponse,
	}

	return &response, nil
}

// GetItemByUuid returns a single item from the database by UUID
func (s *ItemsService) GetItemByUuid(ctx context.Context, uuid string) (*queries.GetItemByUuidRow, error) {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return nil, err
	}

	item, err := s.q.GetItemByUuid(ctx, *pgUuid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error getting item by uuid")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewAPIError(http.StatusNotFound, "item not found")
	}

	return &item, nil
}

// AddItem adds an item to the database
func (s *ItemsService) AddItem(ctx context.Context, request types.AddItemRequest) (*queries.AddItemRow, error) {
	if request.ItemTitle == "" {
		return nil, types.NewAPIError(http.StatusBadRequest, "title is required")
	}

	item, err := s.q.AddItem(ctx, request.ItemTitle)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error adding item")
	}

	return &item, nil
}

// UpdateItem updates an item by UUID
func (s *ItemsService) UpdateItem(ctx context.Context, uuid, title string) (*queries.UpdateItemRow, error) {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return nil, err
	}

	if title == "" {
		currentItem, err := s.q.GetItemByUuid(ctx, *pgUuid)
		if err != nil {
			return nil, types.NewAPIError(http.StatusInternalServerError, "error getting item by uuid")
		}
		title = currentItem.Title
	}

	item, err := s.q.UpdateItem(ctx, queries.UpdateItemParams{
		ItemUuid:  *pgUuid,
		ItemTitle: title,
	})
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "error updating item")
	}

	return &item, nil
}

// DeleteItem deletes an item from the database
func (s *ItemsService) DeleteItem(ctx context.Context, uuid string) error {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return err
	}

	err = s.q.DeleteItem(ctx, *pgUuid)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "error deleting item")
	}

	return nil
}

// GetItemsCount returns the count of all items
func (s *ItemsService) GetItemsCount(ctx context.Context) (int64, error) {
	return s.q.GetItemsCount(ctx)
}
