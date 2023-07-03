package position_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/learn-video/continue-watching-api/position"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRecordOK(t *testing.T) {
	db, mock := redismock.NewClientMock()
	userID := "bda031c0-4e7d-493a-92ba-6fc1eb3e6216"
	videoID := "a1fbd2af-ab5e-44ac-9e5d-1a24051f89cf"
	key := fmt.Sprintf("%s_%s", userID, videoID)

	mock.ExpectSet(key, 1.0, 1*time.Minute).SetVal("OK")

	err := position.Record(db, userID, videoID, 1.0)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestRecordRedisError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	userID := "bda031c0-4e7d-493a-92ba-6fc1eb3e6216"
	videoID := "a1fbd2af-ab5e-44ac-9e5d-1a24051f89cf"
	key := fmt.Sprintf("%s_%s", userID, videoID)
	mock.ExpectSet(key, 1.0, 1*time.Minute).SetErr(errors.New("failed to set key"))

	err := position.Record(db, userID, videoID, 1.0)

	assert.NotNil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestFetchOK(t *testing.T) {
	db, mock := redismock.NewClientMock()
	userID := "bda031c0-4e7d-493a-92ba-6fc1eb3e6216"
	videoID := "a1fbd2af-ab5e-44ac-9e5d-1a24051f89cf"
	key := fmt.Sprintf("%s_%s", userID, videoID)
	mock.ExpectGet(key).SetVal("1")

	pos, err := position.Fetch(db, userID, videoID)

	assert.Nil(t, err)
	assert.Equal(t, pos, 1)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestFetchPositionNotFound(t *testing.T) {
	db, mock := redismock.NewClientMock()
	userID := "bda031c0-4e7d-493a-92ba-6fc1eb3e6216"
	videoID := "a1fbd2af-ab5e-44ac-9e5d-1a24051f89cf"
	key := fmt.Sprintf("%s_%s", userID, videoID)
	mock.ExpectGet(key).SetErr(redis.Nil)

	pos, err := position.Fetch(db, userID, videoID)

	assert.Equal(t, err, position.ErrNotFound)
	assert.Equal(t, pos, 0)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestFetchRedisError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	userID := "bda031c0-4e7d-493a-92ba-6fc1eb3e6216"
	videoID := "a1fbd2af-ab5e-44ac-9e5d-1a24051f89cf"
	key := fmt.Sprintf("%s_%s", userID, videoID)
	mock.ExpectGet(key).SetErr(errors.New("failed to get key"))

	pos, err := position.Fetch(db, userID, videoID)

	assert.Equal(t, err, errors.New("failed to get key"))
	assert.Equal(t, pos, 0)
	assert.Nil(t, mock.ExpectationsWereMet())
}
