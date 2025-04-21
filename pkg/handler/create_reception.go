package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/project/pkg/models"
)

type CreateReceptionRequest struct {
	PvzID uuid.UUID `json:"pvzId" binding:"required"`
}

func (h *Handler) CreateReception(c *gin.Context) {

	var req CreateReceptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос."})
		return
	}
	var exists bool
	err := h.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM pvz WHERE id = $1)", req.PvzID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ПВЗ не найден."})
		return
	}

	var openExists bool

	err = h.Db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM receptions 
			WHERE pvz_id = $1 AND reception_status = 'Открыта'
		)
	`, req.PvzID).Scan(&openExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке приемки."})
		return
	}

	if openExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Есть незакрытая приемка."})
		return
	}

	id := uuid.New()
	dateTime := time.Now()
	status := "Открыта"

	_, err = h.Db.Exec(`
		INSERT INTO receptions (id, date_time, pvz_id, reception_status) 
		VALUES ($1, $2, $3, $4)
	`, id, dateTime, req.PvzID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании приемки."})
		return
	}
	c.JSON(http.StatusCreated, models.Reception{
		ID:       id,
		DateTime: dateTime,
		PvzID:    req.PvzID,
		Status:   status,
	})

}
