package schema

type ProductListRequest struct {
	// TODO: add filter
}

type ProductListResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Category    string `json:"category"`
}
