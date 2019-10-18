package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	reuse "github.com/jbenet/go-reuseport"
)

const HAND_SHAKE_MSG = "我是来打洞的"

//const SERVER_IP = "49.233.79.223"
const SERVER_IP = "127.0.0.1"

// const SERVER_IP = "192.168.10.173"

func main() {
	name := flag.String("name", "tx", "name of client")
	clientPort := flag.Int("port", 9982, "port")
	flag.Parse()
	srcAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: *clientPort,
	}
	// l1, _ := reuse.Listen("tcp", "127.0.0.1:1234")
	/* srcAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", net.IPv4zero.String(), *clientPort))
	if err != nil {
		fmt.Println("Can't resolve address:", err)
		return
	} */
	dstAddr := &net.UDPAddr{
		IP:   net.ParseIP(SERVER_IP),
		Port: 9981,
	}
	// d := net.Dialer{Timeout: time.Duration(time.Second * 5), LocalAddr: srcAddr}

	srcStr := fmt.Sprintf("%s:%d", net.IPv4zero.String(), *clientPort)
	dstStr := fmt.Sprintf("%s:%d", SERVER_IP, 9981)
	conn, err := reuse.Dial("udp", srcStr, dstStr)
	// conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err = conn.Write([]byte("hello, I'm new peer:" + *name)); err != nil {
		log.Panic(err)
	}
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Printf("error during read: %s\r\n", err.Error())
		return
	}

	otherPeer := parseAddr(string(data[:n]))
	fmt.Printf("local: %s, server: %s, other: %s\r\n", srcAddr.String(), remoteAddr.String(), otherPeer.String())

	srcAddr1, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", net.IPv4zero.String(), *clientPort))
	if err != nil {
		fmt.Println("Can't resolve address:", err)
		return
	}
	bidirectionHole(srcAddr1, otherPeer, *name)
	conn.Close()
}

func parseAddr(addr string) *net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return &net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

func bidirectionHole(srcAddr *net.UDPAddr, otherAddr *net.UDPAddr, name string) {
	conn, err := net.DialUDP("udp", srcAddr, otherAddr)
	if err != nil {
		fmt.Printf("make hole failed(dialudp): %s\r\n", err.Error())
		return
	}
	defer conn.Close()
	if _, err = conn.Write([]byte(HAND_SHAKE_MSG)); err != nil {
		log.Println("send handshake:", err)
	}

	go func() {
		for {
			time.Sleep(2 * time.Second)
			if _, err = conn.Write([]byte("from [" + name + "]")); err != nil {
				log.Println("send message failed", err)
			}
		}
	}()

	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\r\n", err)
		} else {
			log.Printf("recv data:%s\r\n", data[:n])
		}
	}
}
