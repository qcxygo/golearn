package controller

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rc *redis.Client

func init() {

	rc = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1" + ":" + "6379",
		Password: "",
		DB:       0,
	})
}
func RCGet(key string) *redis.StringCmd {
	return rc.Get(ctx, key)
}
func RCExists(key string) bool {
	return rc.Exists(ctx, key).Val() != 0
}
func RCSet(key string, value interface{}, expiration time.Duration) {
	if RCExists(key) {
		rc.Expire(ctx, key, expiration)
		return
	}
	rc.Set(ctx, key, value, expiration)
}
