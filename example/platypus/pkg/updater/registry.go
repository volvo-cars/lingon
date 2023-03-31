package updater

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"sort"
)

func GetLatestVersion(repository string, constraint string) (string, error) {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return "", fmt.Errorf("parsing constraint: %w", err)
	}

	repo, err := name.NewRepository(repository)
	if err != nil {
		return "", fmt.Errorf("parsing repository: %w", err)
	}
	remoteVersions, err := remote.List(repo)
	if err != nil {
		return "", fmt.Errorf("listing repository versions: %w", err)
	}

	semVersions := make([]*semver.Version, 0)
	for _, v := range remoteVersions {
		sv, err := semver.NewVersion(v)
		if err != nil {
			continue
		}
		if c.Check(sv) {
			semVersions = append(semVersions, sv)
		}
	}
	sort.Sort(semver.Collection(semVersions))

	if len(semVersions) == 0 {
		return "", errors.New("no new versions")
	}
	return semVersions[len(semVersions)-1].String(), nil
}
