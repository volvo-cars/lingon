package providerschemas

//go:generate go run -mod=readonly github.com/volvo-cars/go-terriyaki/cmd/tools/filtersb -out . -provider aws=hashicorp/aws:4.49.0 -include-resources aws_iam_role -include-data-sources aws_iam_role
