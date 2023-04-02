package kube

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/parser"
	"github.com/zegl/kube-score/renderer/human"
	"github.com/zegl/kube-score/score"
	"github.com/zegl/kube-score/scorecard"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type inputReader struct {
	io.Reader
}

func (inputReader) Name() string {
	return "input"
}

func Score(r io.Reader) (*scorecard.Scorecard, error) {
	reader := &inputReader{
		Reader: r,
	}

	cnf := config.Configuration{
		AllFiles: []domain.NamedReader{reader},
		// EnabledOptionalTests: structMap,
		// KubernetesVersion: config.Semver{
		// 	Major: 1,
		// 	Minor: 23,
		// },
	}
	p, err := parser.New()
	if err != nil {
		return nil, err
	}
	parsed, err := p.ParseFiles(cnf)
	if err != nil {
		return nil, err
	}

	card, err := score.Score(parsed, cnf)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func RenderScoreCard(card *scorecard.Scorecard, w io.Writer, useColor bool) error {
	if !useColor {
		color.NoColor = true
	}
	output, err := human.Human(card, 0, 110)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, output); err != nil {
		return err
	}
	return nil
}

func podSpec(o runtime.Object) (*corev1.PodSpec, error) {
	switch t := o.(type) {
	case *corev1.Pod:
		return &t.Spec, nil
	case *corev1.ReplicationController:
		if t.Spec.Template != nil {
			return &t.Spec.Template.Spec, nil
		}
		return nil, fmt.Errorf("replication controller %s/%s has no template", t.Namespace, t.Name)
	case *appsv1.DaemonSet:
		return &t.Spec.Template.Spec, nil
	case *appsv1.Deployment:
		return &t.Spec.Template.Spec, nil
	case *batchv1.Job:
		return &t.Spec.Template.Spec, nil
	case *batchv1.CronJob:
		return &t.Spec.JobTemplate.Spec.Template.Spec, nil
	case *appsv1.StatefulSet:
		return &t.Spec.Template.Spec, nil
	default:
		return nil, fmt.Errorf("unknown object type %T", o)
	}
}

func podTemplateSpec(o runtime.Object) (*corev1.PodTemplateSpec, error) {
	switch t := o.(type) {
	case *corev1.ReplicationController:
		return t.Spec.Template, nil
	case *appsv1.DaemonSet:
		return &t.Spec.Template, nil
	case *appsv1.Deployment:
		return &t.Spec.Template, nil
	case *batchv1.Job:
		return &t.Spec.Template, nil
	case *batchv1.CronJob:
		return &t.Spec.JobTemplate.Spec.Template, nil
	case *appsv1.StatefulSet:
		return &t.Spec.Template, nil
	default:
		return nil, fmt.Errorf("unknown object type %T", o)
	}
}
