package main

import (
	"bufio"
	"bytes"
	// "encoding/binary"
	"fmt"
	"net"
	"os"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Type server address: ")
	scanner.Scan()
	serverHost := scanner.Text()
	serverPort := "9001"
	serverAddr, err := net.ResolveUDPAddr("udp", serverHost+":"+serverPort)
	if err != nil {
		fmt.Println(err)
		return
	}

	host := ""
	fmt.Print("Type client port: ")
	scanner.Scan()
	port := scanner.Text()
	if port == "" {
		port = "9050"
	}
	addr, err := net.ResolveUDPAddr("udp", host+":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("Type your name: ")
	scanner.Scan()
	username := scanner.Text()
	usernamelen := len(username)

	conn, err := net.DialUDP("udp", addr, serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	//応答を受信
	go readMsg(conn)

	for {
		fmt.Print("Type message to send: ")
		// FIXME: 入力待ちの時にメッセージを受信するとその文字列が入力されてしまう
		// 		  本当に標準入力に出力されているのか、表示上の問題なのかを切り分ける
		//		  ←表示上の問題っぽい
		//		  ←type message to sendを消せば見た目上はおかしくならなそう
		scanner.Scan()
		message := scanner.Text()

		if message == "!exit" {
			return
		}

		fmt.Printf("Sending %v\n", message)

		// usernamelen: int + username: string + message: stringのバイナリ
		fullMessageBuf := bytes.NewBuffer([]byte{})
		fullMessageBuf.Write([]byte{byte(usernamelen)})
		fullMessageBuf.Write([]byte(username))
		fullMessageBuf.Write([]byte(message))

		// connのアドレスに送信する
		_, err = conn.Write(fullMessageBuf.Bytes())
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func readMsg(conn *net.UDPConn) {
	for {

		readBuf := make([]byte, 4096)
		n, _, err := conn.ReadFromUDP(readBuf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(string(readBuf[:n]))
		// fmt.Fprintln(os.Stderr, string(readBuf[:n]))
	}
}
