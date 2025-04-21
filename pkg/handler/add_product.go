package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/project/pkg/models"
)

type AddProductRequest struct {
	Type  string    `json:"type" binding:"required,oneof=Электроника Одежда Обувь"`
	PvzID uuid.UUID `json:"pvzId" binding:"required,uuid"`
}

func (h *Handler) AddProduct(c *gin.Context) {
	var req AddProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный запрос.",
		})
		return
	}

	var reception models.Reception
	query := `SELECT id, date_time, pvz_id, reception_status FROM receptions
	          WHERE pvz_id = $1 AND reception_status = 'Открыта'
	          ORDER BY date_time DESC LIMIT 1`
	row := h.Db.QueryRow(query, req.PvzID)
	err := row.Scan(&reception.ID, &reception.DateTime, &reception.PvzID, &reception.Status)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Нет активной приемки для указанного PVZ.",
		})
		return
	}

	product := models.Product{
		ID:          uuid.New(),
		Type:        req.Type,
		ReceptionID: reception.ID,
		DateTime:    time.Now(),
	}

	insertQuery := `INSERT INTO products (id, date_time, product_type, reception_id) VALUES ($1, $2, $3, $4)`
	_, err = h.Db.Exec(insertQuery, product.ID, product.DateTime, product.Type, product.ReceptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось добавить товар.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Товар успешно добавлен.",
		"product": product,
	})
}
