package types

type ItemsResponse struct {
	UUID  string `json:"uuid" example:"00000000-0000-0000-0000-000000000000"`
	Title string `json:"title" example:"Item title"`
}

// PaginatedItemsResponse represents a response containing a list of items
//
//	@Description	a response containing a list of items and a pagination object
type PaginatedItemsResponse struct {
	Pagination Pagination      `json:"pagination"`
	Items      []ItemsResponse `json:"items"`
}

// AddItemRequest represents the request body for adding a new item
//
//	@Description	a request body for adding a new item
type AddItemRequest struct {
	ItemTitle string `json:"title" example:"Item title"`
}

// UpdateItemRequest represents the request body for updating an item
// @Description a request body for updating an item
type UpdateItemRequest struct {
	ItemTitle string `json:"title" example:"Item title"`
}

// ItemDeletedResponse represents a success message for an item deletion
//
//	@Description	A success message confirming the item was deleted
type ItemDeletedResponse struct {
	Message string `json:"success" example:"Item deleted: 77b62cff-0020-43d9-a90c-5d35bff89f7a"`
}
