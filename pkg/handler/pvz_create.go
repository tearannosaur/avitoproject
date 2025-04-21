package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreatePvzRequest struct {
	City             string    `json:"city" binding:"required,oneof=Москва Санкт-Петербург Казань"`
	RegistrationDate time.Time `json:"registrationDate"`
	ID               uuid.UUID `json:"id"`
}

func (h *Handler) CreatePvz(c *gin.Context) {
	var pvz CreatePvzRequest
	if err := c.ShouldBindJSON(&pvz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос."})
		return
	}
	pvz.ID = uuid.New()
	pvz.RegistrationDate = time.Now()

	_, err := h.Db.Exec(`
	INSERT INTO pvz (id,registration_date,city)
	VALUES ($1,$2,$3)
	`, pvz.ID, pvz.RegistrationDate, pvz.City)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении в базу данных."})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "ПВЗ создан.",
		"pvz":     pvz,
	})

}
