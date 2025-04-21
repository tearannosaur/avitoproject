package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/project/pkg/models"
)

func (h *Handler) GetPVZList(c *gin.Context) {
	var startDate, endDate time.Time
	var err error

	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	now := time.Now()

	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат времени."})
			return
		}
	} else {
		startDate = now.AddDate(-1, 0, 0)
	}

	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат времени."})
			return
		}
	} else {
		endDate = now
	}

	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "startDate не может быть позже endDate."})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректное значение страницы."})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 30 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректное значение лимита."})
		return
	}

	startDateSQL := startDate.Format("2006-01-02 15:04:05")
	endDateSQL := endDate.Format("2006-01-02 15:04:05")

	var pvzList []models.PVZ
	query := `
		SELECT DISTINCT p.* FROM pvz p
		JOIN receptions r ON p.id = r.pvz_id
		WHERE r.date_time >= $1 AND r.date_time <= $2
		ORDER BY p.registration_date DESC
		LIMIT $3 OFFSET $4
	`

	err = h.Db.Select(&pvzList, query, startDateSQL, endDateSQL, limit, (page-1)*limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении PVZ."})
		return
	}

	var response []gin.H
	for _, pvz := range pvzList {
		var receptions []models.Reception

		receptionQuery := `
			SELECT * FROM receptions 
			WHERE pvz_id = $1 AND date_time >= $2 AND date_time <= $3
		`

		err := h.Db.Select(&receptions, receptionQuery, pvz.ID, startDateSQL, endDateSQL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении приемок."})
			return
		}

		var receptionData []gin.H
		for _, reception := range receptions {
			var products []models.Product

			productQuery := `SELECT * FROM products WHERE reception_id = $1`
			err := h.Db.Select(&products, productQuery, reception.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении товаров."})
				return
			}

			receptionData = append(receptionData, gin.H{
				"reception": reception,
				"products":  products,
			})
		}

		if len(receptionData) > 0 {
			response = append(response, gin.H{
				"pvz":        pvz,
				"receptions": receptionData,
			})
		}
	}

	if len(response) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Нет данных по указанному времени."})
		return
	}

	c.JSON(http.StatusOK, response)
}
