package main

import (
	"../../business-structures"
	"./dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var router *gin.Engine
var ordersDao dao.OrdersDao = &dao.MapOrdersDao{}

func main() {
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

	// TODO Communicate with Catalog service to make sure spare parts references exist.

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
