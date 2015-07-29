package redisutil

import (
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
)

const (
	ENV_REDIS_SERVER = "REDIS_SERVER"

	SET_COMMAND  = "SET"
	HSET_COMMAND = "HSET"
)

// GetRedisServer reads the environment variable to return the address of redis.
func GetRedisServer() string {
	if os.Getenv(ENV_REDIS_SERVER) != "" {
		return os.Getenv(ENV_REDIS_SERVER)
	} else {
		return "127.0.0.1:6379"
	}
}

// GetString performs get command to return string.
func GetString(key string) string {
	c, err := redis.Dial("tcp", GetRedisServer())
	if err != nil {
		panic(err)
	}
	defer c.Close()

	value, err := redis.String(c.Do("GET", key))
	if err != nil {
		fmt.Println("key not found")
	}

	return value
}

// WriteLogsToRedis take the array of log to write to redis.
func WriteLogsToRedis(buildId int64, logs []string) bool {

	c, err := redis.Dial("tcp", GetRedisServer())
	if err != nil {
		panic(err)
	}
	defer c.Close()

	index := 0

	c.Do(HSET_COMMAND, buildId, "finish", false)

	//Set one build's logs in one hash
	for _, log := range logs {
		fmt.Println(log)

		c.Do(HSET_COMMAND, buildId, index, log)

		c.Do(HSET_COMMAND, buildId, "current", index)

		index++
	}

	c.Do(HSET_COMMAND, buildId, "finish", true)

	return true
}
