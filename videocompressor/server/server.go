package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"compressor/protocol"
)

// 終わったらリファクタリングする

func main() {
	// 通信を確立
	psock, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := psock.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) error {
	defer conn.Close()

	const MAX_PACKET_SIZE int = 1400
	var fileLen uint32
	var packetSum uint32
	var statusCode uint16
	buf := bufio.NewReaderSize(conn, MAX_PACKET_SIZE)

	file, err := os.Create("video.mp4") // 名前は動的生成するようにする
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("server launched")

	// ファイルのバイト数を受け取る
	// 最初のパケットの最初の4バイトからバイト数を確認
	filelenBytes := make([]byte, 4) // FIXME: マジックナンバーが出てきてるのでvidp型に抽象化したい

	_, err = buf.Read(filelenBytes)
	if err != nil {
		return err
	}
	fileLen = binary.BigEndian.Uint32(filelenBytes)

	for {
		if packetSum > fileLen {
			statusCode = protocol.StatusOK
			break
		}

		// ファイルに逐次書き込み
		fileBuf := []byte{}
		buf.Read(fileBuf)
		n, err := writeToFile(file, fileBuf)

		if err != nil {
			// エラーコード
			statusCode = protocol.StatusNG
			break
		}
		packetSum += uint32(n)
	}

	if statusCode == 1 {
		err := os.Remove("video.mp4") // ハードコード
		if err != nil {
			return err
		}
	}
	// ステータスメッセージを返す
	_, err = conn.Write([]byte{byte(statusCode)})
	if err != nil {
		return err
	}
	return nil
}

func writeToFile(file *os.File, bytes []byte) (int, error) {
	n, err := file.Write(bytes)
	if err != nil {
		return 0, err
	}
	return n, nil
}
