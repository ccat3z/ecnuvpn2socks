package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/creack/pty"
)

var vpnHost = flag.String("vpn-host", "vpn.shnu.edu.cn", "hostname of vpn")
var vpnPort = flag.String("vpn-port", "443", "port of vpn")
var username = flag.String("vpn-username", "", "username of vpn")
var password = flag.String("vpn-password", "", "password of vpn")
var motionProLib = flag.String("motionpro-lib", "/usr/local/share/motionpro/", "motionpro lib path")

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if *username == "" {
		log.Fatal("-vpn-username is empty")
	}

	if *password == "" {
		log.Fatal("-vpn-password is empty")
	}
}

var motionProIO *os.File = nil

func RunMotionPro(ready chan<- struct{}) error {
	defer close(ready)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	c := exec.Command("bash")
	f, err := pty.Start(c)

	motionProIO = f
	defer func() {
		motionProIO = nil
	}()

	if err != nil {
		return err
	}

	// write commands into bash
	go func() {
		wg.Add(1)
		defer wg.Done()
		f.WriteString("export PS1='# '\n")
		f.WriteString("rm -v /etc/*.array\n")
		f.WriteString(fmt.Sprintf("cd '%v' || exit 1\n", *motionProLib))
		f.WriteString(fmt.Sprintf("./vpn_cmdline -h '%v' -o '%v' -u '%v' -p '%v'; exit\n", *vpnHost, *vpnPort, *username, *password))
	}()

	// wait motion pro ready
	waitingCtx, cancelWait := context.WithCancel(context.Background())
	defer cancelWait()
	go func() {
		wg.Add(1)
		defer wg.Done()
		err := waitMotionProReady(waitingCtx)
		if err == context.Canceled {
			return
		}
		if err != nil {
			panic(err)
		}
		ready <- struct{}{}
	}()

	// redirect bash output
	io.Copy(os.Stdout, f)

	log.Print("waiting motion pro stop")
	c.Wait()
	log.Print("motion pro is stopped")
	return nil
}

func KillMotionPro() {
	log.Print("killing motion pro")
	if motionProIO == nil {
		log.Print("no running motion pro")
		return
	}

	motionProIO.Write([]byte{3})
	motionProIO.WriteString("exit\n")
	motionProIO.Write([]byte{4})
}

func waitMotionProReady(ctx context.Context) error {
	for {
		log.Println("waiting motionpro...")

		select {
		case <-ctx.Done():
			return context.Canceled
		case <-time.NewTimer(2 * time.Second).C:
		}

		_, err := os.Stat("/etc/resolv.conf.array")
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			break
		}
	}

	resolv, err := ioutil.ReadFile("/etc/resolv.conf.array")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/resolv.conf", resolv, 0644)
	if err != nil {
		return err
	}

	log.Println("motionpro ready!")

	return nil
}
