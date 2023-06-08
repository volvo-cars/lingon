// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package benthos

var DemoConfig = `http:
  address: 0.0.0.0:4195
  enabled: true
  root_path: /benthos
  debug_endpoints: false
  cors:
    enabled: true
    allowed_origins: ["*"]
input:
  generate:
    interval: 1s
    mapping: |
      #!blobl
      # https://www.benthos.dev/docs/guides/bloblang/about
      root.id = uuid_v4()
      root.doc.id = uuid_v4()
      root.doc.received_at = now()
      root.doc.host = hostname()
      root.topic = "${NATS_TOPIC:demo.output}"
      meta hellew = {
        "interesting": "true",
        "timestamp": now(),
        "nats_kv": 1234
      }
pipeline:
  processors:
  - bloblang: |
      #!blobl
      root = this
      root.id = count("id_for_counter") # https://www.benthos.dev/docs/guides/bloblang/functions/
      root.meta_obj = @  # @ represent all metadata k/v
      root.interesting = @.hellew.interesting
      meta = deleted()
      meta ts = now()
      meta topic = root.topic
      root.topic = deleted()

output:
  nats:
    urls: [ "${NATS_URLs:nats://nats.nats.svc.cluster.local:4222}" ]
    subject: ${!meta("topic")}
    headers:
      Content-Type: application/json
      Timestamp: ${!meta("ts")}

metrics:
  prometheus:
    add_process_metrics: true
    add_go_metrics: true
shutdown_delay: ""
shutdown_timeout: 20s

`

// DemoConfigSimple is a simple Benthos configuration for demo purposes
// Benthos configuration can be generated with `benthos create "in/pipe/out"`.
const DemoConfigSimple = `http:
  address: 0.0.0.0:4195
  enabled: true
  root_path: /benthos
  debug_endpoints: false
input:
  generate:
    interval: 1s
    mapping: |
    root.id = uuid_v4()
output:
  nats:
    urls: [ ${NATS_URL:-"nats://nats.nats.svc.cluster.local:4222" } ]
    subject: demo.output
metrics:
  prometheus:
	add_process_metrics: true
	add_go_metrics: true
shutdown_delay: ""
shutdown_timeout: 20s
`

// DemoOutputStdOut is a Benthos output to simply print the messages to stdout.
// It is useful for testing/validating.
const DemoOutputStdOut = `
  stdout:
    codec: lines
`
