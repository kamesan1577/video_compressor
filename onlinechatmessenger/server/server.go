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

	clients := map[net.Addr]client{}

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
		// FIXME: 一覧にあるはずなのに毎回アドレスが追加される
		_, ok := clients[addr]
		if !ok {
			fmt.Println("new addr added: ", addr.String())
			clients[addr] = client{time.Now(), 0}
		}

		for k, v := range clients {
			// しばらく送信がないか、連続で失敗した場合、クライアント一覧から消す
			// 最終更新時間が現在よりもtimeout秒以上前なら削除
			subSec := time.Now().Sub(v.LastMessageAt).Seconds()
			timeout := 60
			if subSec >= float64(timeout) || v.ErrorCount >= 3 {
				delete(clients, k)
				fmt.Println("client deleted")
			}
		}

		// 接続中の送信者以外のクライアントにリレーする
		// FIXME: 同じアドレスが複数存在していてメッセージ送信時に自分のメッセージが大量に届く
		for k, v := range clients {
			// FIXME: 送信者のアドレスにもメッセージが送信されている
			if v == clients[addr] {
				continue
			}
			//kのアドレスにそうしんする
			_, err = conn.WriteToUDP(append([]byte(fmt.Sprintf("%v: ", username)), []byte(message)...), k.(*net.UDPAddr))
			if err != nil {
				fmt.Println(err)
				v = client{v.LastMessageAt, v.ErrorCount + 1}
				return
			}
			v = client{time.Now(), 0}
			fmt.Println("message sent to: ", k)
		}

	}
}

type client struct {
	LastMessageAt time.Time
	ErrorCount    int
}
