package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

var vpnHost = flag.String("vpn-host", "vpn-cn.ecnu.edu.cn", "hostname of vpn")
var username = flag.String("vpn-username", "", "username of vpn")
var password = flag.String("vpn-password", "", "password of vpn")

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

// StartVpn start ecnu vpn
func StartVPN(ctx context.Context) (err error) {
	// Start openconnect
	log.Print("Exec openconnect")
	c := exec.Command("openconnect", "-u", *username, "--passwd-on-stdin", *vpnHost)
	f, err := pty.Start(c)
	if err != nil {
		return
	}
	defer func() {
		if c.ProcessState != nil {
			c.Process.Kill()
		}
		log.Println("Openconnect stopped")
	}()
	// Gracefully kill cmd when ctx done
	go func() {
		<-ctx.Done()
		c.Process.Kill()
	}()

	// Redirect cmd output
	go io.Copy(os.Stdout, f)

	// Apply password
	_, err = f.WriteString(*password + "\n")
	if err != nil {
		return fmt.Errorf("failed to apply password %w", err)
	}

	return c.Wait()
}
