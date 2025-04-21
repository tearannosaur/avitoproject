package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/project/pkg/handler"
	"github.com/stretchr/testify/assert"
)

func TestDummyLoginHandler_InvalidRole(t *testing.T) {
	h := &handler.Handler{}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/dummy_login", h.DummyLoginHandler)

	body := handler.DummyLoginRequest{Role: "admin"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/dummy_login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверный запрос")
}

func TestDummyLoginHandler_MissingRole(t *testing.T) {
	h := &handler.Handler{}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/dummy_login", h.DummyLoginHandler)

	req := httptest.NewRequest(http.MethodPost, "/dummy_login", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверный запрос")
}
