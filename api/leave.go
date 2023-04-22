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

	orgID := c.GetInt("organizationID")
	empID := c.GetInt("employeeID")

	var leave model.Leave
	if err := c.ShouldBindJSON(&leave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	leave.EmployeeId = empID
	leave.OrganizationId = orgID
	leave.Total = int(leave.EndDate.Sub(leave.StartDate).Hours() / 24)
	leave.Status = "pending"

	id, err := storage.LeaveCreate(&leave)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id}) //todo return in response leaveID
}

func LeavesStatusRead(c *gin.Context) {

	orgID := c.GetInt("organizationID")
	empID := c.GetInt("employeeID")

	status := c.Query("status")
	eid := c.Query("employeeID")

	var leaveRequests []model.Leave
	// svaki zaposleni moze da vidi svoje leave reqeuste, a admin moze da vidi sve leave requeste iz svoje organizacije
	filter := map[string]string{"organizationID": fmt.Sprintf("%d", orgID)}
	if status != "" {
		filter["status"] = status
	}
	if eid != "" {
		filter["employeeID"] = eid
	} else if !isAdmin(c) {
		filter["employeeID"] = fmt.Sprintf("%d", empID)
	}

	leaveRequests, err := storage.LeaveRead(filter)
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

	orgID := c.GetInt("organizationID")
	empID := c.GetInt("employeeID")

	leave, err := getLeaveRequest(c, leaveID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if leave.OrganizationId != orgID && leave.EmployeeId != empID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update this leave request"})
		return
	}

	if err := c.ShouldBindJSON(&leave); err != nil { //todo first bind then go to db and merge leave struct
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	leave.EmployeeId = empID
	leave.OrganizationId = orgID
	leave.Total = int(leave.EndDate.Sub(leave.StartDate).Hours() / 24)
	leave.Status = "pending"

	err = storage.LeaveUpdate(leave)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"status": "success", "message": "Leave request updated successfully"})

}

func LeaveStatusAdminUpdate(c *gin.Context) {
	leaveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	orgID := c.GetInt("organizationID")
	empID := c.GetInt("employeeID")

	if !isAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update status of the leave request"})
		return
	}

	leaveRequest, err := getLeaveRequest(c, leaveID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if leaveRequest.OrganizationId != orgID && leaveRequest.EmployeeId != empID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update this leave request"})
		return
	}

	var updatedLeave model.Leave
	if err := c.ShouldBindJSON(&updatedLeave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if updatedLeave.Status != "approved" && updatedLeave.Status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid status"})
		return
	}

	leaveRequest.Status = updatedLeave.Status

	err = storage.LeaveUpdateStatus(leaveRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update leave request"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"status": "success", "message": "Leave request updated successfully"})
}
