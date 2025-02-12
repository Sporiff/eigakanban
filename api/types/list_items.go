package types

type ListItemsResponse struct {
	UUID     string `json:"uuid" example:"00000000-0000-0000-0000-000000000000"`
	ListUUID string `json:"list_uuid" example:"00000000-0000-0000-0000-000000000001"`
	ItemUUID string `json:"item_uuid" example:"00000000-0000-0000-0000-000000000002"`
	Status   string `json:"status" example:"Backlog"`
	Position int32  `json:"position" example:"0"`
}

type PaginatedListItemsResponse struct {
	Pagination Pagination          `json:"pagination"`
	ListItems  []ListItemsResponse `json:"list_items"`
}
