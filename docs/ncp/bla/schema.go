package bla

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"ncp/templates"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/ko/pkg/build"
	"github.com/google/ko/pkg/publish"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type SchemaActor struct {
	// Conn is a NATS connection to the actor account.
	Conn *nats.Conn

	js     jetstream.JetStream
	bucket nats.KeyValue
	subs   []*nats.Subscription
}

func (sc *SchemaActor) Close() {
	for _, sub := range sc.subs {
		sub.Unsubscribe()
	}
	sc.Conn.Close()
}

func RegisterSchemaActor(ctx context.Context, actor *SchemaActor) error {
	js, err := jetstream.New(actor.Conn)
	if err != nil {
		return fmt.Errorf("jetstream instance: %w", err)
	}
	actor.js = js

	// Use legacy (but not that legacy) Jetstream API to create
	// a KeyValue bucket, because the new JetStream API doesn't
	// support it...
	jsc, err := actor.Conn.JetStream()
	if err != nil {
		return fmt.Errorf("jetstream context: %w", err)
	}
	bucket, err := jsc.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "schema",
	})
	if err != nil {
		return fmt.Errorf("create key value: %w", err)
	}
	actor.bucket = bucket

	name := "schema"
	{
		action := "publish"
		sub, err := actor.Conn.QueueSubscribe(
			fmt.Sprintf(ActionSubjectSubscribeForAccount, name, action),
			name,
			func(msg *nats.Msg) {
				reply, err := actor.publishSchema(ctx, msg)
				if err != nil {
					slog.Error("action", "actor", name, "action", action, "error", err)
					reply = &SchemaPublishReply{
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
		slog.Info("subscribed action", "actor", name, "action", action, "subject", sub.Subject)
		actor.subs = append(actor.subs, sub)
	}
	{
		action := "list"
		sub, err := actor.Conn.QueueSubscribe(
			fmt.Sprintf(ActionSubjectSubscribeForAccount, name, action),
			name,
			func(msg *nats.Msg) {
				reply, err := actor.schemaList(ctx, msg)
				if err != nil {
					slog.Error("action", "actor", name, "action", action, "error", err)
					reply = &SchemaListReply{
						Schemas: []Schema{},
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
		slog.Info("subscribed action", "actor", name, "action", action, "subject", sub.Subject)
		actor.subs = append(actor.subs, sub)
	}
	{
		action := "get"
		sub, err := actor.Conn.QueueSubscribe(
			fmt.Sprintf(ActionSubjectSubscribeForAccount, name, action),
			name,
			func(msg *nats.Msg) {
				reply, err := actor.schemaGet(ctx, msg)
				if err != nil {
					slog.Error("action", "actor", name, "action", action, "error", err)
					reply = &SchemaGetReply{
						Schema: Schema{},
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
		slog.Info("subscribed action", "actor", name, "action", action, "subject", sub.Subject)
		actor.subs = append(actor.subs, sub)
	}

	return nil
}

func (sa *SchemaActor) schemaList(
	ctx context.Context,
	msg *nats.Msg,
) (*SchemaListReply, error) {
	var data SchemaListMsg
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	keys, err := sa.bucket.Keys()
	if err != nil {
		if errors.Is(err, nats.ErrNoKeysFound) {
			return &SchemaListReply{
				Schemas: []Schema{},
			}, nil
		}
		return nil, fmt.Errorf("keys: %w", err)
	}
	schemas := make([]Schema, len(keys))
	for i, key := range keys {
		kve, err := sa.bucket.Get(key)
		if err != nil {
			return nil, fmt.Errorf("get: %w", err)
		}
		var schema Schema
		if err := json.Unmarshal(kve.Value(), &schema); err != nil {
			return nil, fmt.Errorf("unmarshal: %w", err)
		}
		schemas[i] = schema
	}
	return &SchemaListReply{
		Schemas: schemas,
	}, nil
}

func (sa *SchemaActor) schemaGet(
	ctx context.Context,
	msg *nats.Msg,
) (*SchemaGetReply, error) {
	var data SchemaGetMsg
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	kve, err := sa.bucket.Get(data.Key)
	if err != nil {
		if !errors.Is(err, nats.ErrKeyNotFound) {
			return nil, fmt.Errorf("get: %w", err)
		}
		return nil, fmt.Errorf("schema not found")
	}
	var schema Schema
	if err := json.Unmarshal(kve.Value(), &schema); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &SchemaGetReply{
		Schema: schema,
	}, nil
}

func (sa *SchemaActor) publishSchema(
	ctx context.Context,
	msg *nats.Msg,
) (*SchemaPublishReply, error) {
	// Get account public key from subject
	accountPubKey, err := AccountFromSubject(msg.Subject)
	if err != nil {
		return nil, fmt.Errorf("account from subject: %w", err)
	}
	var data SchemaPublishMsg
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

	schema := Schema{
		Spec:   data,
		Events: []Event{},
	}
	schema.Events = append(schema.Events, Event{
		Status:  EventStatusInfo,
		Message: "schema published",
	})

	// Stream is named by the event and contains all the events for all versions of that schema.
	streamName := data.Name
	// A durable consumer is created for each version of the schema, and that consumer can be used
	// for monitoring and metrics.
	consumerName := strings.ReplaceAll(data.Version, ".", "_")
	schemaSubject := fmt.Sprintf("ingest.%s.%s", streamName, consumerName)

	// Create stream for the ingestion service in the target account.
	// First we need to create some credentials for the target account, and use those
	// to establish a nats connection.
	userReply, err := SendUserCreateForAccountMsg(sa.Conn, accountPubKey, UserCreateMsg{
		Name: "ingest",
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	nc, err := nats.Connect(
		nats.DefaultURL,
		nats.UserJWTAndSeed(userReply.JWT, string(userReply.NKey)),
	)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	defer nc.Close()
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("jetstream instance: %w", err)
	}

	// Ensure the stream exists for the schema name.
	// We use the same stream for all versions of the schema.
	// Each version has its own consumer.
	streamCfg := jetstream.StreamConfig{
		Name: streamName,
		// Subject is ingest.<schema-name>.<schema-version>
		Subjects: []string{
			fmt.Sprintf("ingest.%s.*", streamName),
		},
	}
	if _, err := js.CreateStream(ctx, streamCfg); err != nil {
		if !errors.Is(err, jetstream.ErrStreamNameAlreadyInUse) {
			return nil, fmt.Errorf("create stream: %w", err)
		}
		// Update stream if it already exists
		if _, err = js.UpdateStream(ctx, streamCfg); err != nil {
			return nil, fmt.Errorf("update stream: %w", err)
		}
		slog.Info("updated stream", "name", streamName)
	} else {
		slog.Info("created stream", "name", streamName)
	}
	// Create the consumer for the schema version.
	if _, err := js.CreateOrUpdateConsumer(ctx, streamName, jetstream.ConsumerConfig{
		Name:    consumerName,
		Durable: consumerName,
		Description: fmt.Sprintf(
			"Consumer for schema %s:%s",
			schema.Spec.Name,
			schema.Spec.Version,
		),
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: schemaSubject,
	}); err != nil {
		return nil, fmt.Errorf("create consumer: %w", err)
	}
	slog.Info("created or updated consumer", "consumer", consumerName, "stream", streamName)

	// Build the schema into a container
	ref, err := buildSchema(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("build schema: %w", err)
	}
	schema.Events = append(schema.Events, Event{
		Status:  EventStatusInfo,
		Message: "schema built into container image " + ref,
	})

	// Create a user for the container by sending a message to the account actor
	// reply, err := SendUserCreateForAccountMsg(sa.Conn, accountPubKey, UserCreateMsg{
	// 	Name: "bob-for-ingest",
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("create user: %w", err)
	// }

	// Run the container
	containerID, err := runSchema(
		ctx,
		ref,
		streamName,
		consumerName,
		userReply.JWT,
		string(userReply.NKey),
	)
	if err != nil {
		return nil, fmt.Errorf("run schema: %w", err)
	}
	schema.Events = append(schema.Events, Event{
		Status:  EventStatusInfo,
		Message: "schema running in container " + containerID,
	})

	bSchema, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("marshal schema: %w", err)
	}
	if _, err := sa.bucket.Put(schema.Key(), bSchema); err != nil {
		return nil, fmt.Errorf("put schema: key %s: %w", schema.Key(), err)
	}

	return &SchemaPublishReply{
		Name: data.Name,
	}, nil
}

func buildSchema(ctx context.Context, schema SchemaPublishMsg) (string, error) {
	dir, err := os.MkdirTemp("", schema.Name)
	if err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}
	// defer os.RemoveAll(dir)
	slog.Info("created temp dir", "dir", dir)
	// Unpack ingester service template
	if err := UnpackTxtar(dir, templates.IngesterTxtar); err != nil {
		return "", fmt.Errorf("unpack ingester: %w", err)
	}
	// Unpack protobuf schema
	schemaDir := filepath.Join(dir, "schema")
	if err := UnpackTxtar(schemaDir, schema.Schema); err != nil {
		return "", fmt.Errorf("unpack: %w", err)
	}
	if err := protoBuild(dir); err != nil {
		return "", fmt.Errorf("proto build: %w", err)
	}

	ref, err := koBuild(ctx, dir)
	if err != nil {
		return "", fmt.Errorf("ko build: %w", err)
	}
	return ref, nil
}

func koBuild(ctx context.Context, dir string) (string, error) {
	plat, err := v1.ParsePlatform("linux/arm64")
	if err != nil {
		return "", fmt.Errorf("parse platform: %w", err)
	}
	ref, err := name.ParseReference("cgr.dev/chainguard/static:latest")
	if err != nil {
		return "", fmt.Errorf("parse reference: %w", err)
	}
	desc, err := remote.Get(ref, remote.WithContext(ctx), remote.WithPlatform(*plat))
	if err != nil {
		return "", fmt.Errorf("get remote container image descriptor: %w", err)
	}
	base, err := desc.Image()
	if err != nil {
		return "", fmt.Errorf("get image: %w", err)
	}
	bi, err := build.NewGo(
		ctx,
		dir,
		build.WithDisabledSBOM(),
		build.WithPlatforms(plat.String()),
		build.WithBaseImages(
			func(ctx context.Context, s string) (name.Reference, build.Result, error) {
				return ref, base, nil
			},
		),
	)
	if err != nil {
		return "", fmt.Errorf("new go builder: %w", err)
	}
	importpath, err := bi.QualifyImport(".")
	if err != nil {
		return "", fmt.Errorf("qualify import: %w", err)
	}
	if err := bi.IsSupportedReference(importpath); err != nil {
		return "", fmt.Errorf("is supported reference: %w", err)
	}
	result, err := bi.Build(ctx, importpath)
	if err != nil {
		return "", fmt.Errorf("build: %w", err)
	}

	pi, err := publish.NewDaemon(
		packageWithMD5,
		// func(s1, s2 string) string {
		// 	return s1 + "-whatever-" + s2
		// },
		[]string{},
	)
	if err != nil {
		return "", fmt.Errorf("new publish daemon: %w", err)
	}
	pr, err := pi.Publish(ctx, result, importpath)
	if err != nil {
		return "", fmt.Errorf("publish: %w", err)
	}
	return pr.Name(), nil
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

func runSchema(
	ctx context.Context,
	ref string,
	stream string,
	consumer string,
	jwt string,
	nkey string,
) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("new docker client: %w", err)
	}
	defer cli.Close()

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: ref,
		Env: []string{
			"NATS_URL=" + "nats://host.docker.internal:4222",
			"NATS_JWT=" + jwt,
			"NATS_NKEY=" + nkey,
		},
		Cmd: []string{
			"-stream=" + stream,
			"-consumer=" + consumer,
		},
		Tty: false,
	}, nil, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("create container: %w", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("start container: %w", err)
	}
	return resp.ID, nil
}

func packageWithMD5(base, importpath string) string {
	hasher := md5.New() // nolint: gosec // No strong cryptography needed.
	hasher.Write([]byte(importpath))
	return path.Join(base, path.Base(importpath)+"-"+hex.EncodeToString(hasher.Sum(nil)))
}

// func goRun(dir string, stream string, consumer string) error {
// 	cmd := exec.Command("go", "run", ".", "-stream="+stream, "-consumer="+consumer)
// 	cmd.Dir = dir
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	slog.Info("running go", "exec", cmd.String())

// 	if err := cmd.Start(); err != nil {
// 		return fmt.Errorf("starting go run: %w", err)
// 	}

// 	return nil
// }

type SchemaListMsg struct {
}

type SchemaListReply struct {
	Schemas []Schema `json:"schemas"`
}

type SchemaGetMsg struct {
	Key string `json:"key"`
}

type SchemaGetReply struct {
	Schema Schema `json:"schema"`
}

type SchemaPublishMsg struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Schema  []byte `json:"schema"`
}

type SchemaPublishReply struct {
	Name  string `json:"name"`
	Error string `json:"error,omitempty"`
}

func SendSchemaListForAccountMsg(
	nc *nats.Conn,
	accountID string,
	msg SchemaListMsg,
) (*SchemaListReply, error) {
	bData, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	replyMsg, err := nc.Request(
		fmt.Sprintf(ActionSubjectSendForAccount, "schema", "list", accountID),
		bData,
		time.Second*10,
	)
	if err != nil {
		return nil, fmt.Errorf("request schema list: %w", err)
	}
	var reply SchemaListReply
	if err := json.Unmarshal(replyMsg.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}

func SendSchemaListMsg(nc *nats.Conn, msg SchemaListMsg) (*SchemaListReply, error) {
	bData, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	replyMsg, err := nc.Request(
		fmt.Sprintf(ActionSubjectSend, "schema", "list"),
		bData,
		time.Second*10,
	)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	var reply SchemaListReply
	if err := json.Unmarshal(replyMsg.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}

func SendSchemaGetForAccountMsg(
	nc *nats.Conn,
	accountID string,
	msg SchemaGetMsg,
) (*SchemaGetReply, error) {
	bData, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	replyMsg, err := nc.Request(
		fmt.Sprintf(ActionSubjectSendForAccount, "schema", "get", accountID),
		bData,
		time.Second*10,
	)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	var reply SchemaGetReply
	if err := json.Unmarshal(replyMsg.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}
func SendSchemaGetMsg(nc *nats.Conn, msg SchemaGetMsg) (*SchemaGetReply, error) {
	bData, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	replyMsg, err := nc.Request(
		fmt.Sprintf(ActionSubjectSend, "schema", "get"),
		bData,
		time.Second*10,
	)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	var reply SchemaGetReply
	if err := json.Unmarshal(replyMsg.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}

func SendSchemaPublishMsg(nc *nats.Conn, msg SchemaPublishMsg) (*SchemaPublishReply, error) {
	bData, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	replyMsg, err := nc.Request(
		fmt.Sprintf(ActionSubjectSend, "schema", "publish"),
		bData,
		time.Second*10,
	)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	var reply SchemaPublishReply
	if err := json.Unmarshal(replyMsg.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}

type Schema struct {
	Spec   SchemaPublishMsg `json:"spec"`
	Events []Event          `json:"events"`
}

func (s Schema) Key() string {
	// Cannot use ":" as it's not allowed, but "=" is apparently...
	return fmt.Sprintf("%s=%s", s.Spec.Name, s.Spec.Version)
}

type EventStatus string

const (
	EventStatusInfo  EventStatus = "info"
	EventStatusError EventStatus = "error"
)

type Event struct {
	Status  EventStatus `json:"status"`
	Message string      `json:"message"`
}
