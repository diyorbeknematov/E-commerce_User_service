package redis_1

import 	"github.com/redis/go-redis/v9"


func RedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return client
}