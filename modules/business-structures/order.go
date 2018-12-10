package models

import "time"

type Order struct {
	ID                int       `json:"id"`
	CustomerId        string    `json:"customerId" binding:"required"`
	ContentReferences []string  `json:"contentReferences" binding:"required"`
	CreationDate      time.Time `json:"creationDate"`
	Validated         bool      `json:"validated"`
}
