package ko_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/ko/pkg/build"
	"github.com/google/ko/pkg/publish"
)

func TestKo(t *testing.T) {
	ctx := context.TODO()

	plat, err := v1.ParsePlatform("linux/amd64")
	if err != nil {
		t.Fatal(err)
	}
	ref, err := name.ParseReference("cgr.dev/chainguard/static:latest")
	if err != nil {
		t.Fatal(err)
	}
	desc, err := remote.Get(ref, remote.WithContext(ctx), remote.WithPlatform(*plat))
	if err != nil {
		t.Fatal(err)
	}
	base, err := desc.Image()
	if err != nil {
		t.Fatal(err)
	}
	bi, err := build.NewGo(
		ctx,
		"",
		build.WithDisabledSBOM(),
		build.WithPlatforms(plat.String()),
		build.WithBaseImages(func(ctx context.Context, s string) (name.Reference, build.Result, error) {
			return ref, base, nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	result, err := bi.Build(ctx, "./ex")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result.Size())

	pi, err := publish.NewDaemon(
		func(s1, s2 string) string {
			return s1 + "-whatever-" + s2
		},
		[]string{},
	)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := pi.Publish(ctx, result, "github.com/jlarfors/whatever")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("REF: ", pr.String())
	spew.Dump(pr)
	fmt.Println(pr.Context().RegistryStr())
	fmt.Println(pr.Context().RepositoryStr())

}
