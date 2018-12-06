package dao

import "../../../business-structures"

// TODO Implement all methods.

type DGraphSpareParsDao struct {
	values []models.SparePart
}

func (dao *DGraphSpareParsDao) Create(article models.SparePart) error {
	dao.values = append(dao.values, article)

	return nil
}

func (dao *DGraphSpareParsDao) GetAll() ([]models.SparePart, error) {
	return dao.values, nil
}

func (dao *DGraphSpareParsDao) GetByReference(reference string) (models.SparePart, error) {
	for _, sparePart := range dao.values {
		if sparePart.Reference == reference {
			return sparePart, nil
		}
	}

	return models.SparePart{}, nil
}
