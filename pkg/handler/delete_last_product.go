package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/project/pkg/models"
)

func (h *Handler) DeleteLastProduct(c *gin.Context) {
	pvzIDStr := c.Param("pvzId")
	pvzID, err := uuid.Parse(pvzIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос."})
		return
	}

	var reception models.Reception
	query := `SELECT id, date_time, pvz_id, reception_status FROM receptions
	          WHERE pvz_id = $1 AND reception_status = 'Открыта'
	          ORDER BY date_time DESC LIMIT 1`
	row := h.Db.QueryRow(query, pvzID)
	err = row.Scan(&reception.ID, &reception.DateTime, &reception.PvzID, &reception.Status)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет активной приемки для указанного PVZ."})
		return
	}

	var products []models.Product
	rows, err := h.Db.Query(`SELECT id, date_time, product_type, reception_id FROM products WHERE reception_id = $1`, reception.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении товаров."})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении товаров."})
			return
		}
		products = append(products, product)
	}

	if len(products) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет товаров для удаления."})
		return
	}

	lastProduct := products[len(products)-1]
	_, err = h.Db.Exec(`DELETE FROM products WHERE id = $1`, lastProduct.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить товар."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Товар успешно удален.",
		"product": lastProduct,
	})
}
