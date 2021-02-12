package main

import (
	"context"
	"flag"
	"log"
	"time"
)

var (
	prepareTimeout = 30 * time.Second
	// TODO: immediateFailureThresholdSecond should greater than prepareTimeout
	immediateFailureThresholdSecond = flag.Int("fail-threshold", 60, "Threshold of immediate failure time in second")
	maxFailureTime                  = flag.Int("max-try", 3, "Maximum attempts of reconnect")
)

func main() {
	go func() {
		log.Print("starting socks server...")
		err := RunSocksServer()
		if err != nil {
			panic(err)
		}
	}()

	lastFail := time.Unix(0, 0)
	failTimes := 0

	ctx := context.Background()
	for {
		err := runMotionPro(ctx)
		if err != nil {
			log.Printf("Motion Pro exited with error: %v", err.Error())
		} else {
			log.Print("Motion Pro exited")
		}

		if time.Now().Sub(lastFail) < time.Duration(*immediateFailureThresholdSecond)*time.Second {
			failTimes++
		} else {
			failTimes = 1
		}
		lastFail = time.Now()

		if failTimes >= maxRetry {
			log.Printf("Max retry reached (%v)", failTimes)
			break
		}

		time.Sleep(2 * time.Second)
	}
}

func runMotionPro(ctx context.Context) error {
	motionCtx, cancelMotion := context.WithCancel(ctx)
	defer cancelMotion()

	unhealthMonCtx, cancelUnhealthMon := context.WithCancel(ctx)
	defer cancelUnhealthMon()

	err := StartMotionPro(motionCtx, time.Second*30)
	defer WaitMotionPro()
	if err != nil {
		return err
	}

	go func() {
		log.Print("Start health montor")
		err := WaitUnhealth(unhealthMonCtx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			panic(err) // unexcept error
		}
		log.Print("Cannot connect to internet")
		cancelMotion()
	}()

	return nil
}
