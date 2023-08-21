// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package natsutil

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nkeys"
)

func StartTestServer(t *testing.T) TestServer {
	ts, err := NewTestServer(
		WithDir(t.TempDir()),
		WithFindAvailablePort(true),
	)
	if err != nil {
		t.Fatal("creating test server: ", err)
	}
	if err := ts.StartUntilReady(); err != nil {
		t.Fatal("starting test server: ", err)
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
	ns.ConfigureLogger()

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
}

func BootstrapServerJWTAuth() (ServerJWTAuth, error) {
	host := "0.0.0.0"
	port := 4222

	accountServerURL := fmt.Sprintf("nats://%s:%d", host, port)

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

	// Create system account
	sakp, err := nkeys.CreateAccount()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("creating system account key pair: %w", err)
	}
	saseed, err := sakp.Seed()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("getting seed for system account: %w", err)
	}
	sapk, err := sakp.PublicKey()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("getting system account public key: %w", err)
	}

	// Create operator JWT
	oc := jwt.NewOperatorClaims(opk)
	oc.Name = "test"
	oc.AccountServerURL = accountServerURL
	oc.SystemAccount = sapk
	ojwt, err := oc.Encode(okp)
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("encoding operator JWT: %w", err)
	}

	// Create system account JWT
	sac := jwt.NewAccountClaims(opk)
	sac.Name = "SYS"
	sac.Subject = sapk
	// Add stream exports according to what the `nsc` tool generates
	sac.Exports.Add(&jwt.Export{
		Name:                 "account-monitoring-streams",
		Subject:              "$SYS.ACCOUNT.*.>",
		Type:                 jwt.Stream,
		AccountTokenPosition: 3,
		Info: jwt.Info{
			Description: "Account specific monitoring stream",
			InfoURL:     "https://docs.nats.io/nats-server/configuration/sys_accounts",
		},
	})
	sac.Exports.Add(&jwt.Export{
		Name:                 "account-monitoring-services",
		Subject:              "$SYS.REQ.ACCOUNT.*.*",
		Type:                 jwt.Service,
		AccountTokenPosition: 4,
		Info: jwt.Info{
			Description: "Request account specific monitoring services for: SUBSZ, CONNZ, LEAFZ, JSZ and INFO",
			InfoURL:     "https://docs.nats.io/nats-server/configuration/sys_accounts",
		},
	})
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

	sujwt, err := suc.Encode(sakp)
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("encoding system user JWT: %w", err)
	}

	suseed, err := sukp.Seed()
	if err != nil {
		return ServerJWTAuth{}, fmt.Errorf("getting seed for system user: %w", err)
	}

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
	}, nil
}

func generateNATSConf(ja ServerJWTAuth, dir string) string {
	const natsConf = `
operator: %s

# configuration of the nats based resolver
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
