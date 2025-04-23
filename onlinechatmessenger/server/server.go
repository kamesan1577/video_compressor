package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:9001")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	clients := map[string]client{}

	for {
		fmt.Println("\nwaiting to receive message")
		buffer := make([]byte, 4096)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("received %v bytes from %v\n", len(buffer[:n]), addr)

		// メッセージの最初の1バイトからユーザー名を特定
		usernamelen := buffer[0]
		if usernamelen <= 0 {
			fmt.Println("Invalid bytes received: usernamelen must be greater than 0")
			continue
		}
		fmt.Println("usernamelen: ", usernamelen)
		username := string(buffer[1 : usernamelen+1])
		fmt.Println("username: ", username)

		message := buffer[1+usernamelen : n]
		fmt.Println("message: ", string(message))

		// クライアント一覧になければ追加
		_, ok := clients[addr.String()]
		if !ok {
			fmt.Println("new addr added: ", addr.String())
			clients[addr.String()] = client{addr, time.Now(), 0}

			enterLogBuf := []byte(username + " entered")
			broadCast(enterLogBuf, clients, conn, *addr)
		}

		for k, v := range clients {
			// しばらく送信がないか、連続で失敗した場合、クライアント一覧から消す
			// 最終更新時間が現在よりもtimeout秒以上前なら削除
			subSec := time.Now().Sub(v.LastMessageAt).Seconds()
			timeout := 60
			if subSec >= float64(timeout) || v.ErrorCount >= 3 {
				exitLogBuf := []byte(username + " exited")
				broadCast(exitLogBuf, clients, conn, *addr)
				delete(clients, k)
				fmt.Println("client deleted")
			}
		}
		messageBuf := append([]byte(fmt.Sprintf("%v: ", username)), []byte(message)...)
		err = broadCast(messageBuf, clients, conn, *addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(clients)

		clients[addr.String()] = client{addr, time.Now(), 0}
	}
}

type client struct {
	Address       *net.UDPAddr
	LastMessageAt time.Time
	ErrorCount    int
}

func broadCast(messageBuf []byte, clients map[string]client, conn *net.UDPConn, addr net.UDPAddr) error {
	// 接続中の送信者以外のクライアントにリレーする
	for k, v := range clients {
		if k == addr.String() {
			continue
		}
		//kのアドレスにそうしんする
		_, err := conn.WriteToUDP(messageBuf, v.Address)
		if err != nil {
			v = client{v.Address, v.LastMessageAt, v.ErrorCount + 1}
			return err
		}
		fmt.Println("message sent to: ", k)
	}
	return nil
}
