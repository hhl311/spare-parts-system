package models

type SparePart struct {
	Reference         string   `json:"reference" binding:"required"`
	Name              string   `json:"name" binding:"required"`
	ContentReferences []string `json:"contentReferences"`
	Price             float64  `json:"price" binding:"required"`
}
