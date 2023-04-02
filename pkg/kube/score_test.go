package kube_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/volvo-cars/lingon/pkg/kube"
	"github.com/volvo-cars/lingon/pkg/kube/testdata/go/tekton"
	"github.com/zegl/kube-score/scorecard"
)

func TestScore(t *testing.T) {
	f, err := os.Open("testdata/tekton.yaml")
	if err != nil {
		t.Fatal(err)
	}

	card, err := kube.Score(f)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Printf("%##v", card)
	// color.NoColor = true
	// output, err := human.Human(card, 0, 110)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	for _, object := range *card {
		name := object.ObjectMeta.Name
		ns := object.ObjectMeta.Namespace
		tv := object.TypeMeta.APIVersion
		tk := object.TypeMeta.Kind
		t.Log(tv + "/" + tk + "/" + ns + "/" + name)
		for _, check := range object.Checks {
			if check.Skipped {
				continue
			}
			printComments(t, check)
			// t.Log("file", file, "comment", check.Comments)
		}
	}
}

func printComments(
	t *testing.T,
	ts scorecard.TestScore,
	// comments []scorecard.TestScoreComment,
) {
	if ts.Grade == scorecard.GradeAllOK {
		return
	}
	grade := greade2Str(ts.Grade)
	for _, comment := range ts.Comments {
		t.Logf("%s \t [%s] %s - %s\n", grade, ts.Check.ID, comment.Summary, comment.Description)
	}
}

func greade2Str(grade scorecard.Grade) string {
	switch grade {
	case scorecard.GradeCritical:
		return "CRITICAL"
	case scorecard.GradeWarning:
		return "WARNING"
	case scorecard.GradeAlmostOK:
		return "ALMOST OK"
	case scorecard.GradeAllOK:
		return "ALL OK"
	default:
		return "UNKNOWN"
	}
}

func TestTxtar2Reader(t *testing.T) {
	tk := tekton.New()
	var buf bytes.Buffer
	if err := kube.Export(
		tk,
		kube.WithExportWriter(&buf),
		kube.WithExportAsSingleFile("input.yaml"),
	); err != nil {
		t.Fatal(err)
	}

	output, err := kube.Score(&buf)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v\n", output)
}
