package dao

import (
	"context"
	"encoding/json"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
	"os"
	"spare-parts-system/modules/business-structures"
)

var Logger = log.New(os.Stdout, "[DGraph spare parts DAO] ", log.Ldate|log.Ltime|log.Lshortfile)

type DGraphSparePartsDao struct {
	DatabaseLocation string
	databaseClient   *dgo.Dgraph
}

func (dao *DGraphSparePartsDao) Create(article models.SparePart) error {
	log.Println("Saving", article, "in the database")
	if err := dao.databaseReady(); err != nil {
		log.Println("Error while preparing the database: ", err)
		return err
	}

	if jsonArticle, err := json.Marshal(prepareArticle(article)); err == nil {
		mutation := &api.Mutation{CommitNow: true}
		mutation.SetJson = jsonArticle

		assigned, err := dao.databaseClient.NewTxn().Mutate(context.Background(), mutation)

		log.Println("Assigned:", assigned)
		return err
	} else {
		log.Println("Failure when marshalling article")
		return err
	}
}

type GetAllResponse struct {
	Articles []TransitionSparePart `json:"Articles"`
}

func (dao *DGraphSparePartsDao) GetAll() ([]models.SparePart, error) {
	log.Println("Getting all Articles from the database")
	if err := dao.databaseReady(); err != nil {
		log.Println("Error while preparing the database: ", err)
		return nil, err
	}

	response, err := dao.databaseClient.NewTxn().Query(context.Background(), `{
  		Articles (func: has (_predicate_) ) {
			reference
    		name
			price
			contentReferences {
				reference
			}
  		}
	}`)

	if err != nil {
		return nil, err
	}

	var unmarshalled GetAllResponse

	if err := json.Unmarshal(response.Json, &unmarshalled); err != nil {
		log.Println("Error while unmarshalling database response: ", string(response.Json))
		return nil, err
	} else {
		var retrieved []models.SparePart
		for _, r := range unmarshalled.Articles {
			retrieved = append(retrieved, retrieveArticle(r))
		}

		return retrieved, nil
	}
}

func (dao *DGraphSparePartsDao) GetByReference(reference string) (models.SparePart, error) {
	log.Println("Getting article with reference", reference, "from the database")
	if err := dao.databaseReady(); err != nil {
		log.Println("Error while preparing the database: ", err)
		return models.SparePart{}, err
	}

	response, err := dao.databaseClient.NewTxn().Query(context.Background(), `{
  		Articles (func: eq(reference, "`+reference+`")) {
			reference
    		name
			price
			contentReferences {
				reference
			}
		}
	}`)

	if err != nil {
		return models.SparePart{}, err
	}

	var unmarshalled GetAllResponse

	if err := json.Unmarshal(response.Json, &unmarshalled); err != nil {
		log.Println("Error while unmarshalling database response: ", string(response.Json))
		return models.SparePart{}, err
	} else if len(unmarshalled.Articles) == 0 {
		return models.SparePart{}, errors.New("NOT FOUND")
	} else {
		return retrieveArticle(unmarshalled.Articles[0]), nil
	}
}

func (dao *DGraphSparePartsDao) databaseReady() error {
	if dao.databaseClient != nil {
		return nil
	}

	if connection, err := grpc.Dial(dao.DatabaseLocation, grpc.WithInsecure()); err == nil {
		log.Println("Connected to database")
		dao.databaseClient = dgo.NewDgraphClient(api.NewDgraphClient(connection))

		return dao.setupDatabase()
	} else {
		dao.databaseClient = nil
		log.Println("Could not connect to database: ", err)
		return err
	}
}

func (dao *DGraphSparePartsDao) setupDatabase() error {
	initialization := &api.Operation{}
	initialization.Schema = `
			reference: string @index(exact) . 
			name: string .
			price: float .
	`

	log.Println("Going to alter database schema")
	return dao.databaseClient.Alter(context.Background(), initialization)
}

type TransitionReference struct {
	Reference string `json:"reference"`
}

type TransitionSparePart struct {
	Reference         string                `json:"reference"`
	Name              string                `json:"name"`
	Price             float64               `json:"price"`
	ContentReferences []TransitionReference `json:"contentReferences"`
}

func prepareArticle(article models.SparePart) TransitionSparePart {
	var references []TransitionReference
	for _, r := range article.ContentReferences {
		references = append(references, TransitionReference{Reference: r})
	}

	return TransitionSparePart{Reference: article.Reference, Name: article.Name, Price: article.Price, ContentReferences: references}
}

func retrieveArticle(transition TransitionSparePart) models.SparePart {
	var references []string
	for _, r := range transition.ContentReferences {
		references = append(references, r.Reference)
	}

	return models.SparePart{Reference: transition.Reference, Name: transition.Name, Price: transition.Price, ContentReferences: references}
}
