package services

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/helpers"
	"codeberg.org/sporiff/eigakanban/types"
	"context"
	"errors"
)

type ItemsService struct {
	q *queries.Queries
}

func NewItemsService(q *queries.Queries) *ItemsService {
	return &ItemsService{q: q}
}

// GetAllItems retrieves all items from the database as a paginated list
func (s *ItemsService) GetAllItems(ctx context.Context, pagination *types.Pagination) ([]queries.GetAllItemsRow, *types.Pagination, error) {
	// Get the total number of items
	total, err := s.GetItemsCount(ctx)
	if err != nil {
		return nil, pagination, errors.New("could not get item count: " + err.Error())
	}

	// Update the pagination with the total count
	pagination.Total = total

	// Fetch the items for the current page
	items, err := s.q.GetAllItems(ctx, queries.GetAllItemsParams{
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	})
	if err != nil {
		return nil, pagination, errors.New("error getting all items: " + err.Error())
	}

	return items, pagination, nil
}

// GetItemByUuid returns a single item from the database by UUID
func (s *ItemsService) GetItemByUuid(ctx context.Context, uuid string) (queries.GetItemByUuidRow, error) {
	emptyItem := queries.GetItemByUuidRow{}

	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return emptyItem, errors.New("error validating uuid: " + err.Error())
	}

	item, err := s.q.GetItemByUuid(ctx, pgUuid)
	if err != nil {
		return emptyItem, errors.New("error getting item by uuid: " + err.Error())
	}

	return item, nil
}

// AddItem adds an item to the database
func (s *ItemsService) AddItem(ctx context.Context, request types.AddItemRequest) (queries.AddItemRow, error) {
	emptyItem := queries.AddItemRow{}

	if request.ItemTitle == "" {
		return emptyItem, errors.New("title is required")
	}

	item, err := s.q.AddItem(ctx, request.ItemTitle)
	if err != nil {
		return emptyItem, errors.New("error adding item: " + err.Error())
	}

	return item, nil
}

// UpdateItem updates an item by UUID
func (s *ItemsService) UpdateItem(ctx context.Context, uuid, title string) (queries.UpdateItemRow, error) {
	emptyItem := queries.UpdateItemRow{}

	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return emptyItem, errors.New("error validating uuid: " + err.Error())
	}

	if title == "" {
		currentItem, err := s.q.GetItemByUuid(ctx, pgUuid)
		if err != nil {
			return emptyItem, errors.New("error getting item by uuid: " + err.Error())
		}
		title = currentItem.Title
	}

	item, err := s.q.UpdateItem(ctx, queries.UpdateItemParams{
		ItemUuid:  pgUuid,
		ItemTitle: title,
	})
	if err != nil {
		return emptyItem, errors.New("error updating item: " + err.Error())
	}
	return item, nil
}

// DeleteItem deletes an item from the database
func (s *ItemsService) DeleteItem(ctx context.Context, uuid string) error {
	pgUuid, err := helpers.ValidateAndConvertUUID(uuid)
	if err != nil {
		return errors.New("error validating uuid: " + err.Error())
	}

	err = s.q.DeleteItem(ctx, pgUuid)
	if err != nil {
		return errors.New("error deleting item: " + err.Error())
	}

	return nil
}

// GetItemsCount returns the count of all items
func (s *ItemsService) GetItemsCount(ctx context.Context) (int64, error) {
	return s.q.GetItemsCount(ctx)
}
