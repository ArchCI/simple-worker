package redisutil

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

const (
	SET_COMMAND  = "SET"
	HSET_COMMAND = "HSET"
)

func GetString(key string) string {
	c, err := redis.Dial("tcp", ":6379")
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

func WriteLogsToRedis(buildId int64, logs []string) bool {

	c, err := redis.Dial("tcp", ":6379")
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
