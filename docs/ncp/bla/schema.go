package bla

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"ncp/templates"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type SchemaActor struct {
	// Conn is a NATS connection to the actor account.
	Conn *nats.Conn

	subs []*nats.Subscription
}

func (sc *SchemaActor) Close() {
	for _, sub := range sc.subs {
		sub.Unsubscribe()
	}
	sc.Conn.Close()
}

func RegisterSchemaActor(ctx context.Context, actor *SchemaActor) error {
	name := "schema"
	{
		action := "publish"
		sub, err := actor.Conn.QueueSubscribe(
			fmt.Sprintf("actor.%s.%s.*", name, action),
			name,
			func(msg *nats.Msg) {
				reply, err := publishSchemaAction(ctx, actor.Conn, msg)
				if err != nil {
					slog.Error("action", "error", err)
					reply = &PublishSchemaReply{
						Error: err.Error(),
					}
				}
				bReply, err := json.Marshal(reply)
				if err != nil {
					slog.Error("marshal", "error", err)
					return
				}
				if err := msg.Respond(bReply); err != nil {
					slog.Error("respond", "error", err)
				}
			},
		)
		if err != nil {
			return fmt.Errorf("subscribe: %w", err)
		}
		actor.subs = append(actor.subs, sub)
	}

	return nil
}

type PublishSchemaMsg struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Schema  []byte `json:"schema"`
}

type PublishSchemaReply struct {
	Name  string `json:"name"`
	Error string `json:"error,omitempty"`
}

func PublishSchema(nc *nats.Conn, msg PublishSchemaMsg) (*PublishSchemaReply, error) {
	bData, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	replyMsg, err := nc.Request("actor.schema.publish", bData, time.Second*10)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	var reply PublishSchemaReply
	if err := json.Unmarshal(replyMsg.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}

func publishSchemaAction(
	ctx context.Context,
	nc *nats.Conn,
	msg *nats.Msg,
) (*PublishSchemaReply, error) {
	var data PublishSchemaMsg
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	if data.Name == "" {
		return nil, fmt.Errorf("name is empty")
	}
	if data.Version == "" {
		return nil, fmt.Errorf("version is empty")
	}
	if len(data.Schema) == 0 {
		return nil, fmt.Errorf("schema is empty")
	}

	// TODO: store schema in KV store

	// Create stream for the ingestion service
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("jetstream: %w", err)
	}
	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     data.Name,
		Subjects: []string{fmt.Sprintf("ingest.%s", data.Name)},
	})
	if err != nil {
		if !errors.Is(err, jetstream.ErrStreamNameAlreadyInUse) {
			return nil, fmt.Errorf("create stream: %w", err)
		}
		// TODO: handle update of stream
	}
	slog.Info("created stream", "config", stream.CachedInfo().Config)

	if err := buildSchema(data); err != nil {
		return nil, fmt.Errorf("build schema: %w", err)
	}

	return &PublishSchemaReply{
		Name: data.Name,
	}, nil
}

func buildSchema(schema PublishSchemaMsg) error {
	dir, err := os.MkdirTemp("", schema.Name)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	// defer os.RemoveAll(dir)
	slog.Info("created temp dir", "dir", dir)
	// Unpack ingester service template
	if err := UnpackTxtar(dir, templates.IngesterTxtar); err != nil {
		return fmt.Errorf("unpack ingester: %w", err)
	}
	// Unpack protobuf schema
	schemaDir := filepath.Join(dir, "schema")
	if err := UnpackTxtar(schemaDir, schema.Schema); err != nil {
		return fmt.Errorf("unpack: %w", err)
	}
	if err := protoBuild(dir); err != nil {
		return fmt.Errorf("proto build: %w", err)
	}
	if err := goRun(dir, schema.Name, schema.Name); err != nil {
		return fmt.Errorf("go run: %w", err)
	}

	return nil
}

func protoBuild(dir string) error {
	bufTmpl := `{"version":"v1","managed":{"enabled":true,"go_package_prefix":{"default":"ingester"}},"plugins":[{"plugin":"go","out":".","opt":"paths=source_relative"}]}`
	cmd := exec.Command("buf", "generate", "--template", bufTmpl)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	slog.Info("generating protobuf files", "exec", cmd.String())

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running buf: %w", err)
	}

	return nil
}

func goRun(dir string, stream string, consumer string) error {
	cmd := exec.Command("go", "run", ".", "-stream="+stream, "-consumer="+consumer)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	slog.Info("running go", "exec", cmd.String())

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting go run: %w", err)
	}

	return nil
}
