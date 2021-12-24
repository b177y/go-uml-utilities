package mconsole

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func min(x uint32, y uint32) uint32 {
	if x < y {
		return x
	}
	return y
}

func recvOutput(sock net.UnixConn) (output string, err error) {
	var reply mconsoleReply
	reply.More = 1
	for reply.More == 1 {
		respBytes := make([]byte, MCONSOLE_MAX_DATA+12)
		_, err := sock.Read(respBytes)
		if err != nil {
			return "", err
		}
		err = binary.Read(bytes.NewBuffer(respBytes), binary.LittleEndian, &reply)
		if err != nil {
			return "", err
		}
		if reply.Err != 0 {
			return "", fmt.Errorf("Error from mconsole: %d", reply.Err)
		}
		output += string(reply.Data[:])
	}
	return output, err
}

func SendCommand(command string, sock net.UnixConn) (output string, err error) {
	req := mconsoleRequest{
		magic:   MCONSOLE_MAGIC,
		version: MCONSOLE_VERSION,
		length:  min(uint32(len(command)), MCONSOLE_MAX_DATA),
	}
	copy(req.data[:], []byte(command)[:req.length])
	req.data[req.length] = byte('\x00')
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, &req)
	if err != nil {
		return "", err
	}
	_, err = sock.Write(buf.Bytes())
	if err != nil {
		return "", err
	}
	return recvOutput(sock)
}

func SendCommandToSock(command string,
	sockpath string) (output string, err error) {
	ra, err := net.ResolveUnixAddr("unixgram", sockpath)
	if err != nil {
		return "", err
	}
	la, err := net.ResolveUnixAddr("unixgram", "@"+fmt.Sprint(os.Getpid())+"@@@@")
	if err != nil {
		return "", err
	}
	conn, err := net.DialUnix("unixgram", la, ra)
	if err != nil {
		return "", err
	}
	return SendCommand(command, *conn)
}
