package tests

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/project/pkg/handler"
	"github.com/stretchr/testify/assert"
)

func TestDeleteLastProduct_Success(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.DELETE("/reception/:pvzId/product", h.DeleteLastProduct)

	pvzID := uuid.New()
	receptionID := uuid.New()
	productID := uuid.New()
	now := time.Now()

	mock.ExpectQuery(`SELECT id, date_time, pvz_id, reception_status FROM receptions`).
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "reception_status"}).
			AddRow(receptionID, now, pvzID, "Открыта"))

	mock.ExpectQuery(`SELECT id, date_time, product_type, reception_id FROM products`).
		WithArgs(receptionID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "product_type", "reception_id"}).
			AddRow(productID, now, "Товар", receptionID))

	mock.ExpectExec(`DELETE FROM products WHERE id = \$1`).
		WithArgs(productID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodDelete, "/reception/"+pvzID.String()+"/product", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Товар успешно удален")
}
func TestDeleteLastProduct_InvalidUUID(t *testing.T) {
	h := &handler.Handler{Db: nil}

	router := gin.Default()
	router.DELETE("/reception/:pvzId/product", h.DeleteLastProduct)

	req := httptest.NewRequest(http.MethodDelete, "/reception/not-a-uuid/product", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверный запрос")
}

func TestDeleteLastProduct_NoReception(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.DELETE("/reception/:pvzId/product", h.DeleteLastProduct)

	pvzID := uuid.New()

	mock.ExpectQuery(`SELECT id, date_time, pvz_id, reception_status FROM receptions`).
		WithArgs(pvzID).
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodDelete, "/reception/"+pvzID.String()+"/product", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Нет активной приемки")
}
func TestDeleteLastProduct_NoProducts(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.DELETE("/reception/:pvzId/product", h.DeleteLastProduct)

	pvzID := uuid.New()
	receptionID := uuid.New()

	mock.ExpectQuery(`SELECT id, date_time, pvz_id, reception_status FROM receptions`).
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "reception_status"}).
			AddRow(receptionID, time.Now(), pvzID, "Открыта"))

	mock.ExpectQuery(`SELECT id, date_time, product_type, reception_id FROM products`).
		WithArgs(receptionID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "product_type", "reception_id"}))

	req := httptest.NewRequest(http.MethodDelete, "/reception/"+pvzID.String()+"/product", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Нет товаров для удаления")
}
func TestDeleteLastProduct_DeleteError(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.DELETE("/reception/:pvzId/product", h.DeleteLastProduct)

	pvzID := uuid.New()
	receptionID := uuid.New()
	productID := uuid.New()

	mock.ExpectQuery(`SELECT id, date_time, pvz_id, reception_status FROM receptions`).
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "reception_status"}).
			AddRow(receptionID, time.Now(), pvzID, "Открыта"))

	mock.ExpectQuery(`SELECT id, date_time, product_type, reception_id FROM products`).
		WithArgs(receptionID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "product_type", "reception_id"}).
			AddRow(productID, time.Now(), "Товар", receptionID))

	mock.ExpectExec(`DELETE FROM products WHERE id = \$1`).
		WithArgs(productID).
		WillReturnError(errors.New("delete error"))

	req := httptest.NewRequest(http.MethodDelete, "/reception/"+pvzID.String()+"/product", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Не удалось удалить товар")
}
