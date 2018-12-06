package models

import "time"

type Order struct {
	ID                string    `json:"reference"`
	CustomerId        string    `json:"customerId" binding:"required"`
	ContentReferences []string  `json:"contentReferences" binding:"required"`
	CreationDate      time.Time `json:"creationDate"`
	Validated         bool      `json:"validated"`
}
