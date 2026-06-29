package main

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestOfficialRedisClient(t *testing.T) {
	listenAddr := ":5002"
	server := NewServer(Config{ListenAddr: listenAddr})

	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Millisecond * 400)

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost" + listenAddr,
	})

	testCases := map[string]string{
		"foo":  "bar",
		"a":    "gg",
		"your": "mom",
		"step": "daddy",
	}

	for key, val := range testCases {
		if err := rdb.Set(context.Background(), key, val, 0).Err(); err != nil {
			t.Fatal(err)
		}
		newVal, err := rdb.Get(context.Background(), key).Result()
		if err != nil {
			t.Fatal(err)
		}
		if newVal != val {
			t.Fatalf("expected %s but got %s", val, newVal)
		}
	}
}
