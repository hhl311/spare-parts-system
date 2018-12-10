package notifications

import "../../../business-structures"

type PackingSkipsNotifier interface {
	Notify(packingSlip models.PackingSlip) error
}
