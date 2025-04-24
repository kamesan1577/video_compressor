// package main

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"fmt"
// 	"net"
// 	"os"

// 	"compressor/protocol"
// )

// // 終わったらリファクタリングする

// func main() {
// 	// 通信を確立
// 	psock, err := net.Listen("tcp", ":9000")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// ch := make(chan error)

// 	for {
// 		conn, err := psock.Accept()
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		go handleConnection(conn)
// 	}
// }

// func handleConnection(conn net.Conn) error {
// 	defer conn.Close()

// 	var fileLen uint32
// 	var packetSum uint32
// 	var statusCode uint16
// 	// buf := bufio.NewReaderSize(conn, protocol.MAX_PACKET_SIZE)
// 	buf := new(bytes.Buffer)

// 	file, err := os.Create("video.txt") // 名前は動的生成するようにする
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	fmt.Println("server launched")

// 	// ファイルのバイト数を受け取る
// 	// 最初のパケットの最初の4バイトからバイト数を確認
// 	filelenBytes := make([]byte, 4) // FIXME: マジックナンバーが出てきてるのでvidp型に抽象化したい

// 	_, err = conn.Read(filelenBytes)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fileLen = binary.BigEndian.Uint32(filelenBytes)

// 	for {
// 		if packetSum > fileLen {
// 			statusCode = protocol.StatusOK
// 			break
// 		}

// 		// ファイルに逐次書き込み
// 		_, err = conn.Read(buf.Bytes())
// 		if err != nil {
// 			panic(err)
// 		}
// 		n, err := file.Write(buf.Bytes())
// 		if err != nil {
// 			// エラーコード
// 			statusCode = protocol.StatusNG
// 			_err := sendStatus(statusCode, conn)
// 			if _err != nil {
// 				panic(_err)
// 			}
// 			panic(err)
// 		}
// 		fmt.Println(buf.String())
// 		// fmt.Println(packetSum)
// 		packetSum += uint32(n)
// 	}

// 	if statusCode == 1 {
// 		err := os.Remove("video.txt") // ハードコード
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	// ステータスメッセージを返す
// 	err = sendStatus(statusCode, conn)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return nil
// }

// func writeToFile(file *os.File, bytes []byte) (int, error) {
// 	n, err := file.Write(bytes)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return n, nil
// }

// func sendStatus(statusCode uint16, conn net.Conn) error {
// 	_, err := conn.Write([]byte{byte(statusCode)})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
