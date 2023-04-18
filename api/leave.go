package api

import (
	"fmt"
	"hr-app-back/model"
	"hr-app-back/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func LeaveCreate(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	employee, orgID, err := getEmployeeByContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// ne znam da l mi je potreban bas ovaj deo orgID jer ja svakako trazim employee na osnovu orgID i mejla
	if employee.OrganizationId != orgID || employee.Id != id {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "you are not authorized to modify this resource"})
		return
	}

	var leave model.Leave
	if err := c.ShouldBindJSON(&leave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	leave.EmployeeId = id
	leave.OrganizationId = employee.OrganizationId
	leave.Total = int(leave.EndDate.Sub(leave.StartDate).Hours() / 24)
	leave.Status = "pending"

	err = storage.LeaveCreate(&leave)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "Created", "message": "Leave request is submitted"})
}

func LeavesStatusRead(c *gin.Context) {

	employee, orgID, err := getEmployeeByContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status := c.Query("status")

	var leaveRequests []model.Leave

	filter := map[string]string{"organizationID": fmt.Sprintf("%d", orgID)}
	if status != "" || !isAdmin(c) {
		filter["status"] = status
	}
	if !isAdmin(c) {
		filter["employeeID"] = fmt.Sprintf("%d", employee.Id)
	}

	leaveRequests, err = storage.LeaveRead(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falied to retrvie leave requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": leaveRequests})

}

func LeaveUpdate(c *gin.Context) {
	leaveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	employee, _, err := getEmployeeByContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	leave, err := getLeaveRequest(c, leaveID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if leave.EmployeeId != employee.Id || leave.OrganizationId != employee.OrganizationId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update this leave request"})
		return
	}

	if err := c.ShouldBindJSON(&leave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	leave.EmployeeId = employee.Id
	leave.OrganizationId = employee.OrganizationId
	leave.Total = int(leave.EndDate.Sub(leave.StartDate).Hours() / 24)
	leave.Status = "pending"

	err = storage.LeaveUpdate(leave)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Leave request updated successfully"})

}

func LeaveStatusAdminUpdate(c *gin.Context) {
	leaveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	employee, _, err := getEmployeeByContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update status of the leave request"})
		return
	}

	leaveRequest, err := getLeaveRequest(c, leaveID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if leaveRequest.OrganizationId != employee.OrganizationId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update status of leave request from another organziation"})
		return
	}

	var updateLeave model.Leave
	if err := c.ShouldBindJSON(&updateLeave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if updateLeave.Status != "approved" && updateLeave.Status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid status"})
		return
	}

	leaveRequest.Status = updateLeave.Status

	err = storage.LeaveUpdateStatus(leaveRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update leave request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Leave request updated successfully"})

}

func LeaveReadAll(c *gin.Context) {

	orgID, ok := c.Get("organizationID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization ID not found in context"})
		return
	}

	tmp := map[string]string{"organizationID": fmt.Sprintf("%d", orgID)}
	leaves, err := storage.LeaveRead(tmp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employees", "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "data": leaves})

}
