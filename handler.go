package node2go

import (
	"encoding/json"
	"reflect"
)

// HandlerFunc for users to implement
type HandlerFunc func(interface{}) (interface{}, error)

// FuncMap to be passed o Runner constructor by users
type FuncMap map[string]*Handler

// Handler holds the func that will be called for specified endpoint
// DataTemplate may be any struct and will be used as a template when decoding
// json.
// If Raw is true, Func fill receive raw byte array instead of parsed data
type Handler struct {
	DataTemplate interface{}
	Func         HandlerFunc
	Raw          bool
}

// parse creates a new instance of DataTemplate and decodes json into it
func (h *Handler) parse(in []byte) (out interface{}, err error) {
	out = reflect.New(reflect.TypeOf(h.DataTemplate).Elem()).Interface()
	err = json.Unmarshal(in, &out)
	return
}
