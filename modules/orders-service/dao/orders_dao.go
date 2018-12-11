package dao

import "github.com/AntoineAube/spare-parts-system/modules/business-structures"

type OrdersDao interface {
	Create(order models.Order) (int, error)
	Validate(orderId int) (bool, error)
	GetOne(orderId int) (models.Order, error)
	GetAll() ([]models.Order, error)
}
