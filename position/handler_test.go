package position_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/labstack/echo/v4"
	"github.com/learn-video/continue-watching-api/position"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type contextSetup struct {
	Method  string
	Path    string
	Body    string
	Cookies []*http.Cookie
}

func setupContext(setup contextSetup) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(setup.Method, setup.Path, strings.NewReader(setup.Body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for _, cookie := range setup.Cookies {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestRecordPositionOK(t *testing.T) {
	setup := contextSetup{
		Method: http.MethodPost,
		Path:   "/",
		Body:   `{"video_id": "123", "position": 1.0}`,
		Cookies: []*http.Cookie{
			{Name: "user_id", Value: "bda031c0-4e7d-493a-92ba-6fc1eb3e6216"},
		},
	}
	c, rec := setupContext(setup)
	db, mock := redismock.NewClientMock()
	mock.Regexp().ExpectSet("bda031c0-4e7d-493a-92ba-6fc1eb3e6216_123", 1.0, 1*time.Minute).
		SetVal("OK")
	h := position.NewHandler(db)

	if assert.NoError(t, h.Record(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestRecordPositionMissingUserID(t *testing.T) {
	setup := contextSetup{
		Method: http.MethodPost,
		Path:   "/",
		Body:   `{"video_id": "123", "position": 1}`,
	}
	c, rec := setupContext(setup)
	db, _ := redismock.NewClientMock()
	h := position.NewHandler(db)
	h.Record(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRecordPositionMissingPayload(t *testing.T) {
	setup := contextSetup{
		Method: http.MethodPost,
		Path:   "/",
	}
	c, rec := setupContext(setup)
	req := c.Request()
	req.AddCookie(&http.Cookie{Name: "user_id", Value: "bda031c0-4e7d-493a-92ba-6fc1eb3e6216"})
	db, _ := redismock.NewClientMock()
	h := position.NewHandler(db)
	h.Record(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRecordPositionRedisError(t *testing.T) {
	setup := contextSetup{
		Method: http.MethodPost,
		Path:   "/",
		Body:   `{"video_id": "123", "position": 1}`,
		Cookies: []*http.Cookie{
			{Name: "user_id", Value: "bda031c0-4e7d-493a-92ba-6fc1eb3e6216"},
		},
	}
	c, rec := setupContext(setup)
	db, mock := redismock.NewClientMock()
	mock.Regexp().ExpectSet("bda031c0-4e7d-493a-92ba-6fc1eb3e6216_123", 1, 1*time.Minute).
		SetErr(errors.New("failed to set key"))
	h := position.NewHandler(db)
	h.Record(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHandlerFetchOK(t *testing.T) {
	setup := contextSetup{
		Method: http.MethodGet,
		Path:   "/?video_id=123",
		Cookies: []*http.Cookie{
			{Name: "user_id", Value: "bda031c0-4e7d-493a-92ba-6fc1eb3e6216"},
		},
	}
	c, rec := setupContext(setup)
	db, mock := redismock.NewClientMock()
	mock.ExpectGet("bda031c0-4e7d-493a-92ba-6fc1eb3e6216_123").SetVal("1")
	h := position.NewHandler(db)

	expectedJSON := `{"position": 1}`
	if assert.NoError(t, h.Fetch(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, expectedJSON, rec.Body.String())
	}
}

func TestHandlerMissingUserID(t *testing.T) {
	setup := contextSetup{
		Method: http.MethodGet,
		Path:   "/?video_id=123",
	}
	c, rec := setupContext(setup)
	db, _ := redismock.NewClientMock()
	h := position.NewHandler(db)

	h.Fetch(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandlerFetchPositionNotFound(t *testing.T) {
	setup := contextSetup{
		Method: http.MethodGet,
		Path:   "/?video_id=123",
		Cookies: []*http.Cookie{
			{Name: "user_id", Value: "bda031c0-4e7d-493a-92ba-6fc1eb3e6216"},
		},
	}
	c, rec := setupContext(setup)
	db, mock := redismock.NewClientMock()
	mock.ExpectGet("bda031c0-4e7d-493a-92ba-6fc1eb3e6216_123").SetErr(redis.Nil)
	h := position.NewHandler(db)

	h.Fetch(c)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandlerFetchPositionRedisError(t *testing.T) {
	setup := contextSetup{
		Method: http.MethodGet,
		Path:   "/?video_id=123",
		Cookies: []*http.Cookie{
			{Name: "user_id", Value: "bda031c0-4e7d-493a-92ba-6fc1eb3e6216"},
		},
	}
	c, rec := setupContext(setup)
	db, mock := redismock.NewClientMock()
	mock.ExpectGet("bda031c0-4e7d-493a-92ba-6fc1eb3e6216_123").SetErr(errors.New("failed to get"))
	h := position.NewHandler(db)

	h.Fetch(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
