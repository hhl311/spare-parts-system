package notifications

import "github.com/AntoineAube/spare-parts-system/modules/business-structures"

type PackingSkipsNotifier interface {
	Notify(packingSlip models.PackingSlip) error
}
