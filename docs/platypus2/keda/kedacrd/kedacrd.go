package kedacrd

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/volvo-cars/lingon/pkg/kube"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var _ kube.Exporter = (*CRD)(nil)

type CRD struct {
	kube.App
	ClustertriggerauthenticationsShCRD *apiextensionsv1.CustomResourceDefinition
	ScaledjobsShCRD                    *apiextensionsv1.CustomResourceDefinition
	ScaledobjectsShCRD                 *apiextensionsv1.CustomResourceDefinition
	TriggerauthenticationsShCRD        *apiextensionsv1.CustomResourceDefinition
}

func New() *CRD {
	return &CRD{
		ClustertriggerauthenticationsShCRD: ClustertriggerauthenticationsShCRD,
		ScaledjobsShCRD:                    ScaledjobsShCRD,
		ScaledobjectsShCRD:                 ScaledobjectsShCRD,
		TriggerauthenticationsShCRD:        TriggerauthenticationsShCRD,
	}
}

// Apply applies the kubernetes objects to the cluster
func (a *CRD) Apply(ctx context.Context) error {
	return Apply(ctx, a)
}

// Export exports the kubernetes objects to YAML files in the given directory
func (a *CRD) Export(dir string) error {
	return kube.Export(a, kube.WithExportOutputDirectory(dir))
}

// Apply applies the kubernetes objects contained in Exporter to the cluster
func Apply(ctx context.Context, km kube.Exporter) error {
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
	cmd.Env = os.Environ()        // inherit environment in case we need to use kubectl from a container
	stdin, err := cmd.StdinPipe() // pipe to pass data to kubectl
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		defer func() {
			err = errors.Join(err, stdin.Close())
		}()
		if errEW := kube.Export(
			km,
			kube.WithExportWriter(stdin),
			kube.WithExportAsSingleFile("stdin"),
		); errEW != nil {
			err = errors.Join(err, errEW)
		}
	}()

	if errS := cmd.Start(); errS != nil {
		return errors.Join(err, errS)
	}

	// waits for the command to exit and waits for any copying
	// to stdin or copying from stdout or stderr to complete
	return errors.Join(err, cmd.Wait())
}
