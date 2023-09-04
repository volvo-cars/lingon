#!/usr/bin/env bash

set -xeuo pipefail

bench='bench foo --kv --no-progress --multisubject --bucket bar --storage file --msgs 65536 --maxbytes 4500MB'

for s in 64 256 1k 4k 16k 64k; do
#  nats kv add bar --max-bucket-size 4500MB --storage file >/dev/null 2>&1

  nats ${bench} --size ${s} --purge --pub 4
  nats ${bench} --size ${s} --sub 1
done
