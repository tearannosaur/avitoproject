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
	"github.com/stretchr/testify/require"
)

func TestCloseLastReception_Success(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать mock для базы данных: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	h := handler.NewHandler(sqlxDB)

	rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "reception_status"}).
		AddRow("00000000-0000-0000-0000-000000000001", time.Now(), "5d8c679b-9997-4b33-83cf-e2d71179328d", "Открыта")

	mock.ExpectQuery(`SELECT id, date_time, pvz_id, reception_status FROM receptions`).
		WithArgs("5d8c679b-9997-4b33-83cf-e2d71179328d").
		WillReturnRows(rows)

	mock.ExpectExec(`UPDATE receptions SET reception_status = 'Закрыта' WHERE id = \$1`).
		WithArgs("00000000-0000-0000-0000-000000000001").
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := gin.Default()
	r.POST("/pvz/:pvzId/close_last_reception", h.CloseLastReception)

	req, _ := http.NewRequest("POST", "/pvz/5d8c679b-9997-4b33-83cf-e2d71179328d/close_last_reception", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "Закрыта")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания были выполнены: %v", err)
	}
}

func TestCloseLastReception_InvalidUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать mock для базы данных: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	h := &handler.Handler{Db: sqlxDB}

	router := gin.Default()
	router.POST("/pvz/:pvzId/close_last_reception", h.CloseLastReception)

	req := httptest.NewRequest(http.MethodPost, "/pvz/invalid-uuid/close_last_reception", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверный запрос")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания были выполнены: %v", err)
	}
}

func TestCloseLastReception_NoOpenReception(t *testing.T) {
	dbRaw, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbRaw.Close()

	db := sqlx.NewDb(dbRaw, "postgres")
	h := &handler.Handler{Db: db}

	router := gin.Default()
	router.POST("/pvz/:pvzId/close_last_reception", h.CloseLastReception)

	pvzID := uuid.New()

	mock.ExpectQuery(`SELECT id, date_time, pvz_id, reception_status FROM receptions WHERE pvz_id = \$1 AND reception_status = 'Открыта' ORDER BY date_time DESC LIMIT 1`).
		WithArgs(pvzID).
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/close_last_reception", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Нет активной приемки")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseLastReception_UpdateError(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/pvz/:pvzId/close_last_reception", h.CloseLastReception)

	pvzID := uuid.New()
	receptionID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "reception_status"}).
		AddRow(receptionID, time.Now(), pvzID, "Открыта")
	mock.ExpectQuery("SELECT id, date_time, pvz_id, reception_status FROM receptions").
		WithArgs(pvzID).
		WillReturnRows(rows)

	mock.ExpectExec("UPDATE receptions SET reception_status = 'Закрыта' WHERE id =").
		WithArgs(receptionID).
		WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/close_last_reception", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Не удалось закрыть приемку")
}
