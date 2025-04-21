package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/project/pkg/handler"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dbRaw, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	password := "123456"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mock.ExpectQuery(`SELECT \* FROM users WHERE email = \$1`).
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash", "user_role"}).
			AddRow("11111111-1111-1111-1111-111111111111", "test@example.com", string(hashedPassword), "employee"))

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/login", h.LoginHandler)

	body := []byte(`{
		"email": "test@example.com",
		"password": "123456"
	}`)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Успешная авторизация")
	assert.Contains(t, w.Body.String(), "token")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginHandler_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dbRaw, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/login", h.LoginHandler)

	body := []byte(`{
		"email": "invalid-email",
		"password": "123"
	}`)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверный запрос.")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dbRaw, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbRaw.Close()
	db := sqlx.NewDb(dbRaw, "postgres")

	correctPass := "correctpassword"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(correctPass), bcrypt.DefaultCost)

	mock.ExpectQuery(`SELECT \* FROM users WHERE email = \$1`).
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "role"}).
			AddRow("11111111-1111-1111-1111-111111111111", "test@example.com", string(hashed), "employee"))

	h := &handler.Handler{Db: db}
	router := gin.Default()
	router.POST("/login", h.LoginHandler)

	body := []byte(`{
		"email": "test@example.com",
		"password": "wrongpassword"
	}`)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Неверные данные.")

	assert.NoError(t, mock.ExpectationsWereMet())
}
