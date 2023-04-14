package api

import (
	"encoding/base64"
	"hr-app-back/storage"
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

		//todo base64decode
		//split :
		//email:password

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

		tmp := map[string]string{"email": email, "password": password}

		employeeFromDB, err := storage.EmployeeRead(tmp)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email of password"})
			return
		}
		if len(employeeFromDB) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": err.Error()})
			return
		}
		//todo select employee with that email and password from db
		//if len 0, unauthorized
		//if len 1, set organizaytionID

		//todo set also email to context
		organizationID := employeeFromDB[0].OrganizationId
		isAdmin := employeeFromDB[0].IsAdmin

		c.Set("organizationID", organizationID)
		c.Set("email", email)
		c.Set("isAdmin", isAdmin)
		c.Next()

	}
}
