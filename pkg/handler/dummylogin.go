package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/project/middleware"
	"github.com/project/pkg/models"
)

type DummyLoginRequest struct {
	Role string `json:"role" binding:"required,oneof=client employee"`
}

func (h *Handler) DummyLoginHandler(c *gin.Context) {
	var reg DummyLoginRequest
	if err := c.ShouldBindJSON(&reg); err != nil {
		log.Printf("Ошибка при обработке запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос."})
		return
	}
	user := models.User{
		ID:   uuid.New(),
		Role: reg.Role,
	}
	token, err := middleware.GenerateToken(user)
	if err != nil {
		log.Printf("Ошибка при генерации токена: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена."})
		return
	}
	log.Printf("Пользователь с ID: %s авторизован, токен сгенерирован.", user.ID.String())
	c.JSON(http.StatusOK, gin.H{
		"message": "Успешная авторизация",
		"token":   token})

}
