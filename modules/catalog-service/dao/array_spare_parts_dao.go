package dao

import "github.com/AntoineAube/spare-parts-system/modules/business-structures"

type ArraySparePartsDao struct {
	values []models.SparePart
}

func (dao *ArraySparePartsDao) Create(article models.SparePart) error {
	dao.values = append(dao.values, article)

	return nil
}

func (dao *ArraySparePartsDao) GetAll() ([]models.SparePart, error) {
	return dao.values, nil
}

func (dao *ArraySparePartsDao) GetByReference(reference string) (models.SparePart, error) {
	for _, sparePart := range dao.values {
		if sparePart.Reference == reference {
			return sparePart, nil
		}
	}

	return models.SparePart{}, nil
}
