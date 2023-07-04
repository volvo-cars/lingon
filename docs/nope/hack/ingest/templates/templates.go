package templates

import (
	_ "embed"
)

//go:embed ingestion.txtar
var IngestionTxtar []byte

//go:embed streamer.txtar
var StreamerTxtar []byte
