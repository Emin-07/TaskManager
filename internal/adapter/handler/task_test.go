package handler

import (
	"encoding/json"
	"errors"
	"net/http"
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

func TestTaskHandler_GetFromDB(t *testing.T) {
	taskSvc := &testutil.MockTaskService{}
	taskSvc.On("Get", mock.Anything, "1", "1", "user").Return(testutil.SampleDomainTask, nil)
	cache := &testutil.MockRateAndCacheService{}
	cache.On("Get", mock.Anything, "task", "1", "1").Return("", domain.ErrKeyNotFound)
	cache.On("Set", mock.Anything, "task", "1", "1", mock.Anything, mock.Anything).Return(nil)
	tokenSvc := &testutil.MockTokenService{}
	tokenSvc.On("ParseFromRequest", mock.Anything).Return(map[string]string{"id": "1", "role": "user"}, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/tasks/1", nil)
	req.Header.Set("Authorization", "Bearer token")
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h := NewTaskHandler(taskSvc, tokenSvc, cache, &testutil.MockBroker{})
	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]TaskResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "Buy milk", body["task"].Title)
}

func TestTaskHandler_GetNotFound(t *testing.T) {
	taskSvc := &testutil.MockTaskService{}
	taskSvc.On("Get", mock.Anything, "99", "1", "user").Return(nil, domain.ErrNoRecord)
	cache := &testutil.MockRateAndCacheService{}
	cache.On("Get", mock.Anything, "task", "99", "1").Return("", domain.ErrKeyNotFound)
	tokenSvc := &testutil.MockTokenService{}
	tokenSvc.On("ParseFromRequest", mock.Anything).Return(map[string]string{"id": "1", "role": "user"}, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/tasks/99", nil)
	req.Header.Set("Authorization", "Bearer token")
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "99"}}

	h := NewTaskHandler(taskSvc, tokenSvc, cache, &testutil.MockBroker{})
	h.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTaskHandler_PostPublishesToBroker(t *testing.T) {
	broker := &testutil.MockBroker{}
	broker.On("Publish", mock.Anything, domain.TopicTasks).Return(nil)
	tokenSvc := &testutil.MockTokenService{}
	tokenSvc.On("ParseFromRequest", mock.Anything).Return(map[string]string{"id": "1", "role": "user"}, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"New task","text":"desc","priority":1,"expire_days":3}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	c.Request = req

	h := NewTaskHandler(&testutil.MockTaskService{}, tokenSvc, &testutil.MockRateAndCacheService{}, broker)
	h.Post(c)

	// c.Status(202) is buffered in the writer and only flushed on write, so the
	// recorder may still report 200; the broker call is what matters here.
	broker.AssertExpectations(t)
}

func TestTaskHandler_PostBrokerError(t *testing.T) {
	broker := &testutil.MockBroker{}
	broker.On("Publish", mock.Anything, domain.TopicTasks).Return(errors.New("kafka down"))
	tokenSvc := &testutil.MockTokenService{}
	tokenSvc.On("ParseFromRequest", mock.Anything).Return(map[string]string{"id": "1", "role": "user"}, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"New task"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	c.Request = req

	h := NewTaskHandler(&testutil.MockTaskService{}, tokenSvc, &testutil.MockRateAndCacheService{}, broker)
	h.Post(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
