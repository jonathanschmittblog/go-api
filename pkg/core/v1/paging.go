package v1

import "math"

type Paging struct {
	ActualPage   int         `json:"page,omitempty"`
	PageCount    int         `json:"pageCount,omitempty"`
	TotalRecords int         `json:"totalRecords,omitempty"`
	Items        interface{} `json:"items,omitempty"`
}

func NewPaging(actualPage int, totalRecords int, pageSize int, items interface{}) Paging {
	return Paging{
		ActualPage:   actualPage,
		TotalRecords: totalRecords,
		Items:        items,
		PageCount:    pageCount(totalRecords, pageSize),
	}
}

func pageCount(totalRecords int, pageSize int) int {
	result := int(math.Floor(float64(totalRecords) / float64(pageSize)))
	if math.Mod(float64(totalRecords), float64(pageSize)) > 0 {
		result++
	}
	return result
}
