package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nats-io/jwt/v2"
	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kubeutil"
	"github.com/volvo-cars/nope/internal/natsutil"
	"golang.org/x/exp/slog"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	var (
		out   string
		force bool
	)

	flag.StringVar(&out, "out", "out", "output directory")
	flag.BoolVar(&force, "force", false, "force overwrite of output directory")
	flag.Parse()

	absOut, err := filepath.Abs(out)
	if err != nil {
		slog.Error("getting absolute path for out", "error", err)
		os.Exit(1)
	}

	if _, err := os.Stat(absOut); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			slog.Error("checking output directory", "error", err)
			os.Exit(1)
		}
		if err := os.MkdirAll(out, os.ModePerm); err != nil {
			slog.Error("creating output directory", "error", err)
			os.Exit(1)
		}
	} else {
		if !force {
			slog.Info("output directory already exists, use -force to overwrite")
			os.Exit(0)
		}
	}

	operatorNKeyFile := filepath.Join(absOut, "operator.nk")
	sysUserCredsFile := filepath.Join(absOut, "sys_user.creds")
	dotEnvFile := filepath.Join(absOut, ".env")
	kubeSecretFile := filepath.Join(absOut, "kube_secret.yaml")

	jwtAuth, err := natsutil.BootstrapServerJWTAuth()
	if err != nil {
		slog.Error("bootstrapping nats jwt auth", "error", err)
		os.Exit(1)
	}

	slog.Info("writing operator NKey", "file", operatorNKeyFile)
	if err := os.WriteFile(operatorNKeyFile, jwtAuth.OperatorNKey, 0o600); err != nil {
		slog.Error("writing operator nk: ", "error", err)
		os.Exit(1)
	}

	slog.Info("writing sys user creds", "file", sysUserCredsFile)
	userCreds, err := jwt.FormatUserConfig(jwtAuth.SysUserJWT, jwtAuth.SysUserNKey)
	if err != nil {
		slog.Error("formatting sys user creds: ", "error", err)
		os.Exit(1)
	}
	if err := os.WriteFile(sysUserCredsFile, userCreds, 0o600); err != nil {
		slog.Error("writing sys user creds: ", "error", err)
		os.Exit(1)
	}

	slog.Info("writing env file", "file", dotEnvFile)
	contents := fmt.Sprintf(
		"export NATS_URL=%s\nexport NATS_CREDS=%s\nexport NATS_OPERATOR_SEED=%s\n",
		jwtAuth.AccountServerURL,
		sysUserCredsFile,
		operatorNKeyFile,
	)
	if err := os.WriteFile(dotEnvFile, []byte(contents), 0o600); err != nil {
		slog.Error("writing env file: ", "error", err)
		os.Exit(1)
	}

	type natsOperatorSecret struct {
		kube.App
		Secret *corev1.Secret
	}

	app := natsOperatorSecret{
		Secret: &corev1.Secret{
			TypeMeta: kubeutil.TypeSecretV1,
			ObjectMeta: metav1.ObjectMeta{
				Name:      "nats-operator",
				Namespace: "system",
			},
			Data: map[string][]byte{
				"operator.nk":    jwtAuth.OperatorNKey,
				"sys_user.creds": userCreds,
			},
		},
	}
	var buffer bytes.Buffer
	if err := kube.Export(
		&app,
		kube.WithExportWriter(&buffer),
		kube.WithExportAsSingleFile("."),
	); err != nil {
		slog.Error("exporting kube secret", "error", err)
		os.Exit(1)
	}
	slog.Info("writing kube secret", "file", kubeSecretFile)
	if err := os.WriteFile(kubeSecretFile, buffer.Bytes(), 0644); err != nil {
		slog.Error("writing kube secret", "error", err)
		os.Exit(1)
	}
}
