package structz

import "math"

type GridResponse struct {
	Data                  interface{} `json:"data"`
	Meta                  Meta        `json:"meta"`
	LastRow               interface{} `json:"lastRow"`
	SecondaryColumnFields interface{} `json:"secondaryColumnFields"`
}

type Meta struct {
	Page    int         `json:"page"`
	Perpage int         `json:"perpage"`
	Total   int         `json:"total"`
	Sort    interface{} `json:"sort"`
	Field   interface{} `json:"field"`
}

func (e *GridResponse) FillPaginate(totalRows int, pageNumber int, pageSize int) {
	// offset := pageNumber * pageSize
	totalPages := int(math.Ceil(float64(totalRows) / float64(pageSize)))

	e.Meta = Meta{
		Page:    totalRows,
		Perpage: pageSize,
		Total:   totalPages,
		Sort:    "",
		Field:   "",
	}
}
