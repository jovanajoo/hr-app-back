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
	protected.POST("request-leave/:id", LeaveCreate)
	protected.PATCH("request-leave/:id", LeaveUpdate)
	protected.PATCH("request-leave/:id/status", LeaveStatusAdminUpdate)

	// departments
	r.GET("departments", ReadDepartments)

	// position
	r.GET("positions", PositionRead)

	return r

}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
