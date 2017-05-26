package main

import (
	"github.com/bmorri12/SmartAqua/pkg/protocol"
)

var StatusChan map[uint64]chan *protocol.Data

func init() {
	StatusChan = make(map[uint64]chan *protocol.Data)
}
