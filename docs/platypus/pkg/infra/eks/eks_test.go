package eks

import (
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"

	aws "github.com/golingon/terraproviders/aws/4.60.0"

	"github.com/volvo-cars/lingon/pkg/terra"
)

func TestEKS(t *testing.T) {
	type awsStack struct {
		terra.Stack

		Provider *aws.Provider
		Cluster  `validate:"required"`
	}
	eks := NewEKSCluster(
		ClusterOpts{
			Name:    "test",
			Version: "1.24",
			VPCID:   "123456",
			PrivateSubnetIDs: [3]string{
				"a", "b", "c",
			},
		},
	)
	stack := awsStack{
		Provider: aws.NewProvider(aws.ProviderArgs{}),
		Cluster:  *eks,
	}
	err := terra.Export(&stack, terra.WithExportOutputDirectory("out"))
	tu.AssertNoError(t, err, "exporting stack")
}
