package dao

import "../../../business-structures"

type OrdersDao interface {
	Create(order models.Order) (string, error)
	Validate(orderId string) (bool, error)
	GetAll() ([]models.Order, error)
}
