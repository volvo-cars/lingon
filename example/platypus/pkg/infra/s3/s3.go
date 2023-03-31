package s3

import (
	"github.com/volvo-cars/lingon/example/platypus/gen/providers/aws"
	"github.com/volvo-cars/lingon/example/platypus/gen/providers/aws/s3bucketserversideencryptionconfiguration"
	"github.com/volvo-cars/lingon/example/platypus/gen/providers/aws/s3bucketversioning"

	"github.com/volvo-cars/lingon/pkg/terra"
)

type Bucket struct {
	S3           *aws.S3Bucket                                  `validate:"required"`
	ACL          *aws.S3BucketAcl                               `validate:"required"`
	Versioning   *aws.S3BucketVersioning                        `validate:"required"`
	PublicAccess *aws.S3BucketPublicAccessBlock                 `validate:"required"`
	SSE          *aws.S3BucketServerSideEncryptionConfiguration `validate:"required"`
}

func NewBucket(bucketName string) *Bucket {
	b := aws.NewS3Bucket(
		"s3", aws.S3BucketArgs{
			Bucket: terra.String(bucketName),
			Tags: terra.Map(
				map[string]terra.StringValue{
					"Name": terra.String("DataOps TF EKS Experiment"),
				},
			),
		},
	)

	bucketID := b.Attributes().Id()

	acl := aws.NewS3BucketAcl(
		"s3", aws.S3BucketAclArgs{
			Bucket: bucketID,
			Acl:    terra.String("private"),
		},
	)

	vv := aws.NewS3BucketVersioning(
		"s3", aws.S3BucketVersioningArgs{
			Bucket: bucketID,
			VersioningConfiguration: &s3bucketversioning.VersioningConfiguration{
				Status: terra.String("Enabled"),
			},
		},
	)

	pab := aws.NewS3BucketPublicAccessBlock(
		"s3", aws.S3BucketPublicAccessBlockArgs{
			Bucket:                bucketID,
			BlockPublicAcls:       terra.Bool(true),
			BlockPublicPolicy:     terra.Bool(true),
			IgnorePublicAcls:      terra.Bool(true),
			RestrictPublicBuckets: terra.Bool(true),
		},
	)

	enc := aws.NewS3BucketServerSideEncryptionConfiguration(
		"s3", aws.S3BucketServerSideEncryptionConfigurationArgs{
			Bucket: bucketID,
			Rule:   RuleEncryptKMS(),
		},
	)

	return &Bucket{
		S3:           b,
		ACL:          acl,
		Versioning:   vv,
		PublicAccess: pab,
		SSE:          enc,
	}
}

func RuleEncryptKMS() []s3bucketserversideencryptionconfiguration.Rule {
	return []s3bucketserversideencryptionconfiguration.Rule{
		{
			ApplyServerSideEncryptionByDefault: &s3bucketserversideencryptionconfiguration.ApplyServerSideEncryptionByDefault{
				SseAlgorithm: terra.String("aws:kms"),
			},
		},
	}
}
