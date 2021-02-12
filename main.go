package main

import (
	"context"
	"log"
	"time"
)

func main() {
	go func() {
		log.Print("starting socks server...")
		err := RunSocksServer()
		if err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()
	for {
		err := runMotionPro(ctx)
		if err != nil {
			log.Printf("Motion Pro exited with error: %v", err.Error())
		} else {
			log.Print("Motion Pro exited")
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
