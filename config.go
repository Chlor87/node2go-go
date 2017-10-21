package node2go

import (
	"flag"
	"log"
)

var (
	socketAddr string
)

func init() {
	flag.StringVar(&socketAddr, "addr", "", "socket address")
	flag.Parse()
	if socketAddr == "" {
		log.Fatal("Please provide socket address.")
	}
}
