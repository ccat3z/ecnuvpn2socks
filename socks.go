package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/txthinking/socks5"
)

var socksPort = flag.Int("socks-port", 1080, "socks server port")
var ip = flag.String("ip", "127.0.0.1", "socks server ip")
var socksUsername = flag.String("socks-username", "", "username of socks account")
var socksPassword = flag.String("socks-password", "", "password of socks account")

func RunSocksServer() error {
	log.Printf("start socks server on :%v", *socksPort)
	addr := fmt.Sprintf("0.0.0.0:%v", *socksPort)
	s, err := socks5.NewClassicServer(addr, *ip, *socksUsername, *socksPassword, 0, 60)
	if err != nil {
		return err
	}
	s.ListenAndServe(nil)
	return nil
}
