package model

type Pagination struct {
	Limit  int
	Offset int
	Order  string
	Sort   string
}

type PaginatedResult struct {
	Total  int         `json:"total"`
	Result interface{} `json:"result"`
}
