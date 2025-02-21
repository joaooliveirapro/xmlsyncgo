package models

import "github.com/joaooliveirapro/xmlsyncgo/initializers"

type PaginatedResponse[T any] struct {
	Page       int   `json:"page"`
	TotalPages int64 `json:"totalPages"`
	TotalItems int64 `json:"totalItems"`
	Items      []T   `json:"items"`
}

func Paginate[T any](pageSize int, pageNumber int, whereQ string, orderQ string, whereA ...interface{}) (*PaginatedResponse[T], error) {
	var list []T
	var totalItems int64
	result := initializers.DB.Model(new(T)).Where(whereQ, whereA...).Count(&totalItems)
	if result.Error != nil {
		return nil, result.Error
	}
	totalPages := (totalItems + int64(pageSize) - 1) / int64(pageSize)
	if pageNumber > int(totalPages) {
		pageNumber = int(totalPages)
	}
	offset := (pageNumber - 1) * pageSize
	result = initializers.DB.Model(new(T)).Where(whereQ, whereA...).Order(orderQ).Offset(offset).Limit(pageSize).Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	response := PaginatedResponse[T]{
		Page:       pageNumber,
		TotalPages: totalPages,
		TotalItems: totalItems,
		Items:      list,
	}

	return &response, nil
}
