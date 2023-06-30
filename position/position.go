package position

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFound = errors.New("no position tracked")
)

func Record(r *redis.Client, userID, videoID string, position int) error {
	key := fmt.Sprintf("%s_%s", userID, videoID)
	return r.Set(context.TODO(), key, position, 1*time.Minute).Err()
}

func Fetch(r *redis.Client, userID, videoID string) (int, error) {
	key := fmt.Sprintf("%s_%s", userID, videoID)
	val, err := r.Get(context.TODO(), key).Result()
	if err == redis.Nil {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	pos, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return pos, nil
}
