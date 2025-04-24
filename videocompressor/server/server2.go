package main

import (
	"bufio"
	"compressor/protocol"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func receiveTCPConn(ln *net.TCPListener) {
	for {
		err := ln.SetDeadline(time.Now().Add(time.Second * 60))
		if err != nil {
			log.Fatal(err)
		}
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go echoHandler(conn)
		go sendHandler(conn)
	}
}

func echoHandler(conn *net.TCPConn) {
	defer conn.Close()
	for {
		_, err := io.WriteString(conn, "Socket Connection!!\n")
		if err != nil {
			return
		}
		time.Sleep(time.Second)
	}
}

func sendHandler(conn *net.TCPConn) {
	defer conn.Close()
	response := protocol.StatusOK
	file, err := os.Create(time.Now().String() + ".mp4")
	if err != nil {
		panic(err)
	}
	for {
		err = saveToFile(file, conn)
		if err != nil {
			response = protocol.StatusNG
		}
		//処理が終わったらレスポンスする
		buf := bufio.NewWriter(conn)
		buf.WriteByte(byte(response))
	}
}

func saveToFile(file *os.File, src io.Reader) error {
	_, err := io.Copy(file, src)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	receiveTCPConn(ln)
}
