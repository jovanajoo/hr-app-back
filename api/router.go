package api

import "github.com/gin-gonic/gin"

func SetUpRouters() *gin.Engine {
	r := gin.Default()

	protected := r.Group("/api")
	protected.Use(authMiddleware("secretkey"))

	// organization
	r.GET("/organization", OrganizationGet)
	r.POST("/organization", OrganizationInsert)
	r.PATCH("organization/:id", OrganizationUpdate)
	r.DELETE("/organization/:id", OrganizationDelete)

	// employee
	protected.GET("/employee", EmployeeGet)
	r.POST("/login", EmployeeLogin)
	protected.POST("/employee", EmployeeInsert)
	protected.PATCH("employee/:id", EmployeeUpdate)
	protected.DELETE("employee/:id", EmployeeDelete)

	return r

}
