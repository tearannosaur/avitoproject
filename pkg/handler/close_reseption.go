package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CloseReception struct {
	ID       uuid.UUID `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzID    uuid.UUID `json:"pvzId"`
	Status   string    `json:"status" binding:"required,oneof=Открыта Закрыта"`
}

func (h *Handler) CloseLastReception(c *gin.Context) {

	pvzIDStr := c.Param("pvzId")
	pvzID, err := uuid.Parse(pvzIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос."})
		return
	}

	var reception CloseReception
	query := `SELECT id, date_time, pvz_id, reception_status FROM receptions
	          WHERE pvz_id = $1 AND reception_status = 'Открыта'
	          ORDER BY date_time DESC LIMIT 1`
	row := h.Db.QueryRow(query, pvzID)
	err = row.Scan(&reception.ID, &reception.DateTime, &reception.PvzID, &reception.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет активной приемки."})
		return
	}

	updateQuery := `UPDATE receptions SET reception_status = 'Закрыта' WHERE id = $1`
	_, err = h.Db.Exec(updateQuery, reception.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось закрыть приемку."})
		return
	}

	reception.Status = "Закрыта"
	c.JSON(http.StatusOK, reception)
}
