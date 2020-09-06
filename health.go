package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

const testURL = "http://g.cn/generate_204"
const testStatus = 204
const testInterval = 15 * time.Second

const maxRetry = 3
const retryInterval = 1 * time.Second

func init() {
	http.DefaultClient.Timeout = 10 * time.Second
}

func IsHealth() bool {
	for i := 1; i <= maxRetry; i++ {
		log.Printf("check health (%v/%v)", i, maxRetry)
		if Reachable(testURL, testStatus) {
			return true
		}
		time.Sleep(retryInterval)
	}

	return false
}

func Reachable(url string, status int) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	return resp.StatusCode == status
}

func WaitUnhealth(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		case <-time.NewTimer(testInterval).C:
		}

		if !IsHealth() {
			break
		}
	}
	return nil
}
