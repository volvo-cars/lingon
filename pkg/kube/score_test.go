package kube_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/fatih/color"
	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kube/testdata/go/tekton"
	"github.com/zegl/kube-score/renderer/human"
	"github.com/zegl/kube-score/scorecard"
)

func TestScore(t *testing.T) {
	f, err := os.Open("testdata/tekton.yaml")
	if err != nil {
		t.Fatal(err)
	}

	card, err := kube.Score(f)
	// fmt.Printf("%##v", card)
	color.NoColor = false
	output, err := human.Human(card, 0, 110)
	if err != nil {
		t.Fatal(err)
	}

	m, ok := (*card).(map[string]scorecard.ScoredObject)
	if !ok {
		t.Fatal("not a map")
	}
	for file, object := range m {
		t.Log("file", file)
		for _, check := range object.Checks {
			_ = check.Comments
		}
	}
	t.Logf("%s\n", output)
}

func TestTxtar2Reader(t *testing.T) {
	tk := tekton.New()
	var buf bytes.Buffer
	if err := kube.Export(tk, kube.WithExportWriter(&buf)); err != nil {
		t.Fatal(err)
	}

	output, err := kube.Score(kube.Txtar2Reader(&buf))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s\n", output)
}
