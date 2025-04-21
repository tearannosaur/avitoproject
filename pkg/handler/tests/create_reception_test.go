package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/project/pkg/handler"
	"github.com/stretchr/testify/assert"
)

func TestCreateReception_Success(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/reception", h.CreateReception)

	pvzID := uuid.New()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec("INSERT INTO receptions").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), pvzID, "Открыта").
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := `{"pvzId": "` + pvzID.String() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/reception", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"Открыта"`)
}

func TestCreateReception_PvzNotFound(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/reception", h.CreateReception)

	pvzID := uuid.New()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	body := `{"pvzId": "` + pvzID.String() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/reception", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ПВЗ не найден")
}
func TestCreateReception_AlreadyOpen(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/reception", h.CreateReception)

	pvzID := uuid.New()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	body := `{"pvzId": "` + pvzID.String() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/reception", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Есть незакрытая приемка")
}
func TestCreateReception_InsertError(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/reception", h.CreateReception)

	pvzID := uuid.New()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec("INSERT INTO receptions").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), pvzID, "Открыта").
		WillReturnError(errors.New("fail"))

	body := `{"pvzId": "` + pvzID.String() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/reception", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Ошибка при создании приемки")
}
