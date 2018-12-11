package dao

import "spare-parts-system/modules/business-structures"

type MapOrdersDao struct {
	content map[int]models.Order
}

func (dao *MapOrdersDao) Create(order models.Order) (int, error) {
	if dao.content == nil {
		dao.content = make(map[int]models.Order)
	}

	order.ID = len(dao.content)

	dao.content[order.ID] = order

	return order.ID, nil
}

func (dao *MapOrdersDao) Validate(orderId int) (bool, error) {
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

func (dao *MapOrdersDao) GetOne(orderId int) (models.Order, error) {
	return dao.content[orderId], nil
}
