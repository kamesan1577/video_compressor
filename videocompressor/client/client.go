package main

import (
	// "bufio"
	// "bufio"
	"compressor/protocol"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	// "io"
	"net"
	"os"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

// 全体の流れ
// サーバーとのセッション開始
// パケットのサイズ分のバッファを作る
// 最初のパケットの頭4バイトにはデータサイズを埋め込む
// for ステータスコードが返ってくる || タイムアウトするまで
// ファイルの中身をバッファの空き分取得する
// バッファの空きにファイルの中身を詰める
// バッファの中身を送信する
// バッファをクリアする
// forend
// セッションを閉じる
// ステータスコードを返す
func (c *Client) SendFile(f *os.File) (string, error) {
	var fileLenCount uint32
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		fmt.Println("net.dial tcp")
		return "", err
	}
	defer conn.Close()

	// buf := bufio.NewReader(f)

	// ファイルをパケットに分けて送信する
	// 最初のパケットにはデータ長を表すバイトを入れる
	buf := make([]byte, protocol.MAX_PACKET_SIZE)

	fileStat, err := f.Stat()
	if err != nil {
		fmt.Println("file.stat")
		return "", err
	}
	filelen := uint32(fileStat.Size())
	fileLenBytes := make([]byte, binary.MaxVarintLen32)
	binary.BigEndian.PutUint32(fileLenBytes, filelen)
	//bufの最初の4バイトにfileLenBytesを入れたい
	buf = append(buf, fileLenBytes...)

	n, err := f.Read(buf)
	if err != nil {
		fmt.Println("file.read")
		return "", err // FIXME: ここで空文字列を返すのはおかしいから本当はnilを返したい
	}
	fileLenCount += uint32(n)

	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println("conn.write")
		return "", err
	}

	// ファイルの最後まで送信する
	i := 0
	for {
		fmt.Println(i)
		n, err = readFromFile(f, &buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("readFromFile")
			return "", err
		}
		_, err = conn.Write(buf)
		if err != nil {
			fmt.Println("conn.Write2")
			if strings.Contains(err.Error(), "connection reset") {
				break
			}
			return "", err
		}
		// fmt.Println(buf)
		fileLenCount += uint32(n)
		i++
	}
	fmt.Println("finished")
	// ステータスコードが返ってきたら終了
	var resp []byte
	_, err = conn.Read(resp)
	if err != nil {
		fmt.Println("status read ")
		return "", err
	}
	status := string(resp)
	return status, nil
}

func readFromFile(file *os.File, out *[]byte) (int, error) {
	n, err := file.Read(*out)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// func sendPacket(buf []byte, conn net.Conn) error {
// 	_, err := conn.Write(buf)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func main() {
	// ファイルパスを取得
	if len(os.Args) != 2 {
		fmt.Println("usage hogehoge")
		return
	}
	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// mp4かどうか確認する

	c := NewClient(":9000")
	resp, err := c.SendFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// respのステータスコードによって出力を変える
	fmt.Println(resp)

}

// パケットのサイズ分のバッファを作る
// バッファにパケットサイズ分のファイルデータを詰める
// バッファを相手に送る
// バッファをクリアする
