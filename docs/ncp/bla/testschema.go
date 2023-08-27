package bla

import (
	_ "embed"
)

//go:generate echo ">>>> Txtaring event"
//go:generate go run ../cmd/txtarfiles/main.go -dir=testdata/event -in=event.proto -out=event.pb.txtar

//go:embed event.pb.txtar
var EventTxtar []byte
