package position_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/learn-video/continue-watching-api/position"
	"github.com/stretchr/testify/assert"
)

func TestRecordOK(t *testing.T) {
	db, mock := redismock.NewClientMock()
	userID := "bda031c0-4e7d-493a-92ba-6fc1eb3e6216"
	videoID := "a1fbd2af-ab5e-44ac-9e5d-1a24051f89cf"
	key := fmt.Sprintf("%s_%s", userID, videoID)

	mock.ExpectSet(key, 1, 1*time.Minute).SetVal("OK")

	err := position.Record(db, userID, videoID, 1)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestRecordRedisError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	userID := "bda031c0-4e7d-493a-92ba-6fc1eb3e6216"
	videoID := "a1fbd2af-ab5e-44ac-9e5d-1a24051f89cf"
	key := fmt.Sprintf("%s_%s", userID, videoID)

	mock.ExpectSet(key, 1, 1*time.Minute).SetErr(errors.New("failed to set key"))

	err := position.Record(db, userID, videoID, 1)
	assert.NotNil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
