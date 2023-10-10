package structz

type GridRequest struct {
	Code         string                 `json:"code"`
	StartRow     int                    `json:"startRow"`
	EndRow       int                    `json:"endRow"`
	FilterModel  map[string]FilterModel `json:"filterModel"`
	SortModel    []SortModel            `json:"sortModel"`
	DefaultParam []DefaultParam         `json:"defaultParam"`
}

type FilterModel struct {
	Filter     string `json:"filter"`
	FilterType string `json:"filterType"` //TYPE: text number date set
	Type       string `json:"type"`
}

type SortModel struct {
	ColId string `json:"colId"`
	Sort  string `json:"sort"`
}

type DefaultParam struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
