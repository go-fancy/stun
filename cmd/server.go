package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 9981,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("local address: %s\r\n", listener.LocalAddr().String())
	peers := make([]net.UDPAddr, 0, 2)
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		log.Printf("remote address: %s, data: %s\r\n", remoteAddr.String(), data[:n])
		peers = append(peers, *remoteAddr)
		if len(peers) == 2 {
			remoteAddrStr1 := peers[0].String()
			remoteAddrStr2 := peers[1].String()
			log.Printf("connect %s <--> %s\r\n", remoteAddrStr1, remoteAddrStr2)
			listener.WriteToUDP([]byte(remoteAddrStr2), &peers[0])
			listener.WriteToUDP([]byte(remoteAddrStr1), &peers[1])
			time.Sleep(time.Minute * 3)
			log.Println("transit server exit")
			return
		}
	}
}
