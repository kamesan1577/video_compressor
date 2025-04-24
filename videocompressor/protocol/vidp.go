package protocol

import (
	"bytes"
	"encoding/binary"
)

type Vidp struct {
	FileLen uint32 // 仕様書通りならuint32を使うべき
	Data    []byte
}

func (v *Vidp) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write([]byte{byte(v.FileLen)})
	buf.Write(v.Data)
	return buf.Bytes()

}

func NewVidp(filelen uint32, data []byte) Vidp {
	vidp := Vidp{filelen, data}
	return vidp
}

func ParseVidp(bytes []byte) Vidp {
	filelen := binary.BigEndian.Uint32(bytes[:3])
	data := bytes[4:]
	return NewVidp(filelen, data)
}

const (
	StatusOK = 0
	StatusNG = 1
)

const (
	MAX_PACKET_SIZE int = 1400
)
