package api

import (
	"hr-app-back/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DepartmentsRead(c *gin.Context) {

	departments, err := storage.DepartmentRead(map[string]string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive departments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": departments})

}
