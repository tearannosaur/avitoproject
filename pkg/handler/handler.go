package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/project/middleware"
)

type Handler struct {
	Db *sqlx.DB
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.POST("/dummyLogin", h.DummyLoginHandler)
	router.POST("/register", h.RegisterHandler)
	router.POST("/login", h.LoginHandler)
	employeeGroup := router.Group("/", middleware.EmployeeAuthMiddleware())
	{
		employeeGroup.POST("/pvz", h.CreatePvz)
		employeeGroup.GET("/pvz", h.GetPVZList)
		employeeGroup.POST("/pvz/:pvzId/close_last_reception", h.CloseLastReception)
		employeeGroup.POST("/receptions", h.CreateReception)
		employeeGroup.POST("/pvz/:pvzId/delete_last_product", h.DeleteLastProduct)
		employeeGroup.POST("/products", h.AddProduct)
	}

	return router

}
func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{Db: db}
}
