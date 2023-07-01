package position_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/labstack/echo/v4"
	"github.com/learn-video/continue-watching-api/position"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRecordPositionOK(t *testing.T) {
	e := echo.New()
	positionJSON := `{"video_id": "123", "position": 1}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(positionJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(&http.Cookie{
		Name:  "user_id",
		Value: "bda031c0-4e7d-493a-92ba-6fc1eb3e6216",
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db, mock := redismock.NewClientMock()
	mock.Regexp().ExpectSet("bda031c0-4e7d-493a-92ba-6fc1eb3e6216_123", 1, 1*time.Minute).
		SetVal("OK")
	h := position.NewHandler(db)

	if assert.NoError(t, h.Record(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestRecordPositionMissingUserID(t *testing.T) {
	e := echo.New()
	positionJSON := `{"video_id": "123", "position": 1}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(positionJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db, _ := redismock.NewClientMock()
	h := position.NewHandler(db)
	h.Record(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRecordPositionMissingPayload(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(&http.Cookie{
		Name:  "user_id",
		Value: "bda031c0-4e7d-493a-92ba-6fc1eb3e6216",
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db, _ := redismock.NewClientMock()
	h := position.NewHandler(db)
	h.Record(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRecordPositionRedisError(t *testing.T) {
	e := echo.New()
	positionJSON := `{"video_id": "123", "position": 1}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(positionJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(&http.Cookie{
		Name:  "user_id",
		Value: "bda031c0-4e7d-493a-92ba-6fc1eb3e6216",
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db, mock := redismock.NewClientMock()
	mock.Regexp().ExpectSet("bda031c0-4e7d-493a-92ba-6fc1eb3e6216_123", 1, 1*time.Minute).
		SetErr(errors.New("failed to set key"))
	h := position.NewHandler(db)
	h.Record(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHandlerFetchOK(t *testing.T) {
	e := echo.New()
	q := make(url.Values)
	q.Set("video_id", "123")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	req.AddCookie(&http.Cookie{
		Name:  "user_id",
		Value: "bda031c0-4e7d-493a-92ba-6fc1eb3e6216",
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
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
	e := echo.New()
	q := make(url.Values)
	q.Set("video_id", "123")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db, _ := redismock.NewClientMock()
	h := position.NewHandler(db)

	h.Fetch(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandlerFetchPositionNotFound(t *testing.T) {
	e := echo.New()
	q := make(url.Values)
	q.Set("video_id", "123")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	req.AddCookie(&http.Cookie{
		Name:  "user_id",
		Value: "bda031c0-4e7d-493a-92ba-6fc1eb3e6216",
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db, mock := redismock.NewClientMock()
	mock.ExpectGet("bda031c0-4e7d-493a-92ba-6fc1eb3e6216_123").SetErr(redis.Nil)
	h := position.NewHandler(db)

	h.Fetch(c)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandlerFetchPositionRedisError(t *testing.T) {
	e := echo.New()
	q := make(url.Values)
	q.Set("video_id", "123")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	req.AddCookie(&http.Cookie{
		Name:  "user_id",
		Value: "bda031c0-4e7d-493a-92ba-6fc1eb3e6216",
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db, mock := redismock.NewClientMock()
	mock.ExpectGet("bda031c0-4e7d-493a-92ba-6fc1eb3e6216_123").SetErr(errors.New("failed to get"))
	h := position.NewHandler(db)

	h.Fetch(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
