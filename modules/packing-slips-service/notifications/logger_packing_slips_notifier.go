package notifications

import (
	"log"
	"os"
	"spare-parts-system/modules/business-structures"
)

var Logger = log.New(os.Stdout, "[Packing slips notification] ", log.Ldate|log.Ltime|log.Lshortfile)

type LoggerPackingSlipsNotifier struct {
}

func (logger *LoggerPackingSlipsNotifier) Notify(packingSlip models.PackingSlip) error {
	Logger.Println("# New packing slip:")
	Logger.Println("# - Associated order:", packingSlip.OrderID)
	Logger.Println("# - Sent date:", packingSlip.SentDate)
	Logger.Println("# - References of content:", packingSlip.ContentReferences)

	return nil
}
