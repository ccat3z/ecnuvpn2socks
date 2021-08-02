package main

import (
	"context"
	"flag"
	"log"
	"time"
)

var (
	immediateFailureThresholdSecond = flag.Int("fail-threshold", 60, "Threshold of immediate failure time in second")
	maxFailureTime                  = flag.Int("max-try", 3, "Maximum attempts of reconnect")
)

func main() {
	// Socks server
	go func() {
		log.Print("Starting socks server...")
		err := RunSocksServer()
		if err != nil {
			panic(err)
		}
	}()

	// VPN
	lastFail := time.Unix(0, 0)
	failTimes := 0

	ctx := context.Background()
	for {
		err := StartVPN(ctx)
		if err != nil {
			log.Printf("VPN exited with error: %v", err.Error())
		} else {
			log.Print("VPN exited")
		}

		if time.Since(lastFail) < time.Duration(*immediateFailureThresholdSecond)*time.Second {
			failTimes++
		} else {
			failTimes = 1
		}
		lastFail = time.Now()

		if failTimes >= *maxFailureTime {
			log.Printf("Max retry reached (%v)", failTimes)
			break
		}

		time.Sleep(2 * time.Second)
	}
}
