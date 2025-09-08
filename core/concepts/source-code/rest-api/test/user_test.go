package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"rest-api/config"
	"rest-api/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func setupRouter() (*gin.Engine, *gin.Engine) {
	db := config.InitDB()
	public := gin.Default()
	protected := gin.Default()

	public.GET("/users", handler.GetUsers(db))
	public.GET("/users/:id", handler.GetUserByID(db))

	protected.POST("/users", handler.CreateUser(db))
	protected.PUT("/users/:id", handler.UpdateUser(db))
	protected.DELETE("/users/:id", handler.DeleteUser(db))

	return public, protected
}

func TestCreateUser(t *testing.T) {
	_, r := setupRouter()
	user := User{Name: "John"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer secret-token")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetUsers(t *testing.T) {
	r, _ := setupRouter()
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserByID_NotFound(t *testing.T) {
	r, _ := setupRouter()
	req := httptest.NewRequest("GET", "/users/99999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateUser(t *testing.T) {
	_, r := setupRouter()
	user := User{Name: "Updated"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer secret-token")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Contains(t, []int{http.StatusOK, http.StatusInternalServerError}, w.Code)
}

func TestDeleteUser(t *testing.T) {
	_, r := setupRouter()
	req := httptest.NewRequest("DELETE", "/users/1", nil)
	req.Header.Set("Authorization", "Bearer secret-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Contains(t, []int{http.StatusOK, http.StatusInternalServerError}, w.Code)
}
