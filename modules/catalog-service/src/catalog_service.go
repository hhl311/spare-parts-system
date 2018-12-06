package main

import (
	"../../business-structures"
	"./dao"
	"github.com/gin-gonic/gin"
	"net/http"
)

var router *gin.Engine
var sparePartsDao dao.SparePartsDao = &dao.ArraySparePartsDao{}

func main() {
	router = gin.Default()

	router.POST("spare-parts/", createSparePart)
	router.GET("spare-parts/", getAllSpareParts)
	router.GET("spare-parts/:reference", getOneSparePart)

	_ = router.Run()
}

func createSparePart(c *gin.Context) {
	received := models.SparePart{ContentReferences: []string{}}

	if bindingErr := c.Bind(&received); bindingErr != nil {
		return
	}

	if daoErr := sparePartsDao.Create(received); daoErr == nil {
		c.JSON(http.StatusCreated, received)
	} else {
		c.JSON(http.StatusInternalServerError, daoErr)
	}
}

func getOneSparePart(c *gin.Context) {
	reference := c.Param("reference")

	found, err := sparePartsDao.GetByReference(reference)

	if err == nil {
		if found.Reference == reference {
			c.JSON(http.StatusOK, found)
		} else {
			c.JSON(http.StatusNoContent, nil)
		}
	} else {
		c.JSON(http.StatusInternalServerError, err)
	}
}

func getAllSpareParts(c *gin.Context) {
	found, err := sparePartsDao.GetAll()

	if err == nil {
		c.JSON(http.StatusOK, found)
	} else {
		c.JSON(http.StatusInternalServerError, err)
	}
}
