package bla

import (
	"context"
	"crypto/md5"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"ncp/templates"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/ko/pkg/build"
	"github.com/google/ko/pkg/publish"
	"github.com/nats-io/jsm.go/api"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

//go:embed schema.tmpl.html
var schemaHTML string

//go:embed schemas.tmpl.html
var schemasHTML string

type SchemaActor struct {
	// Conn is a NATS connection to the actor account.
	Conn *nats.Conn

	ctx context.Context

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
	actor.ctx = context.Background()
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

	{
		r := NewRouter()
		r.Get("/", actor.renderSchemas)
		r.Get("/{schemaID}", actor.renderSchema)
		sub, err := SubscribeToSubject[Request, Reply](
			actor.Conn,
			fmt.Sprintf(ActorSubjectHTTPRender, "schema"),
			r.Serve,
		)
		if err != nil {
			return fmt.Errorf("subscribing: %w", err)
		}
		actor.subs = append(actor.subs, sub)
	}
	{
		sub, err := SubscribeToSubjectWithAccount[SchemaPublishMsg, SchemaPublishReply](
			actor.Conn,
			fmt.Sprintf(ActionSubjectSubscribeForAccount, "schema", "publish"),
			actor.publishSchema,
		)
		if err != nil {
			return fmt.Errorf("subscribing: %w", err)
		}
		actor.subs = append(actor.subs, sub)
	}
	{
		sub, err := SubscribeToSubjectWithAccount[SchemaListMsg, SchemaListReply](
			actor.Conn,
			fmt.Sprintf(ActionSubjectSubscribeForAccount, "schema", "list"),
			actor.schemaList,
		)
		if err != nil {
			return fmt.Errorf("subscribing: %w", err)
		}
		actor.subs = append(actor.subs, sub)
	}
	{
		sub, err := SubscribeToSubjectWithAccount[SchemaGetMsg, SchemaGetReply](
			actor.Conn,
			fmt.Sprintf(ActionSubjectSubscribeForAccount, "schema", "get"),
			actor.schemaGet,
		)
		if err != nil {
			return fmt.Errorf("subscribing: %w", err)
		}
		actor.subs = append(actor.subs, sub)
	}

	return nil
}

func (sa *SchemaActor) renderSchemas(req *Request) (*Reply, error) {
	schemaListReply, err := SendSchemaListForAccountMsg(
		sa.Conn,
		req.Account.ID,
		SchemaListMsg{},
	)
	if err != nil {
		return nil, fmt.Errorf("getting schema list: %w", err)
	}
	type Data struct {
		UserInfo UserInfo
		Account  Account
		Schemas  []Schema
	}
	data := Data{
		UserInfo: req.UserInfo,
		Account:  req.Account,
		Schemas:  schemaListReply.Schemas,
	}

	body, err := renderTemplate(schemasHTML, data)
	if err != nil {
		return nil, fmt.Errorf("rendering template: %w", err)
	}
	return &Reply{
		Body: body,
	}, nil
}

func (sa *SchemaActor) renderSchema(req *Request) (*Reply, error) {
	schemaID, ok := req.Params["schemaID"]
	if !ok {
		return nil, fmt.Errorf("no schemaID")
	}
	schemaReply, err := SendSchemaGetForAccountMsg(
		sa.Conn,
		req.Account.ID,
		SchemaGetMsg{
			Key: schemaID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("getting schema: %w", err)
	}
	schema := schemaReply.Schema
	consumerInfo, err := ConsumerState(
		sa.Conn,
		req.Account.ID,
		schema.Spec.Name,
		strings.ReplaceAll(schema.Spec.Version, ".", "_"),
	)
	if err != nil {
		return nil, fmt.Errorf("getting consumer state: %w", err)
	}
	type Data struct {
		UserInfo     UserInfo
		Account      Account
		Schema       Schema
		ConsumerInfo *api.ConsumerInfo
	}
	d := Data{
		UserInfo:     req.UserInfo,
		Account:      req.Account,
		Schema:       schemaReply.Schema,
		ConsumerInfo: consumerInfo,
	}
	body, err := renderTemplate(schemaHTML, d)
	if err != nil {
		return nil, fmt.Errorf("rendering template: %w", err)
	}
	return &Reply{
		Body: body,
	}, nil
}

func (sa *SchemaActor) schemaList(
	accountID string,
	msg *SchemaListMsg,
) (*SchemaListReply, error) {
	// TODO: use accountID to separate schemas out...
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
	accountID string,
	msg *SchemaGetMsg,
) (*SchemaGetReply, error) {
	kve, err := sa.bucket.Get(msg.Key)
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
	accountPubKey string,
	msg *SchemaPublishMsg,
) (*SchemaPublishReply, error) {
	if msg.Name == "" {
		return nil, FromError(&Error{
			Status:  http.StatusBadRequest,
			Message: "Schema name is empty",
		})
	}
	if msg.Version == "" {
		return nil, FromError(&Error{
			Status:  http.StatusBadRequest,
			Message: "Schema version is empty",
		})
	}
	if len(msg.Schema) == 0 {
		return nil, FromError(&Error{
			Status:  http.StatusBadRequest,
			Message: "Schema body is empty",
		})
	}

	schema := Schema{
		Spec:   *msg,
		Events: []Event{},
	}
	schema.Events = append(schema.Events, Event{
		Status:  EventStatusInfo,
		Message: "schema published",
	})

	// Stream is named by the event and contains all the events for all versions of that schema.
	streamName := msg.Name
	// A durable consumer is created for each version of the schema, and that consumer can be used
	// for monitoring and metrics.
	consumerName := strings.ReplaceAll(msg.Version, ".", "_")
	schemaSubject := fmt.Sprintf("ingest.%s.%s", streamName, consumerName)

	// Create stream for the ingestion service in the target account.
	// First we need to create some credentials for the target account, and use those
	// to establish a nats connection.
	userReply, err := SendUserCreateForAccountMsg(
		sa.Conn,
		accountPubKey,
		UserCreateMsg{
			Name: "ingest",
		},
	)
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
	if _, err := js.CreateStream(sa.ctx, streamCfg); err != nil {
		if !errors.Is(err, jetstream.ErrStreamNameAlreadyInUse) {
			return nil, fmt.Errorf("create stream: %w", err)
		}
		// Update stream if it already exists
		if _, err = js.UpdateStream(sa.ctx, streamCfg); err != nil {
			return nil, fmt.Errorf("update stream: %w", err)
		}
		slog.Info("updated stream", "name", streamName)
	} else {
		slog.Info("created stream", "name", streamName)
	}
	// Create the consumer for the schema version.
	if _, err := js.CreateOrUpdateConsumer(sa.ctx, streamName, jetstream.ConsumerConfig{
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
	slog.Info(
		"created or updated consumer",
		"consumer",
		consumerName,
		"stream",
		streamName,
	)

	// Build the schema into a container
	ref, err := buildSchema(sa.ctx, msg)
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
		sa.ctx,
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
		Name: msg.Name,
	}, nil
}

func buildSchema(
	ctx context.Context,
	schema *SchemaPublishMsg,
) (string, error) {
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
	desc, err := remote.Get(
		ref,
		remote.WithContext(ctx),
		remote.WithPlatform(*plat),
	)
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
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
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
	return path.Join(
		base,
		path.Base(importpath)+"-"+hex.EncodeToString(hasher.Sum(nil)),
	)
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
	return RequestSubject[SchemaListMsg, SchemaListReply](
		nc,
		msg,
		fmt.Sprintf(ActionSubjectSendForAccount, "schema", "list", accountID),
	)
}

func SendSchemaListMsg(
	nc *nats.Conn,
	msg SchemaListMsg,
) (*SchemaListReply, error) {
	return RequestSubject[SchemaListMsg, SchemaListReply](
		nc,
		msg,
		fmt.Sprintf(ActionSubjectSend, "schema", "list"),
	)
}

func SendSchemaGetForAccountMsg(
	nc *nats.Conn,
	accountID string,
	msg SchemaGetMsg,
) (*SchemaGetReply, error) {
	return RequestSubject[SchemaGetMsg, SchemaGetReply](
		nc,
		msg,
		fmt.Sprintf(ActionSubjectSendForAccount, "schema", "get", accountID),
	)
}

func SendSchemaGetMsg(
	nc *nats.Conn,
	msg SchemaGetMsg,
) (*SchemaGetReply, error) {
	return RequestSubject[SchemaGetMsg, SchemaGetReply](
		nc,
		msg,
		fmt.Sprintf(ActionSubjectSend, "schema", "get"),
	)
}

func SendSchemaPublishMsg(
	nc *nats.Conn,
	msg SchemaPublishMsg,
) (*SchemaPublishReply, error) {
	return RequestSubject[SchemaPublishMsg, SchemaPublishReply](
		nc,
		msg,
		fmt.Sprintf(ActionSubjectSend, "schema", "publish"),
	)
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
