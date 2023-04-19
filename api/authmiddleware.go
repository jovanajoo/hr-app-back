package api

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func authMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "basic" { //todo basic
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(authHeaderParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid base64-encoded credentials"})
			return
		}

		emailPassword := strings.Split(string(decoded), ":")
		if len(emailPassword) != 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials format"})
			return
		}

		email := emailPassword[0]
		password := emailPassword[1]

		employeeFromDB, err := authEmployee(email, password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		organizationID := employeeFromDB.OrganizationId
		employeeID := employeeFromDB.Id
		isAdmin := employeeFromDB.IsAdmin

		c.Set("organizationID", organizationID)
		c.Set("employeeID", employeeID)
		c.Set("email", email)
		c.Set("isAdmin", isAdmin)
		c.Next()

	}
}
