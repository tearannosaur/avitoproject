package tests

import (
	"bytes"
	"database/sql"
	"fmt"
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

func TestAddProduct_Success(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/product", h.AddProduct)

	pvzID := uuid.New()
	receptionID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "reception_status"}).
		AddRow(receptionID, now, pvzID, "Открыта")
	mock.ExpectQuery("SELECT id, date_time, pvz_id, reception_status FROM receptions").
		WithArgs(pvzID).
		WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO products").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Электроника", receptionID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := fmt.Sprintf(`{"type": "Электроника", "pvzId": "%s"}`, pvzID)
	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Товар успешно добавлен.")
}

func TestAddProduct_InvalidInput(t *testing.T) {
	h := &handler.Handler{}
	router := gin.Default()
	router.POST("/product", h.AddProduct)

	body := `{"type": "Еда", "pvzId": "not-a-uuid"}`
	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewBufferString(body))
	req.Header.Set("Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверный запрос")
}
func TestAddProduct_NoOpenReception(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/product", h.AddProduct)

	pvzID := uuid.New()

	mock.ExpectQuery("SELECT id, date_time, pvz_id, reception_status FROM receptions").
		WithArgs(pvzID).
		WillReturnError(sql.ErrNoRows)

	body := fmt.Sprintf(`{"type": "Обувь", "pvzId": "%s"}`, pvzID)
	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Нет активной приемки")
}

func TestAddProduct_InsertError(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/product", h.AddProduct)

	pvzID := uuid.New()
	receptionID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "reception_status"}).
		AddRow(receptionID, now, pvzID, "Открыта")
	mock.ExpectQuery("SELECT id, date_time, pvz_id, reception_status FROM receptions").
		WithArgs(pvzID).
		WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO products").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Одежда", receptionID).
		WillReturnError(sql.ErrConnDone)

	body := fmt.Sprintf(`{"type": "Одежда", "pvzId": "%s"}`, pvzID)
	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Не удалось добавить товар")
}
