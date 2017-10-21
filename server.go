package node2go

import (
	"bufio"
	"bytes"
	"net"
)

type socket struct {
	net.Listener
}

func (s *socket) read() (c chan *token, e chan *Err) {
	c, e = make(chan *token), make(chan *Err)

	go func() {
		for {
			conn, err := s.Listener.Accept()
			if err != nil {
				e <- newErr(conn, nilID, err)
				continue
			}

			s := bufio.NewScanner(conn)
			s.Split(bufio.ScanLines)

			for s.Scan() {
				if err := s.Err(); err != nil {
					e <- newErr(conn, nilID, err)
					continue
				}

				split := bytes.Split(s.Bytes(), []byte(";"))
				if len(split) != 3 {
					e <- newErr(conn, nilID, ErrMalformedMessage)
					continue
				}

				if err != nil {
					e <- newErr(conn, nilID, err)
				}

				c <- &token{
					id:       split[0],
					Conn:     conn,
					FuncName: string(split[1]),
					Data:     split[2],
				}
			}
		}
	}()

	return

}

func newSocket(addr string) (s *socket, err error) {
	l, err := net.Listen("unix", addr)
	if err != nil {
		return
	}
	s = &socket{l}
	return
}
