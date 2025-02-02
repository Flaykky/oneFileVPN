package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"
)

func main() {
	proxy := flag.String("proxy", "", "Proxy address in format host:port")
	credentials := flag.String("auth", "", "Authentication in format user:password")
	flag.Parse()

	if *proxy == "" {
		fmt.Println("Proxy address is required")
		os.Exit(1)
	}

	parsedProxy, err := url.Parse(fmt.Sprintf("http://%s", *proxy))
	if err != nil {
		fmt.Println("Invalid proxy address:", err)
		os.Exit(1)
	}

	if *credentials != "" {
		auth := strings.Split(*credentials, ":")
		if len(auth) != 2 {
			fmt.Println("Invalid credentials format")
			os.Exit(1)
		}
		parsedProxy.User = url.UserPassword(auth[0], auth[1])
	}

	proxyURL := parsedProxy.String()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("Error starting local listener:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("Local SOCKS5 proxy started on %s\n", listener.Addr())

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(clientConn, proxyURL)
	}
}

func handleClient(clientConn net.Conn, proxyURL string) {
	defer clientConn.Close()

	targetConn, err := net.Dial("tcp", proxyURL)
	if err != nil {
		fmt.Println("Error connecting to proxy:", err)
		return
	}
	defer targetConn.Close()

	go io.Copy(clientConn, targetConn)
	io.Copy(targetConn, clientConn)
}
