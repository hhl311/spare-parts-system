package main

import (
	"fmt"
	"github.com/AntoineAube/spare-parts-system/modules/business-structures"
	"github.com/AntoineAube/spare-parts-system/modules/communication"
	"github.com/AntoineAube/spare-parts-system/modules/orders-service/consumers"
	"github.com/AntoineAube/spare-parts-system/modules/orders-service/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	sparePartsServiceLocation = "SPARE_PARTS_SERVICE_LOCATION"

	validatedOrdersChannelName    = "VALIDATED_ORDERS_CHANNEL"
	validatedOrdersBusLocation    = "VALIDATED_ORDERS_BUS_LOCATION"
	validatedOrdersBusCredentials = "VALIDATED_ORDERS_BUS_CREDENTIALS"

	databaseLocation = "DATABASE_LOCATION"
	databaseUser     = "DATABASE_USER"
	databasePassword = "DATABASE_PWD"
)

var router *gin.Engine
var ordersDao dao.OrdersDao
var sparePartsConsumer communication.SparePartsConsumer
var ordersSender consumers.ValidatedOrdersSender

func main() {
	ordersDao = &dao.PostgreSQLOrdersDao{DatabaseLocation: os.Getenv(databaseLocation),
		DatabaseUser:     os.Getenv(databaseUser),
		DatabasePassword: os.Getenv(databasePassword)}
	sparePartsConsumer = communication.SparePartsConsumer{
		ServiceLocation: os.Getenv(sparePartsServiceLocation)}
	ordersSender = consumers.ValidatedOrdersSender{
		ChannelName:    os.Getenv(validatedOrdersChannelName),
		BusLocation:    os.Getenv(validatedOrdersBusLocation),
		BusCredentials: os.Getenv(validatedOrdersBusCredentials)}

	router = gin.Default()

	router.POST("orders/", createOrder)
	router.GET("orders/", getAllOrders)
	router.PUT("orders/:id", validateOneOrder)

	_ = router.Run()
}

func createOrder(c *gin.Context) {
	var received models.Order

	if bindingErr := c.Bind(&received); bindingErr != nil {
		return
	}

	received.CreationDate = time.Now()
	received.Validated = false

	for _, reference := range received.ContentReferences {
		if _, err := sparePartsConsumer.GetSparePart(reference); err != nil {
			c.JSON(http.StatusBadRequest, "Ordered reference(s) do not exist")
			return
		}
	}

	if id, err := ordersDao.Create(received); err == nil {
		received.ID = id

		c.JSON(http.StatusCreated, received)
	} else {
		c.JSON(http.StatusInternalServerError, "Error while creating an order")
	}
}

func validateOneOrder(c *gin.Context) {
	orderId := c.Param("id")

	orderIdInt, err := strconv.Atoi(orderId)

	if err != nil {
		c.JSON(http.StatusBadRequest, "Parameter ID should be an integer")
	}

	if validated, err := strconv.ParseBool(c.Query("validate")); err != nil || !validated {
		c.JSON(http.StatusNotModified, "Nothing is queried.")
		return
	}

	modified, err := ordersDao.Validate(orderIdInt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error while validating the order")
		return
	}

	if modified {
		updated, err := ordersDao.GetOne(orderIdInt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, "Error while getting the validated order")
			return
		}

		err = ordersSender.Send(updated)

		if err != nil {
			c.JSON(http.StatusInternalServerError, "Error while sending the validated order on the bus")
			return
		}

		fmt.Println("Validated order", orderIdInt)
		c.JSON(http.StatusOK, "Validated!")
	} else {
		c.JSON(http.StatusNotModified, "Not found or already validated!")
	}
}

func getAllOrders(c *gin.Context) {
	if found, err := ordersDao.GetAll(); err == nil {
		c.JSON(http.StatusOK, found)
	} else {
		c.JSON(http.StatusInternalServerError, "Error while getting all orders")
	}
}
