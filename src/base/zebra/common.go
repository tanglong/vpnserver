package zebra

import (
	"errors"
	"time"
)

const (
	PACK_HEAD_LEN         = 20
	READ_TIME_OUT         = 60 * time.Second
	WRITE_TIME_OUT        = 10 * time.Second
	HIGH_WATER_MARK_SCALE = 0.9
)

var (
	NoConnect     = errors.New("tcp_conn: no connect")
	ReadOverflow  = errors.New("tcp_conn: read buffer overflow")
	WriteOverflow = errors.New("tcp_conn: write buffer overflow")
	ErrorMsgType  = errors.New("tcp_conn: error msg type")
)

//Big Endian
func DecodeUint32(data []byte) uint32 {
	return (uint32(data[0]) << 24) | (uint32(data[1]) << 16) | (uint32(data[2]) << 8) | uint32(data[3])
}

//Big Endian
func DecodeUint64(data []byte) uint64 {
	return (uint64(data[0]) << 56) | (uint64(data[1]) << 48) | (uint64(data[2]) << 40) | (uint64(data[3]) << 32) | (uint64(data[4]) << 24) | (uint64(data[5]) << 16) | (uint64(data[6]) << 8) | uint64(data[7])
}

//Big Endian
func EncodeUint32(n uint32, b []byte) {
	b[3] = byte(n & 0xFF)
	b[2] = byte((n >> 8) & 0xFF)
	b[1] = byte((n >> 16) & 0xFF)
	b[0] = byte((n >> 24) & 0xFF)
}

//Big Endian
func EncodeUint64(n uint64, b []byte) {
	b[7] = byte(n & 0xff)
	b[6] = byte((n >> 8) & 0xFF)
	b[5] = byte((n >> 16) & 0xFF)
	b[4] = byte((n >> 24) & 0xFF)
	b[3] = byte((n >> 32) & 0xFF)
	b[2] = byte((n >> 40) & 0xFF)
	b[1] = byte((n >> 48) & 0xFF)
	b[0] = byte((n >> 56) & 0xFF)
}
