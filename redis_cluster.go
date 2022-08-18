package main

import (
	"context"
	errs "errors"
	"github.com/go-redis/redis/v9"
	"strings"
	"time"
)

func ClusterConnect(redisHost string) *redis.ClusterClient {
	addr := strings.Split(redisHost, ",")
	rdb := redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: addr,
		RouteRandomly: true,
		ReadTimeout:   5 * time.Second,
	})

	if err := rdb.Ping(context.TODO()).Err(); err != nil {
		panic("unable to connect to redis " + err.Error())
	}

	return rdb
}

func ClusterDelete(conn *redis.ClusterClient, key string) error {
	if conn == nil {
		return errs.New("redis cluster connection is required")
	}

	if key == "" {
		return errs.New("redis key is missing")
	}

	return conn.Del(context.TODO(), key).Err()
}
