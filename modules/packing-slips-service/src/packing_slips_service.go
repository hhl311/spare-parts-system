package main

import (
	"../../business-structures"
	"./consumers"
	"log"
)

func main() {
	ordersReceiver := consumers.ValidatedOrdersReceiver{
		ChannelName:    "validated_orders",
		BusLocation:    "localhost:5672",
		BusCredentials: "guest:guest"}

	ordersReceiver.LaunchAcknowledgment(func(order models.Order) {
		if packingSlip, err := createPackingSlip(order); err == nil {
			err := notify(packingSlip)

			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func createPackingSlip(order models.Order) (models.PackingSlip, error) {
	// TODO Implement this.
	// TODO Fetch all article and replace it by its content if necessary
	return models.PackingSlip{}, nil
}

func notify(packingSlip models.PackingSlip) error {
	log.Println(packingSlip)

	return nil
}
