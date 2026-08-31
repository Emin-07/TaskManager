package handler

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Emin-07/TaskManager/internal/core/domain"
	"github.com/Emin-07/TaskManager/internal/testutil"
)

func TestUserHandler_Authenticate(t *testing.T) {
	userSvc := &testutil.MockUserService{}
	userSvc.On("Authenticate", mock.Anything, "alice@example.com", "password1").Return(testutil.SampleDomainUser, nil)
	tokenSvc := &testutil.MockTokenService{}
	tokenSvc.On("CreateToken", "1", "admin").Return("jwt-token", nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("POST", "/login", strings.NewReader("email=alice%40example.com&password=password1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	u := NewUserHandler(userSvc, tokenSvc, &testutil.MockRateAndCacheService{}, &testutil.MockBroker{})
	u.Authenticate(c)

	assert.Equal(t, 200, w.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "jwt-token", body["token"])
}

func TestUserHandler_AuthenticateBadCreds(t *testing.T) {
	userSvc := &testutil.MockUserService{}
	userSvc.On("Authenticate", mock.Anything, "alice@example.com", "wrongpass").Return(nil, domain.ErrNoRecord)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("POST", "/login", strings.NewReader("email=alice%40example.com&password=wrongpass"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	u := NewUserHandler(userSvc, &testutil.MockTokenService{}, &testutil.MockRateAndCacheService{}, &testutil.MockBroker{})
	u.Authenticate(c)

	assert.Equal(t, 400, w.Code)
}

func TestUserHandler_GetById(t *testing.T) {
	userSvc := &testutil.MockUserService{}
	userSvc.On("GetById", mock.Anything, "1").Return(testutil.SampleDomainUser, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/users/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	u := NewUserHandler(userSvc, &testutil.MockTokenService{}, &testutil.MockRateAndCacheService{}, &testutil.MockBroker{})
	u.GetById(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "alice")
}

func TestUserHandler_ListUsers(t *testing.T) {
	userSvc := &testutil.MockUserService{}
	userSvc.On("List", mock.Anything).Return([]*domain.User{testutil.SampleDomainUser}, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/users", nil)

	u := NewUserHandler(userSvc, &testutil.MockTokenService{}, &testutil.MockRateAndCacheService{}, &testutil.MockBroker{})
	u.ListUsers(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "alice")
}

func TestUserHandler_SignUpPublishesToBroker(t *testing.T) {
	broker := &testutil.MockBroker{}
	broker.On("Publish", mock.Anything, domain.TopicUsers).Return(nil)
	tokenSvc := &testutil.MockTokenService{}
	tokenSvc.On("ParseFromRequest", mock.Anything).Return(map[string]string{"id": "1", "role": "user"}, nil)
	tokenSvc.On("CreateToken", "1", "user").Return("jwt-token", nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(`{"username":"dave","role":"user","email":"dave@example.com","password":"password1"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	c.Request = req

	u := NewUserHandler(&testutil.MockUserService{}, tokenSvc, &testutil.MockRateAndCacheService{}, broker)
	u.SignUp(c)

	assert.Equal(t, 200, w.Code)
	broker.AssertExpectations(t)
}

func TestUserHandler_PatchPublishesToBroker(t *testing.T) {
	broker := &testutil.MockBroker{}
	broker.On("Publish", mock.Anything, domain.TopicUsers).Return(nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("PATCH", "/users/1", strings.NewReader(`{"username":"dave","role":"user","email":"dave@example.com","password":"password1"}`))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	u := NewUserHandler(&testutil.MockUserService{}, &testutil.MockTokenService{}, &testutil.MockRateAndCacheService{}, broker)
	u.Patch(c)

	// c.Status(202) is buffered and only flushed on body write; the broker call
	// is the important assertion here.
	broker.AssertExpectations(t)
}
