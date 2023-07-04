package templates

//go:generate echo ">>>> Txtaring ingestion"
//go:generate go run ../../txtarfiles/main.go -dir=ingestion -in=main.go -in=go.mod -in=go.sum -out=ingestion.txtar

//go:generate echo ">>>> Txtaring streamer"
//go:generate go run ../../txtarfiles/main.go -dir=streamer -in=main.go -in=go.mod -in=go.sum -out=streamer.txtar
