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

var motionProLock sync.Mutex

// StartMotionPro start motion pro client in subprocess,
// and wait motion pro ready. Then do some extract operation to make vpn work.
func StartMotionPro(ctx context.Context, timeout time.Duration) (err error) {
	motionProLock.Lock()
	log.Print("Starting Motion Pro")

	// build pty bash process
	c := exec.Command("bash")
	f, err := pty.Start(c)
	if err != nil {
		return
	}

	kill := func() {
		if c.ProcessState != nil {
			return
		}

		log.Print("Killing motion pro")
		f.Write([]byte{3})
		f.WriteString("exit\n")
		f.Write([]byte{4})
	}

	// redirect cmd output
	go io.Copy(os.Stdout, f)

	cmdExited := make(chan struct{})

	// release resources after cmd finished
	go func() {
		c.Wait()
		close(cmdExited)
		motionProLock.Unlock()
		log.Print("Motion Pro was stopped")
	}()

	// gracefully kill cmd when ctx done
	go func() {
		select {
		case <-cmdExited:
			break
		case <-ctx.Done():
			kill()
		}
	}()

	// kill cmd if prepare failed
	defer func() {
		if err != nil {
			kill()
		}
	}()

	// prepare operations
	prepareCtx, cancelPrepare := context.WithCancel(context.Background())
	defer cancelPrepare()
	go func() {
		<-cmdExited
		cancelPrepare()
	}()

	// send commands into bash
	f.WriteString("export PS1='# '\n")
	f.WriteString("rm -v /etc/*.array\n")
	f.WriteString(fmt.Sprintf("cd '%v' || exit 1\n", *motionProLib))
	f.WriteString(fmt.Sprintf("./vpn_cmdline -h '%v' -o '%v' -u '%v' -p '%v'; exit\n", *vpnHost, *vpnPort, *username, *password))

	err = waitMotionProReady(prepareCtx)
	if err != nil {
		err = fmt.Errorf("Failed to start motion pro > %w", err)
		return
	}

	err = fixSystemConfigAfterMotionProStarted()
	if err != nil {
		return
	}

	return nil
}

// WaitMotionPro wait current motion pro instance exit.
// If no motion pro is running, it will return immedately.
func WaitMotionPro() {
	motionProLock.Lock()
	motionProLock.Unlock()
}

func waitMotionProReady(ctx context.Context) error {
	for {
		log.Println("waiting motionpro...")

		select {
		case <-ctx.Done():
			return ctx.Err()
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

	log.Println("motionpro ready!")
	return nil
}

func fixSystemConfigAfterMotionProStarted() error {
	log.Println("fix resolv.conf")

	resolv, err := ioutil.ReadFile("/etc/resolv.conf.array")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/etc/resolv.conf", resolv, 0644)
	if err != nil {
		return err
	}

	return nil
}
