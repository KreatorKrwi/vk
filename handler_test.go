package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Auth(req *AuthReq) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func (m *MockService) Registration(req *AuthReq) (*User, error) {
	args := m.Called(req)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockService) NewObj(obj *ObjReqWLogin) (*ObjExport, error) {
	args := m.Called(obj)
	return args.Get(0).(*ObjExport), args.Error(1)
}

func (m *MockService) GetItems(filters *AdsFilters, login string) ([]Ad, error) {
	args := m.Called(filters, login)
	return args.Get(0).([]Ad), args.Error(1)
}

func TestHandler_Auth(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	testReq := &AuthReq{Login: "test", Password: "test123"}
	mockService.On("Auth", testReq).Return("test_token", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody, _ := json.Marshal(testReq)
	c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Auth(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "test_token", response["token"])

	mockService.AssertExpectations(t)
}

func TestHandler_Registration_InvalidLogin(t *testing.T) {
	handler := NewHandler(nil)

	tests := []struct {
		name        string
		login       string
		expectedErr string
	}{
		{"Login too long", "this_login_is_way_too_long_for_validation", "Login len must be < 20"},
		{"Login with spaces", "login with spaces", "Login contains spaces"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			reqBody, _ := json.Marshal(AuthReq{Login: tt.login, Password: "test123"})
			c.Request = httptest.NewRequest("POST", "/registration", bytes.NewBuffer(reqBody))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Registration(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedErr, response["error"])
		})
	}
}

func TestHandler_NewObj_Validation(t *testing.T) {
	handler := NewHandler(nil)

	tests := []struct {
		name        string
		obj         ObjReq
		expectedErr string
	}{
		{
			"Header too long",
			ObjReq{Header: "This header is definitely longer than 25 characters", Body: "test", Price: 100, Image: "test.jpg"},
			"Header len must be < 25",
		},
		{
			"Invalid image format",
			ObjReq{Header: "Test", Body: "test", Price: 100, Image: "test.png"},
			"Image must be .jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("user_login", "testuser")

			reqBody, _ := json.Marshal(tt.obj)
			c.Request = httptest.NewRequest("POST", "/new", bytes.NewBuffer(reqBody))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.NewObj(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedErr, response["error"])
		})
	}
}
