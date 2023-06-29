package position

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Position struct {
	UserID   string `json:"user_id" validate:"required"`
	VideoID  string `json:"video_id" validate:"required"`
	Position int    `json:"position" validate:"required"`
}

func Record(r *redis.Client, userID, videoID string, position int) error {
	ctx := context.TODO()
	key := fmt.Sprintf("%s_%s", userID, videoID)
	return r.Set(ctx, key, position, 1*time.Minute).Err()
}
