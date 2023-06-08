#!/usr/bin/env bash

echo
echo "Not meant to be run."
echo "See the script instead."
echo
exit 1

# Start nats
nats-server -V

# run server
go run .

# try the service

# list all remote endpoints
grpcurl -plaintext 0.0.0.0:7015 list

# send data
grpcurl -plaintext -d @ 0.0.0.0:7015 natspb.EnvelopeService/Ingest <<EOM
{
  "topic": "demo.top",
  "msg": {
    "id": "id1",
    "author_id": "author1",
    "title": "nice demo",
    "content": "nice content"
  }
}
EOM
