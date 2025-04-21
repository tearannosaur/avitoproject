package tests

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/project/pkg/handler"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	mock.ExpectExec(`INSERT INTO users`).
		WithArgs(sqlmock.AnyArg(), "test@example.com", "employee", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	router := gin.Default()
	router.POST("/register", h.RegisterHandler)

	body := []byte(`{
		"email": "test@example.com",
		"password": "1234",
		"role": "employee"
	}`)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Пользователь успешно зарегистрирован")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Не все ожидаемые запросы были выполнены: %v", err)
	}
}

func TestRegisterHandler_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	dbRaw, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock базы данных: %v", err)
	}
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}

	router := gin.Default()
	router.POST("/register", h.RegisterHandler)
	body := []byte(`{
		"email": "not-an-email",
		"password": "1234",
		"role": "employee"
	}`)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверный запрос")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Неожиданные запросы к базе данных: %v", err)
	}
}

func TestRegisterHandler_HashError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dbRaw, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock базы данных: %v", err)
	}
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/register", h.RegisterHandler)
	longPass := make([]byte, 73)
	for i := range longPass {
		longPass[i] = 'a'
	}

	body := []byte(`{
		"email": "test@example.com",
		"password": "` + string(longPass) + `",
		"role": "employee"
	}`)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Ошибка при хэшировании пароля")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Неожиданные запросы к базе данных: %v", err)
	}
}

func TestRegisterHandler_SaveUserError(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	mock.ExpectExec(`INSERT INTO users`).
		WithArgs(sqlmock.AnyArg(), "test@example.com", "employee", sqlmock.AnyArg()).
		WillReturnError(errors.New("mock db error"))

	router := gin.Default()
	router.POST("/register", h.RegisterHandler)

	body := []byte(`{
		"email": "test@example.com",
		"password": "1234",
		"role": "employee"
	}`)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Ошибка при сохранении в базу данных или пользователь существует.")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Не все ожидаемые запросы были выполнены: %v", err)
	}
}
