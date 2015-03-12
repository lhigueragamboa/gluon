package gluon

import (
	"log"
	"os"
)

var (
	Logger           = log.New(os.Stderr, "gluon: ", log.LstdFlags|log.Lmicroseconds)
	LogHandlerErrors bool
)
