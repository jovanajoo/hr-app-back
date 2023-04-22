package api

import (
	"encoding/base64"
	"hr-app-back/model"
	"hr-app-back/storage"
	"hr-app-back/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func EmployeeLogin(c *gin.Context) {

	var employeeFromDB model.Employee

	if err := c.ShouldBindJSON(&employeeFromDB); err != nil {
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
		return
	}

	employee, err := authEmployee(employeeFromDB.Email, employeeFromDB.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authToken := base64.StdEncoding.EncodeToString([]byte(employee.Email + ":" + employee.Password))

	c.JSON(http.StatusOK, gin.H{"token": authToken})
}

func EmployeeReadByOrg(c *gin.Context) {

	empByOrg, err := getEmpAndOrgByContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": empByOrg})
}

func EmployeeInsert(c *gin.Context) {

	// get organization ID from the context
	orgID, ok := c.Get("organizationID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization ID not found in context"})
		return
	}

	if !isAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user is not an admin"})
		return
	}

	var employee model.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	employee.OrganizationId = orgID.(int)
	var newPassword = utility.RandomPassword(8)
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	employee.Password = string(hash)

	err = storage.EmployeeInsert(&employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	sendPasswordEmail(&employee, newPassword)

	c.JSON(201, gin.H{"status": "Created"})
}

func EmployeeUpdate(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	orgID, ok := c.Get("organizationID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization ID not found in context"})
		return
	}

	employee, err := getEmployeeIdByContext(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user is not an admin"})
		return
	}

	if employee.OrganizationId != orgID {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "you are not authorized to modify this resource"})
		return
	}

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	employee.Id = id

	err = storage.EmployeeUpdate(employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func EmployeeDelete(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	orgID, ok := c.Get("organizationID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization ID not found in context"})
		return
	}

	employee, err := getEmployeeIdByContext(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user is not an admin"})
		return
	}

	if employee.OrganizationId != orgID {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "you are not authorized to modify this resource"})
		return
	}

	err = storage.EmployeeDelete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.Status(204)
}
