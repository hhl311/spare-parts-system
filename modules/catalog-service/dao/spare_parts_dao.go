package dao

import "github.com/AntoineAube/spare-parts-system/modules/business-structures"

type SparePartsDao interface {
	Create(article models.SparePart) error
	GetAll() ([]models.SparePart, error)
	GetByReference(reference string) (models.SparePart, error)
}
