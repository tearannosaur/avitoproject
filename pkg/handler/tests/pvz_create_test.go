package tests

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/project/pkg/handler"
	"github.com/stretchr/testify/assert"
)

func TestCreatePvz_Success(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/pvz", h.CreatePvz)

	body := []byte(`{"city": "Москва"}`)
	mock.ExpectExec("INSERT INTO pvz").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Москва").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "ПВЗ создан.")
}

func TestCreatePvz_InvalidInput(t *testing.T) {
	h := &handler.Handler{}
	router := gin.Default()
	router.POST("/pvz", h.CreatePvz)

	body := []byte(`{"city": "Лондон"}`)
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверный запрос")
}

func TestCreatePvz_InsertError(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/pvz", h.CreatePvz)

	body := []byte(`{"city": "Казань"}`)
	mock.ExpectExec("INSERT INTO pvz").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Казань").
		WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Ошибка при сохранении в базу данных.")
}
