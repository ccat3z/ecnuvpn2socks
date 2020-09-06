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

	for {
		err := RunMotionProWithHelper(context.Background())
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
	}
}

func RunMotionProWithHelper(ctx context.Context) error {
	r := make(chan struct{})

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		select {
		case <-r:
			err := WaitUnhealth(ctx)

			if err == context.Canceled {
				return
			}

			if err != nil {
				panic(err)
			}

			log.Print("cannot connect to internet")
			KillMotionPro()
		case <-ctx.Done():
		}
	}()

	return RunMotionPro(r)
}
