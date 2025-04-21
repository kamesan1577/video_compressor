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
		fmt.Println("usernamelen: ", usernamelen)
		username := string(buffer[1 : usernamelen+1])
		fmt.Println("username: ", username)

		message := buffer[1+usernamelen : n]
		fmt.Println("message: ", string(message))

		// クライアント一覧になければ追加
		_, ok := clients[addr]
		if !ok {
			clients[addr] = client{time.Now(), 0}
		}

		for k, v := range clients {
			// しばらく送信がないか、連続で失敗した場合、クライアント一覧から消す
			// 最終更新時間が現在よりも10秒以上前なら削除
			subSec := time.Now().Sub(v.LastMessageAt).Seconds()
			if subSec >= 10 || v.ErrorCount >= 3 {
				delete(clients, k)
				fmt.Println("client deleted")
			}
		}

		// 接続中のクライアントにリレーする
		for k, v := range clients {
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
