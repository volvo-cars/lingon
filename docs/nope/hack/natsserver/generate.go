package main

//go:generate echo ">>>> Bootstrapping NATS operator"
//go:generate go run ../tools/bootstrapnatsoperator/main.go -out=out/nats -force
