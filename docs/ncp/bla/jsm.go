package bla

import (
	"fmt"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/jsm.go/api"
	"github.com/nats-io/nats.go"
)

func ConsumerState(
	nc *nats.Conn,
	account string,
	stream string,
	consumer string,
) (*api.ConsumerInfo, error) {
	mgr, err := jsm.New(
		nc,
		// Prefix the API with the account public key to get the state from
		// that specific account.
		// This is mapped via an import to the target accounts $JS.API subject.
		jsm.WithAPIPrefix("$JS.API."+account),
	)
	if err != nil {
		return nil, fmt.Errorf("creating jsm manager: %w", err)
	}
	jstream, err := mgr.LoadStream(stream)
	if err != nil {
		return nil, fmt.Errorf("loading stream: %w", err)
	}
	jconsumer, err := jstream.LoadConsumer(consumer)
	if err != nil {
		return nil, fmt.Errorf("loading consumer: %w", err)
	}
	ci, err := jconsumer.LatestState()
	if err != nil {
		return nil, fmt.Errorf("getting consumer latest state: %w", err)
	}
	return &ci, nil
	// si, err := jstream.LatestInformation()
	// if err != nil {
	// 	return nil, fmt.Errorf("getting stream latest information: %w", err)
	// }
	// return si, nil
}
