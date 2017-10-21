package node2go

import "net"

type token struct {
	net.Conn
	FuncName string
	Data     []byte
	id       []byte
}
