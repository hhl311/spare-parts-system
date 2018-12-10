package communication

import (
	"../business-structures"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var Logger = log.New(os.Stdout, "[Spare parts consumer] ", log.Ldate|log.Ltime|log.Lshortfile)

type SparePartsConsumer struct {
	ServiceLocation string
}

func (consumer *SparePartsConsumer) GetSparePart(reference string) (models.SparePart, error) {
	Logger.Println("Get spare part with reference", reference)

	response, getErr := http.Get("http://" + consumer.ServiceLocation + "/spare-parts/" + reference)
	if getErr != nil {
		Logger.Fatal("Get error:", getErr)
		return models.SparePart{}, getErr
	}

	if response.StatusCode == 204 {
		Logger.Println("No known spare part with reference ", reference)
		return models.SparePart{}, errors.New("SPARE PART NOT FOUND")
	}

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		Logger.Fatal("Read error:", readErr)
		return models.SparePart{}, getErr
	}

	var sparePart models.SparePart
	if jsonErr := json.Unmarshal(body, &sparePart); jsonErr != nil {
		Logger.Fatal("JSON error:", jsonErr)
		return models.SparePart{}, jsonErr
	}

	return sparePart, nil
}
