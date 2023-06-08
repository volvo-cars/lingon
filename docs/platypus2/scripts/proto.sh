#!/usr/bin/env bash
# Copyright (c) 2023 Volvo Car Corporation
# SPDX-License-Identifier: Apache-2.0

set -e
set -x
set -u
set -o pipefail

command -v protoc > /dev/null
command -v protoc-gen-go > /dev/null
command -v protoc-gen-go-grpc > /dev/null

ROOT_DIR=$(git rev-parse --show-toplevel)

PROTO="$ROOT_DIR/docs/platypus2/proto"
OUT="$ROOT_DIR/docs/platypus2/cmd/tools/nats"
PACKAGE="github.com/volvo-cars/lingoneks"

protoc -I="$PROTO" \
  --go_opt=module="$PACKAGE" \
  --go_out="$OUT" \
  --go-grpc_opt=module="$PACKAGE" \
  --go-grpc_out="$OUT" \
  "$PROTO"/*.proto

set +x
echo
echo "    protobuf files compiled"
echo
echo "    in: $PROTO"
echo "    out: $OUT"
echo


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