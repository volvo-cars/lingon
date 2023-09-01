package bla

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

func StartTestServer(t *testing.T) TestServer {
	ts, err := NewTestServer(
		WithDir(t.TempDir()),
		WithFindAvailablePort(true),
		// WithConfigureLogger(true),
	)
	if err != nil {
		t.Fatal("creating test server: ", err)
	}
	if err := ts.StartUntilReady(); err != nil {
		t.Fatal("starting test server: ", err)
	}

	if err := ts.PublishActorAccount(); err != nil {
		t.Fatal("publishing actor account: ", err)
	}

	t.Cleanup(func() {
		ts.NS.Shutdown()
	})
	return ts
}

type TestServer struct {
	NS   *server.Server
	Auth ServerJWTAuth
}

func (t TestServer) StartUntilReady() error {
	timeout := 4 * time.Second
	t.NS.Start()
	if !t.NS.ReadyForConnections(timeout) {
		return fmt.Errorf("server not ready after %s", timeout.String())
	}
	return nil
}

func (t TestServer) SysUserConn() (*nats.Conn, error) {
	nc, err := nats.Connect(
		t.NS.ClientURL(),
		UserJWTOption(t.Auth.SysUserJWT, t.Auth.SysUserKeyPair),
	)
	if err != nil {
		return nil, fmt.Errorf("connecting to nats: %w", err)
	}
	return nc, nil
}

func (t TestServer) ActorUserConn() (*nats.Conn, error) {
	nc, err := nats.Connect(
		t.NS.ClientURL(),
		UserJWTOption(t.Auth.ActorUserJWT, t.Auth.ActorUserKeyPair),
	)
	if err != nil {
		return nil, fmt.Errorf("connecting to nats: %w", err)
	}
	return nc, nil
}

// PublishActorAccount publishes the actor account JWT to the server.
// This is required to be able to connect to the server with the actor account
// and therefore needs to happen once the server has been started.
func (t TestServer) PublishActorAccount() error {
	nc, err := t.SysUserConn()
	if err != nil {
		return fmt.Errorf("connecting to system_account nats: %w", err)
	}
	defer nc.Close()
	// Check if the account already exists
	lookupMessage, err := nc.Request(
		fmt.Sprintf("$SYS.REQ.ACCOUNT.%s.CLAIMS.LOOKUP", t.Auth.ActorAccountPublicKey),
		nil,
		time.Second,
	)
	if err != nil {
		return fmt.Errorf(
			"looking up account with ID %s: %w",
			t.Auth.ActorAccountPublicKey,
			err,
		)
		// If no responders it means the account does not exist yet... ???
		// if !errors.Is(err, nats.ErrNoResponders) {
		// }
	}
	// If the account exists, do not re-create it because it contains imports from the accounts that.
	//
	// Please find a more elegant way of doing this.
	if string(lookupMessage.Data) != "" {
		slog.Info("found existing actor account", "account_id", t.Auth.ActorAccountPublicKey)
		return nil
	}
	slog.Info("creating actor account", "account_id", t.Auth.ActorAccountPublicKey)
	if _, err := nc.Request(
		"$SYS.REQ.CLAIMS.UPDATE",
		[]byte(t.Auth.ActorAccountJWT),
		time.Second,
	); err != nil {
		return fmt.Errorf("requesting new account: %w", err)
	}
	return nil
}

func UserJWTOption(userJWT string, userKeyPair nkeys.KeyPair) nats.Option {
	return nats.UserJWT(
		func() (string, error) {
			return userJWT, nil
		},
		func(bytes []byte) ([]byte, error) {
			return userKeyPair.Sign(bytes)
		},
	)
}

func NewTestServer(opts ...ServerOption) (TestServer, error) {
	sOpt := serverOptionDefaults
	for _, opt := range opts {
		opt(&sOpt)
	}
	if sOpt.Dir == "" {
		return TestServer{}, fmt.Errorf("dir is required")
	}
	if sOpt.FindAvailablePort {
		port, err := findAvailablePort()
		if err != nil {
			return TestServer{}, fmt.Errorf("finding available port: %w", err)
		}
		sOpt.Port = port
	}

	var jwtAuth ServerJWTAuth
	if sOpt.UseJWTAuth {
		jwtAuth = sOpt.JWTAuth
	} else {
		var err error
		jwtAuth, err = BootstrapServerJWTAuth()
		if err != nil {
			return TestServer{}, fmt.Errorf("bootstrapping server JWT auth: %w", err)
		}
	}

	// dirAccResolver, err := server.NewDirAccResolver(t.TempDir(), 0, 2*time.Minute, server.NoDelete)
	// if err != nil {
	// 	return testServer{}, fmt.Errorf("creating dirAccResolver: %w", err)
	// }
	// opts := &server.Options{
	// 	TrustedOperators: []*jwt.OperatorClaims{
	// 		oc,
	// 	},
	// 	SystemAccount:   sapk,
	// 	AccountResolver: dirAccResolver,

	// 	// JetStream: true,
	// }

	natsConf := generateNATSConf(jwtAuth, sOpt.Dir)
	natsConfFile := filepath.Join(sOpt.Dir, "resolver.conf")
	if err := os.MkdirAll(filepath.Dir(natsConfFile), os.ModePerm); err != nil {
		return TestServer{}, fmt.Errorf("creating dir for nats conf file: %w", err)
	}
	if err := os.WriteFile(natsConfFile, []byte(natsConf), os.ModePerm); err != nil {
		return TestServer{}, fmt.Errorf("writing nats conf: %w", err)
	}

	nsOpts := &server.Options{}
	if err := nsOpts.ProcessConfigFile(natsConfFile); err != nil {
		return TestServer{}, fmt.Errorf("processing config file: %w", err)
	}
	nsOpts.SystemAccount = jwtAuth.SysAccountPublicKey
	nsOpts.Debug = true
	nsOpts.Port = sOpt.Port
	nsOpts.JetStream = true
	nsOpts.StoreDir = sOpt.Dir

	if err := os.RemoveAll(natsConfFile); err != nil {
		return TestServer{}, fmt.Errorf("removing NATS conf file %s: %w", natsConfFile, err)
	}

	ns, err := server.NewServer(nsOpts)
	if err != nil {
		return TestServer{}, fmt.Errorf("new server: %w", err)
	}
	if sOpt.ConfigureLogger {
		ns.ConfigureLogger()
	}

	return TestServer{
		NS:   ns,
		Auth: jwtAuth,
	}, nil
}

type ServerOption func(*serverOption)

type serverOption struct {
	FindAvailablePort bool
	Port              int

	Dir string

	ConfigureLogger bool

	JWTAuth    ServerJWTAuth
	UseJWTAuth bool
}

var serverOptionDefaults = serverOption{
	FindAvailablePort: false,
	Port:              4222,
	Dir:               "",
}

func WithFindAvailablePort(enabled bool) ServerOption {
	return func(so *serverOption) {
		so.FindAvailablePort = enabled
	}
}

func WithPort(port int) ServerOption {
	return func(so *serverOption) {
		so.Port = port
	}
}

func WithDir(dir string) ServerOption {
	return func(so *serverOption) {
		so.Dir = dir
	}
}

func WithJWTAuth(jwtAuth ServerJWTAuth) ServerOption {
	return func(so *serverOption) {
		so.JWTAuth = jwtAuth
		so.UseJWTAuth = true
	}
}

func WithConfigureLogger(logger bool) ServerOption {
	return func(so *serverOption) {
		so.ConfigureLogger = true
	}
}

type ServerJWTAuth struct {
	// AccountServerURL is the configured `account_server_url` for the operator
	AccountServerURL string `json:"account_server_url"`

	OperatorNKey      []byte        `json:"operator_nkey"`
	OperatorKeyPair   nkeys.KeyPair `json:"-"`
	OperatorPublicKey string        `json:"operator_public_key"`
	OperatorJWT       string        `json:"operator_jwt"`

	SysAccountNKey      []byte        `json:"sys_account_nkey"`
	SysAccountKeyPair   nkeys.KeyPair `json:"-"`
	SysAccountPublicKey string        `json:"sys_account_public_key"`
	SysAccountJWT       string        `json:"sys_account_jwt"`

	SysUserNKey      []byte        `json:"sys_user_nkey"`
	SysUserKeyPair   nkeys.KeyPair `json:"-"`
	SysUserPublicKey string        `json:"sys_user_public_key"`
	SysUserJWT       string        `json:"sys_user_jwt"`

	ActorAccountNKey      []byte        `json:"actor_account_nkey"`
	ActorAccountKeyPair   nkeys.KeyPair `json:"-"`
	ActorAccountPublicKey string        `json:"actor_account_public_key"`
	ActorAccountJWT       string        `json:"actor_account_jwt"`

	ActorUserNKey      []byte        `json:"actor_user_nkey"`
	ActorUserKeyPair   nkeys.KeyPair `json:"-"`
	ActorUserPublicKey string        `json:"actor_user_public_key"`
	ActorUserJWT       string        `json:"actor_user_jwt"`
}

func (s *ServerJWTAuth) LoadKeyPairs() error {
	var err error
	s.OperatorKeyPair, err = nkeys.FromSeed(s.OperatorNKey)
	if err != nil {
		return fmt.Errorf("getting operator key pair from seed: %w", err)
	}
	s.SysAccountKeyPair, err = nkeys.FromSeed(s.SysAccountNKey)
	if err != nil {
		return fmt.Errorf("getting system account key pair from seed: %w", err)
	}
	s.SysUserKeyPair, err = nkeys.FromSeed(s.SysUserNKey)
	if err != nil {
		return fmt.Errorf("getting system user key pair from seed: %w", err)
	}
	s.ActorAccountKeyPair, err = nkeys.FromSeed(s.ActorAccountNKey)
	if err != nil {
		return fmt.Errorf("getting actor account key pair from seed: %w", err)
	}
	s.ActorUserKeyPair, err = nkeys.FromSeed(s.ActorUserNKey)
	if err != nil {
		return fmt.Errorf("getting actor user key pair from seed: %w", err)
	}
	return nil
}

func BootstrapServerJWTAuth() (ServerJWTAuth, error) {
	host := "0.0.0.0"
	port := 4222

	accountServerURL := fmt.Sprintf("nats://%s:%d", host, port)
	vr := jwt.ValidationResults{}

	// Create operator
	okp, err := nkeys.CreateOperator()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("creating operator key pair: %w", err)
	}
	oseed, err := okp.Seed()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("getting seed for operator: %w", err)
	}
	opk, err := okp.PublicKey()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("getting operator public key: %w", err)
	}
	// Create system account key pair
	sakp, err := nkeys.CreateAccount()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("creating system account key pair: %w", err)
	}
	sapk, err := sakp.PublicKey()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("getting system account public key: %w", err)
	}
	saseed, err := sakp.Seed()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("getting seed for system account: %w", err)
	}

	// Create operator JWT
	oc := jwt.NewOperatorClaims(opk)
	oc.Name = "test"
	oc.AccountServerURL = accountServerURL
	oc.SystemAccount = sapk
	oc.Validate(&vr)
	ojwt, err := oc.Encode(okp)
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("encoding operator JWT: %w", err)
	}

	// Create system account JWT
	sac := jwt.NewAccountClaims(sapk)
	sac.Name = "SYS"
	// Export some services to import in actor account
	sac.Exports.Add(&jwt.Export{
		Name: "account-monitoring-streams",
		// 	Subject:              "$SYS.ACCOUNT.*.>",
		Subject:              "$SYS.ACCOUNT.>",
		Type:                 jwt.Stream,
		AccountTokenPosition: 3,
		Info: jwt.Info{
			Description: "Account specific monitoring stream",
			InfoURL:     "https://docs.nats.io/nats-server/configuration/sys_accounts",
		},
	})
	sac.Exports.Add(&jwt.Export{
		Name: "account-monitoring-services",
		// 	Subject:              "$SYS.REQ.ACCOUNT.*.*",
		Subject:              "$SYS.REQ.ACCOUNT.>",
		Type:                 jwt.Service,
		AccountTokenPosition: 4,
		Info: jwt.Info{
			Description: "Request account specific monitoring services for: SUBSZ, CONNZ, LEAFZ, JSZ and INFO",
			InfoURL:     "https://docs.nats.io/nats-server/configuration/sys_accounts",
		},
	})
	sac.Validate(&vr)
	sajwt, err := sac.Encode(okp)
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("encoding system account JWT: %w", err)
	}

	// Create system user
	sukp, err := nkeys.CreateUser()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("creating system user: %w", err)
	}
	supk, err := sukp.PublicKey()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("getting system user public key: %w", err)
	}
	// Create system user JWT
	suc := jwt.NewUserClaims(supk)
	suc.Name = "sys"
	suc.IssuerAccount = sapk
	suc.Validate(&vr)

	sujwt, err := suc.Encode(sakp)
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("encoding system user JWT: %w", err)
	}

	suseed, err := sukp.Seed()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("getting seed for system user: %w", err)
	}

	//
	// Create Actor Account
	//
	aakp, err := nkeys.CreateAccount()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("creating actor account key pair: %w", err)
	}
	aapk, err := aakp.PublicKey()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("getting actor account public key: %w", err)
	}
	aaseed, err := aakp.Seed()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("getting seed for actor account: %w", err)
	}
	// Create actor account JWT
	aac := jwt.NewAccountClaims(aapk)
	aac.Name = "actor"
	aac.Limits.JetStreamLimits.Consumer = -1
	aac.Limits.JetStreamLimits.DiskMaxStreamBytes = -1
	aac.Limits.JetStreamLimits.DiskStorage = -1
	aac.Limits.JetStreamLimits.MaxAckPending = -1
	aac.Limits.JetStreamLimits.MemoryMaxStreamBytes = -1
	aac.Limits.JetStreamLimits.MemoryStorage = -1
	aac.Limits.JetStreamLimits.Streams = -1
	aac.Exports.Add(&jwt.Export{
		Type: jwt.Service,
		Name: "actor",
		Info: jwt.Info{
			Description: "Export actors for other account to import",
		},
		Subject: jwt.Subject(ActionExportSubject),
	})
	// aac.Exports.Add(&jwt.Export{
	// 	Type: jwt.Service,
	// 	Name: "ingest",
	// 	Info: jwt.Info{
	// 		Description: "Export ingesters for other account to import",
	// 	},
	// 	// Subject: ingest.<account>.<schema>.<version>.
	// 	Subject: jwt.Subject("ingest.*.*.*"),
	// })
	// Add stream exports according to what the `nsc` tool generates
	// aac.Exports.Add(&jwt.Export{
	// 	Name:                 "account-monitoring-streams",
	// 	Subject:              "$SYS.ACCOUNT.*.>",
	// 	Type:                 jwt.Stream,
	// 	AccountTokenPosition: 3,
	// 	Info: jwt.Info{
	// 		Description: "Account specific monitoring stream",
	// 		InfoURL:     "https://docs.nats.io/nats-server/configuration/sys_accounts",
	// 	},
	// })
	// aac.Exports.Add(&jwt.Export{
	// 	Name:                 "account-monitoring-services",
	// 	Subject:              "$SYS.REQ.ACCOUNT.*.*",
	// 	Type:                 jwt.Service,
	// 	AccountTokenPosition: 4,
	// 	Info: jwt.Info{
	// 		Description: "Request account specific monitoring services for: SUBSZ, CONNZ, LEAFZ, JSZ and INFO",
	// 		InfoURL:     "https://docs.nats.io/nats-server/configuration/sys_accounts",
	// 	},
	// })
	aac.Validate(&vr)
	aajwt, err := aac.Encode(okp)
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("encoding actor account JWT: %w", err)
	}

	//
	// Create Actor User
	//
	aukp, err := nkeys.CreateUser()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("creating actor user: %w", err)
	}
	aupk, err := aukp.PublicKey()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("getting actor user public key: %w", err)
	}
	auseed, err := aukp.Seed()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("getting seed for actor user: %w", err)
	}
	// Create actor user JWT
	auc := jwt.NewUserClaims(aupk)
	auc.Name = "actor"
	auc.IssuerAccount = aapk
	auc.Validate(&vr)
	aujwt, err := auc.Encode(aakp)
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("encoding actor user JWT: %w", err)
	}

	//
	// Create Actor Account
	//
	return ServerJWTAuth{
		AccountServerURL: accountServerURL,

		OperatorNKey:      oseed,
		OperatorKeyPair:   okp,
		OperatorPublicKey: opk,
		OperatorJWT:       ojwt,

		SysAccountNKey:      saseed,
		SysAccountKeyPair:   sakp,
		SysAccountPublicKey: sapk,
		SysAccountJWT:       sajwt,

		SysUserNKey:      suseed,
		SysUserKeyPair:   sukp,
		SysUserPublicKey: supk,
		SysUserJWT:       sujwt,

		ActorAccountNKey:      aaseed,
		ActorAccountKeyPair:   aakp,
		ActorAccountPublicKey: aapk,
		ActorAccountJWT:       aajwt,

		ActorUserNKey:      auseed,
		ActorUserKeyPair:   aukp,
		ActorUserPublicKey: aupk,
		ActorUserJWT:       aujwt,
	}, nil
}

func generateNATSConf(ja ServerJWTAuth, dir string) string {
	const natsConf = `
operator: %s

# Configuration of the nats based resolver
resolver {
    type: full
    dir: '%s'
    allow_delete: true
}

resolver_preload: {
	%s: %s,
}
`
	return fmt.Sprintf(
		natsConf,
		ja.OperatorJWT,
		dir,
		ja.SysAccountPublicKey,
		ja.SysAccountJWT,
	)
}

func findAvailablePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return -1, fmt.Errorf("listen: %w", err)
	}
	l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}
