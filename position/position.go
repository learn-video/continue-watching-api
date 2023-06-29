package position

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func Record(r *redis.Client, userID, videoID string, position int) error {
	ctx := context.TODO()
	key := fmt.Sprintf("%s_%s", userID, videoID)
	r.Set(ctx, key, position, 1*time.Minute)
	return nil
}
