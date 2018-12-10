package dao

import (
	"../../../business-structures"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
	"log"
	"os"
)

// TODO Implement all methods.

const (
	DATABASE_LOCATION = "DATABASE_LOCATION"
)

var Logger = log.New(os.Stdout, "[DGraph spare parts DAO] ", log.Ldate|log.Ltime|log.Lshortfile)

type DGraphSpareParsDao struct {
	values         []models.SparePart
	databaseClient *dgo.Dgraph
}

func (dao *DGraphSpareParsDao) Create(article models.SparePart) error {
	dao.values = append(dao.values, article)

	return nil
}

func (dao *DGraphSpareParsDao) GetAll() ([]models.SparePart, error) {
	return dao.values, nil
}

func (dao *DGraphSpareParsDao) GetByReference(reference string) (models.SparePart, error) {
	for _, sparePart := range dao.values {
		if sparePart.Reference == reference {
			return sparePart, nil
		}
	}

	return models.SparePart{}, nil
}

func (dao *DGraphSpareParsDao) databaseReady() bool {
	if connection, err := grpc.Dial("DATABASE_LOCATION", grpc.WithInsecure()); err != nil {
		dao.databaseClient = dgo.NewDgraphClient(api.NewDgraphClient(connection))

		initialization := &api.Operation{}
		initialization.Schema= `
			reference: string @index(exact) . 
			name: string 
			contentReferences: string[]
			price: float
		`

		return true
	} else {
		return false
	}
}
