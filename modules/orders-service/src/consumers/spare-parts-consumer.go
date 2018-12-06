package consumers

import (
	"../../../business-structures"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type SparePartsConsumer struct {
	ServiceLocation string
}

func (consumer *SparePartsConsumer) GetSparePart(reference string) (models.SparePart, error) {
	response, getErr := http.Get("http://" + consumer.ServiceLocation + "/spare-parts/" + reference)
	if getErr != nil {
		log.Fatal("Get error: ", getErr)
		return models.SparePart{}, getErr
	}

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Fatal("Read error: ", readErr)
		return models.SparePart{}, getErr
	}

	var sparePart models.SparePart
	if jsonErr := json.Unmarshal(body, &sparePart); jsonErr != nil {
		log.Fatal("JSON error: ", jsonErr)
		return models.SparePart{}, jsonErr
	}

	return sparePart, nil
}
