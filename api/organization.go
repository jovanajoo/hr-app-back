package api

import (
	"hr-app-back/model"
	"hr-app-back/storage"
	"strconv"

	"github.com/gin-gonic/gin"
)

func OrganizationInsert(c *gin.Context) {

	var organization model.Organization

	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := storage.OrganizationInsert(&organization)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})

}

func OrganizationGet(c *gin.Context) {
	var organization []model.Organization

	orgs, err := storage.OrganizationGet(organization)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": &orgs})

}

func OrganizationUpdate(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var organization model.Organization
	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	organization.Id = id

	err = storage.OrganizationUpdate(&organization)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})

}

func OrganizationDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	err = storage.OrganizationDelete(id)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})

}
