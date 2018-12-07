package dao

import (
	"../../../business-structures"
	"strconv"
)

type MapOrdersDao struct {
	content map[string]models.Order
}

func (dao *MapOrdersDao) Create(order models.Order) (string, error) {
	if dao.content == nil {
		dao.content = make(map[string]models.Order)
	}

	order.ID = strconv.Itoa(len(dao.content))

	dao.content[order.ID] = order

	return order.ID, nil
}

func (dao *MapOrdersDao) Validate(orderId string) (bool, error) {
	order, isKnown := dao.content[orderId]

	if isKnown {
		order.Validated = true
		dao.content[orderId] = order

		return true, nil
	} else {
		return false, nil
	}
}

func (dao *MapOrdersDao) GetAll() ([]models.Order, error) {
	var values []models.Order
	for _, value := range dao.content {
		values = append(values, value)
	}

	return values, nil
}

func (dao *MapOrdersDao) GetOne(orderId string) (models.Order, error) {
	return dao.content[orderId], nil
}
