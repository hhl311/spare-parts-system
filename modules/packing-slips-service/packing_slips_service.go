package packing_slips_service

import (
	"log"
	"os"
	"spare-parts-system/modules/business-structures"
	"spare-parts-system/modules/communication"
	"spare-parts-system/modules/packing-slips-service/consumers"
	"spare-parts-system/modules/packing-slips-service/notifications"
	"time"
)

const (
	sparePartsServiceLocation = "SPARE_PARTS_SERVICE_LOCATION"

	validatedOrdersChannelName    = "VALIDATED_ORDERS_CHANNEL"
	validatedOrdersBusLocation    = "VALIDATED_ORDERS_BUS_LOCATION"
	validatedOrdersBusCredentials = "VALIDATED_ORDERS_BUS_CREDENTIALS"
)

var sparePartsConsumer communication.SparePartsConsumer
var notifier notifications.PackingSkipsNotifier

func main() {
	sparePartsConsumer = communication.SparePartsConsumer{
		ServiceLocation: os.Getenv(sparePartsServiceLocation)}

	ordersReceiver := consumers.ValidatedOrdersReceiver{
		ChannelName:    os.Getenv(validatedOrdersChannelName),
		BusLocation:    os.Getenv(validatedOrdersBusLocation),
		BusCredentials: os.Getenv(validatedOrdersBusCredentials)}

	notifier = &notifications.LoggerPackingSlipsNotifier{}

	ordersReceiver.LaunchAcknowledgment(func(order models.Order) {
		if packingSlip, err := createPackingSlip(order); err == nil {
			if err := notifier.Notify(packingSlip); err != nil {
				log.Println("Error while notifying the new packing slip", err)
			}
		} else {
			log.Println("Error while getting a validated order", err)
		}
	})
}

func createPackingSlip(order models.Order) (models.PackingSlip, error) {
	toBeSent := models.PackingSlip{OrderID: order.ID}

	for _, reference := range order.ContentReferences {
		log.Println("Getting spare parts associated to", reference)
		if sparePart, err := sparePartsConsumer.GetSparePart(reference); err != nil {
			log.Println("Failed to get spare parts associated to", reference)
			return models.PackingSlip{}, err
		} else {
			log.Println("Got spare part associated to", reference, ":", sparePart)

			if sparePart.ContentReferences == nil || len(sparePart.ContentReferences) == 0 {
				log.Println("Spare part", reference, "is not composed of others spare parts")
				toBeSent.ContentReferences = append(toBeSent.ContentReferences, sparePart.Reference)
			} else {
				log.Println("Spare part", reference, "is composed of", sparePart.ContentReferences)
				toBeSent.ContentReferences = append(toBeSent.ContentReferences, sparePart.ContentReferences...)
			}
		}
	}

	toBeSent.SentDate = time.Now()

	return toBeSent, nil
}
