package main

import (
	"../../business-structures"
	"./consumers"
	"./dao"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

var router *gin.Engine
var ordersDao dao.OrdersDao = &dao.MapOrdersDao{}
var sparePartsConsumer consumers.SparePartsConsumer

func main() {
	sparePartsConsumer = consumers.SparePartsConsumer{ServiceLocation: os.Getenv("SPARE_PARTS_SERVICE_LOCATION")}

	router = gin.Default()

	router.POST("orders/", createOrder)
	router.GET("orders/", getAllOrders)
	router.PUT("orders/:reference", validateOneOrder)

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
		c.JSON(http.StatusInternalServerError, err)
	}
}

func validateOneOrder(c *gin.Context) {
	orderId := c.Param("reference")

	if validated, err := strconv.ParseBool(c.Query("validate")); err != nil || !validated {
		c.JSON(http.StatusNotModified, "Nothing is queried.")
		return
	}

	modified, err := ordersDao.Validate(orderId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if modified {
		// TODO Enqueue the validated order in the bus.
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
		c.JSON(http.StatusInternalServerError, err)
	}
}
