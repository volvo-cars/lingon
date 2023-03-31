package updater

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetLatestVersion(t *testing.T) {
	lv, err := GetLatestVersion("grafana/grafana", ">9.3.6, <9.4")
	require.NoError(t, err)
	fmt.Println(lv)
}
