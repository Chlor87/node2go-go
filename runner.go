package node2go

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// ack is the message that'll be send to stdout for the NodeJS counterpart
// to act upon
var ack = "ACK"

// Runner embeds the socket and holds the user-provided function map
type Runner struct {
	*socket
	funcMap FuncMap
}

func (r *Runner) handleReads(tok *token) {
	h, ok := r.funcMap[tok.FuncName]
	if !ok {
		tok.Write(
			newErr(nil, tok.id, ErrFuncNotImplemented(tok.FuncName)).json(),
		)
		return
	}

	var d interface{}

	if h.Raw {
		d = tok.Data
	} else {
		var err error
		d, err = h.parse(tok.Data)
		if err != nil {
			tok.Write(newErr(nil, tok.id, err).json())
		}
	}

	res, err := h.Func(d)
	if err != nil {
		tok.Write(newErr(nil, tok.id, err).json())
		return
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		tok.Write(newErr(nil, tok.id, err).json())
		return
	}

	tok.Write(formatResponse(jsonRes, tok.id))

}

// Run reads data and errors from socket channels and either returns errors
// back to NodeJS or calls user-provided functions with parsed or raw data.
// Each function is called in a separate goroutine
// see handleReads
func (r *Runner) Run() {
	tokC, errC := r.read()

	for {
		select {
		case tok := <-tokC:
			go r.handleReads(tok)

		case err := <-errC:
			err.Conn.Write(err.json())
		}
	}
}

// Close waits for the signal terminating the process in order to gracefully
// shutdown removing the socket file.
func (r *Runner) Close() {
	defer r.Listener.Close()
	sigC := make(chan os.Signal)
	signal.Notify(sigC, os.Interrupt, os.Kill, syscall.SIGTERM)
	log.Println("Got signal", <-sigC)
}

// NewRunner constructs the runner with a unix domain socket
// (the address comes from `addr` flag)
// and sends the acknowledgement to NodeJS (or any other caller)
func NewRunner(funcMap FuncMap) (r *Runner, err error) {
	s, err := newSocket(socketAddr)
	if err != nil {
		return
	}
	fmt.Fprintln(os.Stdout, ack)
	return &Runner{s, funcMap}, nil
}
