package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/project/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4"`
	Role     string `json:"role" binding:"required,oneof=employee client"`
}

func (h *Handler) RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Ошибка при обработке запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка при хэшировании пароля: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хэшировании пароля."})
		return
	}

	user := models.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	err = h.saveUser(&user)
	if err != nil {
		log.Printf("Ошибка при сохранении пользователя в базу данных: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении в базу данных или пользователь существует."})
		return
	}

	log.Printf("Пользователь с ID: %s зарегистрирован.", user.ID.String())
	c.JSON(http.StatusCreated, gin.H{
		"message": "Пользователь успешно зарегистрирован",
	})
}

func (h *Handler) saveUser(user *models.User) error {

	query := `INSERT INTO users (id, email, user_role, password_hash) VALUES ($1, $2, $3, $4)`

	_, err := h.Db.Exec(query, user.ID, user.Email, user.Role, user.Password)
	return err
}
