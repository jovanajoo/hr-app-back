package api

import (
	"encoding/base64"
	"fmt"
	"hr-app-back/model"
	"hr-app-back/storage"
	"hr-app-back/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func EmployeeLogin(c *gin.Context) {
	var employee model.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Get the employee from the database by email

	tmp := map[string]string{"email": employee.Email, "password": employee.Password}
	employees, err := storage.EmployeeRead(tmp)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email of password"})
		return
	}

	if len(employees) == 0 {
		//todo status code unathorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authToken := base64.StdEncoding.EncodeToString([]byte(employees[0].Email + ":" + employees[0].Password))

	c.JSON(http.StatusOK, gin.H{"token": authToken})

}

func EmployeeGet(c *gin.Context) {

	//todo get orgid from context
	orgID, ok := c.Get("organizationID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization ID not found in context"})
		return
	}
	//todo set as filter into map to employee read

	tmp := map[string]string{"organizationID": fmt.Sprintf("%d", orgID)}
	employees, err := storage.EmployeeRead(tmp)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve employees"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "data": employees})

}

func EmployeeInsert(c *gin.Context) {

	// get organization ID from the context
	orgID, ok := c.Get("organizationID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization ID not found in context"})
		return
	}
	// todo get from context email
	email, ok := c.Get("email")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Email not found in context"})
		return
	}

	var employee model.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
		return
	}

	tmp := map[string]string{"email": email.(string), "isAdmin": strconv.FormatBool(employee.IsAdmin)}
	//employye read to check if that email is admin for that organization
	admin, err := storage.EmployeeRead(tmp)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}

	//if not return unauthorized
	if len(admin) == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user is not an admin"})
		return
	}

	employee.OrganizationId = orgID.(int)
	employee.Password = utility.RandomPassword(8)

	err = storage.EmployeeInsert(&employee)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"status": "Created"})
}

func EmployeeUpdate(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	orgID, ok := c.Get("organizationID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Organization ID not found in context"})
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Email not found in context"})
		return
	}

	tmp := map[string]string{"organizationID": fmt.Sprintf("%d", orgID), "email": email.(string), "isAdmin": strconv.FormatBool(true)}
	_, err = storage.EmployeeRead(tmp)
	if err != nil {
		c.JSON(500, gin.H{"error": "Employee is not admin"})
		return
	}

	tmp2 := map[string]string{"employeeID": fmt.Sprintf("%d", id), "fields": "organizationID"}
	employee, err := storage.EmployeeRead(tmp2)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve employee"})
		return
	}

	if employee[0].OrganizationId != orgID {
		c.JSON(500, gin.H{"status": "error", "message": "you are not authorized to modify this resource"})
		return
	}

	if err := c.ShouldBindJSON(&employee[0]); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	employee[0].Id = id

	err = storage.EmployeeUpdate(&employee[0])
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}

func EmployeeDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	orgID, ok := c.Get("organizationID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Organization ID not found in context"})
		return
	}

	employee, err := storage.EmployeeGetByOrgId(orgID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if employee.OrganizationId != orgID {
		c.JSON(500, gin.H{"status": "error", "message": "you are not authorized to access this resource"})
		return
	}

	err = storage.EmployeeDelete(id)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}
