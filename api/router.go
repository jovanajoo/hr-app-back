package api

import "github.com/gin-gonic/gin"

func SetUpRouters() *gin.Engine {
	r := gin.Default()
	r.Use(CORS())

	protected := r.Group("/api")
	protected.Use(authMiddleware("secretkey"))

	// organization
	protected.GET("/organization", OrganizationReadById)
	r.POST("/registration", OrganizationRegister)
	protected.PATCH("organization/:id", OrganizationUpdate)
	r.DELETE("/organization/:id", OrganizationDelete)

	// employee
	protected.GET("/employee", EmployeeReadByOrg)
	r.POST("/login", EmployeeLogin)
	protected.POST("/employee", EmployeeInsert)
	protected.PATCH("employee/:id", EmployeeUpdate)
	protected.DELETE("employee/:id", EmployeeDelete)

	// leave
	protected.GET("leaves/status", LeavesStatusRead)
	protected.POST("request-leave", LeaveCreate)
	protected.PATCH("request-leave/:id", LeaveUpdate)
	protected.PATCH("request-leave/:id/status", LeaveStatusAdminUpdate)

	// departments
	r.GET("departments", DepartmentsRead)

	// position
	r.GET("positions", PositionRead)

	return r

}
