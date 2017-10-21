package node2go

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
)

// General errors
var (
	ErrMalformedMessage = errors.New(
		"message must constist of a function name to be called, colon and json body",
	)
	ErrFuncNotImplemented = func(name string) error {
		return fmt.Errorf("handler for `%s` not implemented", name)
	}
)

// Err is a wrapper for error, has an id of incoming message and embeds
// the connection it originated in
type Err struct {
	net.Conn `json:"-"`
	id       []byte
	error    error
}

// MarshalJSON implements the marshaller interface and returns json with a
// single error field
func (e *Err) MarshalJSON() ([]byte, error) {
	type Alias Err
	return json.Marshal(&struct {
		Error string `json:"error"`
		*Alias
	}{
		Error: e.error.Error(),
		Alias: (*Alias)(e),
	})
}

func (e *Err) json() []byte {
	j, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}
	return formatResponse(j, e.id)
}

func newErr(conn net.Conn, id []byte, err error) *Err {
	return &Err{conn, id, err}
}
