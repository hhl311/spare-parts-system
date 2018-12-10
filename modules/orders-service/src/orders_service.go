package main

import (
	"../../business-structures"
	"../../communication"
	"./consumers"
	"./dao"
	"fmt"
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
)

var router *gin.Engine
var ordersDao dao.OrdersDao
var sparePartsConsumer communication.SparePartsConsumer
var ordersSender consumers.ValidatedOrdersSender

func main() {
	// TODO Replace with a DB DAO.
	ordersDao = &dao.MapOrdersDao{}
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

	if validated, err := strconv.ParseBool(c.Query("validate")); err != nil || !validated {
		c.JSON(http.StatusNotModified, "Nothing is queried.")
		return
	}

	modified, err := ordersDao.Validate(orderId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error while validating the order")
		return
	}

	if modified {
		updated, err := ordersDao.GetOne(orderId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, "Error while getting the validated order")
			return
		}

		err = ordersSender.Send(updated)

		if err != nil {
			c.JSON(http.StatusInternalServerError, "Error while sending the validated order on the bus")
			return
		}

		fmt.Println("Validated order", orderId)
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
