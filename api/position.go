package api

import (
	"hr-app-back/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PositionRead(c *gin.Context) {

	positions, err := storage.PositionRead(map[string]string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive positions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": positions})
}
