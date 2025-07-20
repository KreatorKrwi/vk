query запрос на отображение ленты:
type AdsFilters struct {
	Page     int    `form:"page" binding:"min=1"`
	PerPage  int    `form:"per_page" binding:"min=1,max=100"`
	SortBy   string `form:"sort_by,omitempty" binding:"omitempty,oneof=date_desc date_asc price_desc price_asc"`
	MinPrice *int   `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice *int   `form:"max_price" binding:"omitempty,min=1"`
}
