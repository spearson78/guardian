package tcpaddr

import (
	"encoding/binary"
	"net"
)

const Size = 16 + 2

func PutTCPAddr(buf []byte, value *net.TCPAddr) int {

	copy(buf[:16], value.IP.To16())
	binary.BigEndian.PutUint16(buf[16:18], uint16(value.Port))

	return Size
}

func TCPAddr(buf []byte) (net.TCPAddr, int) {

	var value net.TCPAddr
	value.IP = make([]byte, 16)
	copy(value.IP, buf[:16])
	value.Port = int(binary.BigEndian.Uint16(buf[16:18]))

	return value, Size

}
