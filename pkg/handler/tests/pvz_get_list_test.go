package tests

import (
	"database/sql"
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

func TestGetPVZList_Success(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.GET("/pvz", h.GetPVZList)

	pvzID := uuid.New()
	receptionID := uuid.New()
	productID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT DISTINCT p.* FROM pvz p").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
			AddRow(pvzID, now, "Москва"))

	mock.ExpectQuery("SELECT \\* FROM receptions").
		WithArgs(pvzID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "reception_status"}).
			AddRow(receptionID, now, pvzID, "Открыта"))

	mock.ExpectQuery("SELECT \\* FROM products").
		WithArgs(receptionID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "product_type", "reception_id"}).
			AddRow(productID, now, "Одежда", receptionID))

	req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Москва")
	assert.Contains(t, w.Body.String(), "Одежда")
}

func TestGetPVZList_EmptyResult(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.GET("/pvz", h.GetPVZList)

	mock.ExpectQuery("SELECT DISTINCT p.* FROM pvz p").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}))

	req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Нет данных по указанному времени")
}

func TestGetPVZList_InvalidDate(t *testing.T) {
	h := &handler.Handler{}
	router := gin.Default()
	router.GET("/pvz", h.GetPVZList)

	req := httptest.NewRequest(http.MethodGet, "/pvz?startDate=неправильная-дата", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверный формат времени")
}

func TestGetPVZList_MainQueryError(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.GET("/pvz", h.GetPVZList)

	mock.ExpectQuery("SELECT DISTINCT p.* FROM pvz p").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 10, 0).
		WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Ошибка при получении PVZ")
}
