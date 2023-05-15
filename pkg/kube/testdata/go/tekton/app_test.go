package tekton

import "testing"

func TestTekton(t *testing.T) {
	if err := New().Export("out"); err != nil {
		t.Error(err)
	}
}
