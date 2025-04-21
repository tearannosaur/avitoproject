package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/project/middleware"
	"github.com/project/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Ошибка при обработке запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос."})
		return
	}

	var user models.User
	err := h.Db.Get(&user, "SELECT * FROM users WHERE email = $1", req.Email)
	if err != nil {
		log.Printf("Ошибка при получении пользователя: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Printf("Неверный пароль: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные."})
		return
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
		"token":   token,
	})
}
