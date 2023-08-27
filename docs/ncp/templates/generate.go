package templates

//go:generate echo ">>>> Txtaring ingester"
//go:generate go run ../cmd/txtarfiles/main.go -dir=ingester -in=main.go -in=go.mod -in=go.sum -out=ingester.txtar
