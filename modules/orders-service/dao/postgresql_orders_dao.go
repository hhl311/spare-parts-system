package dao

import (
	"../../../business-structures"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"log"
)

type PostgreSQLOrdersDao struct {
	DatabaseLocation string
	DatabaseUser     string
	DatabasePassword string
	initialized      bool
}

func (dao *PostgreSQLOrdersDao) Create(order models.Order) (int, error) {
	client, err := dao.makeClient()

	if err != nil {
		log.Println("Failed to create the client")
		return 0, err
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Println("Failed to close PostgreSQL client")
		}
	}()

	log.Println("Going to insert order", order)

	orderPointer := &order

	if err := client.Insert(orderPointer); err != nil {
		log.Println("Failed to insert a new order")

		return 0, err
	}

	return orderPointer.ID, nil
}

func (dao *PostgreSQLOrdersDao) Validate(orderId int) (bool, error) {
	client, err := dao.makeClient()

	if err != nil {
		log.Println("Failed to create the PostgreSQL client", err)
		return false, err
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Println("Failed to close PostgreSQL client")
		}
	}()

	if order, err := dao.GetOne(orderId); err != nil {
		return false, err
	} else if order.Validated == true {
		return false, nil
	} else {
		order.Validated = true

		if err = client.Update(&order); err != nil {
			return false, nil
		}
	}

	return true, nil
}

func (dao *PostgreSQLOrdersDao) GetAll() ([]models.Order, error) {
	client, err := dao.makeClient()

	if err != nil {
		log.Println("Failed to create the PostgreSQL client", err)
		return nil, err
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Println("Failed to close PostgreSQL client")
		}
	}()

	var values []models.Order
	if err := client.Model(&values).Select(); err != nil {
		log.Println("Failed to retrieve all orders")

		return nil, err
	}

	return values, nil
}

func (dao *PostgreSQLOrdersDao) GetOne(orderId int) (models.Order, error) {
	client, err := dao.makeClient()

	if err != nil {
		log.Println("Failed to create the PostgreSQL client", err)
		return models.Order{}, err
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Println("Failed to close PostgreSQL client")
		}
	}()

	var value = new(models.Order)

	err = client.Model(value).
		Where("ID = ?", orderId).
		Select()

	if err != nil {
		log.Println("Failed to retrieve all orders")

		return models.Order{}, err
	}

	return *value, nil
}

func (dao *PostgreSQLOrdersDao) makeClient() (*pg.DB, error) {
	client := pg.Connect(&pg.Options{
		Addr:     dao.DatabaseLocation,
		User:     dao.DatabaseUser,
		Password: dao.DatabasePassword,
	})

	for _, model := range []interface{}{&models.Order{}} {
		err := client.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
			Temp:        false,
		})

		if err != nil {
			log.Println("Failed to create the schema", err)
			return nil, err
		}
	}

	return client, nil
}
