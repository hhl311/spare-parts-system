package models

import "time"

type PackingSlip struct {
	OrderID           string    `json:"orderId"`
	ContentReferences []string  `json:"contentReferences" binding:"required"`
	SentDate          time.Time `json:"sentDate"`
}
