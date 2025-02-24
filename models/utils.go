package models

import (
	"time"

	"github.com/joaooliveirapro/xmlsyncgo/initializers"
	"gorm.io/gorm"
)

type CommonFields struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type PaginatedResponse[T any] struct {
	Page       int   `json:"page"`
	TotalPages int64 `json:"totalPages"`
	Total      int64 `json:"total"`
	Data       []T   `json:"data"`
}

type PaginateArgs struct {
	PageSize   int
	PageNumber int
	WhereQ     string
	OrderQ     string
	Preload    bool
	PreloadQ   string
	WhereA     []interface{}
}

func Paginate[T any](args PaginateArgs) (*PaginatedResponse[T], error) {
	var list []T
	var totalItems int64
	result := initializers.DB.Model(new(T)).Where(args.WhereQ, args.WhereA...).Count(&totalItems)
	if result.Error != nil {
		return nil, result.Error
	}
	totalPages := (totalItems + int64(args.PageSize) - 1) / int64(args.PageSize)
	if args.PageNumber > int(totalPages) {
		args.PageNumber = int(totalPages)
	}
	offset := (args.PageNumber - 1) * args.PageSize
	if args.Preload {
		result = initializers.DB.Model(new(T)).Where(args.WhereQ, args.WhereA...).Order(args.OrderQ).Offset(offset).Limit(args.PageSize).Find(&list)
	} else {
		result = initializers.DB.Model(new(T)).Where(args.WhereQ, args.WhereA...).Order(args.OrderQ).Offset(offset).Limit(args.PageSize).Preload(args.PreloadQ).Find(&list)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	response := PaginatedResponse[T]{
		Page:       args.PageNumber,
		TotalPages: totalPages,
		Total:      totalItems,
		Data:       list,
	}
	return &response, nil
}
