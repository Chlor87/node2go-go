# Node2Go

## Description

This is the Go part of the library. NodeJS counterpart: [node2go-node](https://github.com/Chlor87/node2go-node).

Node2Go is an easy way to call predefined Go functions from NodeJS with json data going back and forth.
Im aware of the fact, that there's a library similar to this one on Github. It's called [gonode](https://github.com/jgranstrom/gonode).
I just wanted to write one myself, for fun mostly.

## Features
* no dependencies
* simple API - not yet documented in this file

To be continued...

## Example

Also found in the example directory

```go
package main

import (
	"errors"
	"log"
	"node2go"
)

// Data will be used when parsing json and will be available in receive function
type Data struct {
	N int `json:"n"`
}

// same old, same old
func fibonacci(n int) int {
	if n == 0 {
		return n
	}
	if n < 2 {
		return 1
	}
	return fibonacci(n-2) + fibonacci(n-1)
}

// receive is a func of type node2go.HandlerFunc
// type assertion converts data to the actual *Data type
func receive(data interface{}) (res interface{}, err error) {
	d, ok := data.(*Data)

	if !ok {
		err = errors.New("failed to parse data")
		return
	}

	return fibonacci(d.N), nil

}

func main() {

	// m holds the functions.
	// Eg. the Handler.Func will be called for (NodeJS):
	// await go.call('fibonacci', {n: 30})
	// and {n: 30} will be loaded into *Data copy
	m := make(map[string]*node2go.Handler)
	m["fibonacci"] = &node2go.Handler{
		Func:         receive,
		DataTemplate: &Data{},
	}

	// create a new Runner
	r, err := node2go.NewRunner(m)
	if err != nil {
		log.Fatal(err)
	}

	// run in goroutine to be able to properly catch signals (Close())
	go r.Run()

	// graceful shutdown (remove the unix domain socket file and exit)
	r.Close()

}

```