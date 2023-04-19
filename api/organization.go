package api

import (
	"fmt"
	"hr-app-back/model"
	"hr-app-back/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func OrganizationReadById(c *gin.Context) {

	orgID, ok := c.Get("organizationID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization ID not found in context"})
		return
	}

	if !isAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user is not an admin"})
		return
	}

	organization, err := storage.OrganizationRead(map[string]string{"organizationID": fmt.Sprintf("%d", orgID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employees"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "data": organization})

}

func OrganizationRegister(c *gin.Context) {

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

func OrganizationUpdate(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	orgID, ok := c.Get("organizationID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization ID not found in context"})
		return
	}

	if !isAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user is not an admin"})
		return
	}

	if orgID != id {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized: user is not an admin"})
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

	if !isAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user is not an admin"})
		return
	}

	err = storage.OrganizationDelete(id)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}
