// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package infra

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"
)

// ValidateName validates the bucket name.
// https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html
func ValidateName(bucketName string) error {
	var errs []error
	// Bucket names must be between 3 (min) and 63 (max) characters long.
	if err := validateLength(bucketName); err != nil {
		errs = append(errs, err)
	}
	// Bucket names can consist only of lowercase letters, numbers, dots (.), and hyphens (-).
	if err := isOnlyLowercaseLettersNumbersDotsOrHyphens(bucketName); err != nil {
		errs = append(errs, err)
	}
	// Bucket names must begin and end with a letter or number.
	if err := beginOrEndWithLetterOrNumber(bucketName); err != nil {
		errs = append(errs, err)
	}
	// Bucket names must not contain two adjacent periods.
	if err := hasAdjacentPeriods(bucketName); err != nil {
		errs = append(errs, err)
	}
	// Bucket names must not be formatted as an IP address (for example, 192.168.5.4).
	if err := isIPAddress(bucketName); err != nil {
		errs = append(errs, err)
	}
	// Bucket names must not start with the prefix xn--.
	if strings.HasPrefix(bucketName, "xn--") {
		errs = append(
			errs,
			fmt.Errorf("bucket name cannot start with the prefix xn--"),
		)
	}
	// Bucket names must not end with the suffix -s3alias. This suffix is reserved for access point alias names.
	// For more information, see Using a bucket-style alias for your S3 bucket access point.
	if strings.HasSuffix(bucketName, "-s3alias") {
		errs = append(
			errs,
			fmt.Errorf("bucket name cannot end with the suffix -s3alias"),
		)
	}
	// Bucket names must be unique across all AWS accounts in all the AWS Regions within a partition.
	// A partition is a grouping of Regions.
	// AWS currently has three partitions: aws (Standard Regions), aws-cn (China Regions), and aws-us-gov (AWS GovCloud (US)).
	// -> CANNOT BE TESTED

	// A bucket name cannot be used by another AWS account in the same partition until the bucket is deleted.
	// -> CANNOT BE TESTED

	// Buckets used with Amazon S3 Transfer Acceleration can't have dots (.) in their names.
	// For more information about Transfer Acceleration, see Configuring fast, secure file transfers using Amazon S3 Transfer Acceleration.
	// -> EDGE CASE TOO FAR OUT OF SCOPE

	if len(errs) == 0 {
		return nil
	}

	var buf bytes.Buffer

	if len(errs) > 1 {
		_, _ = fmt.Fprintf(&buf, "%d errors: ", len(errs))
	}
	for i, err := range errs {
		if i != 0 {
			buf.WriteString("; ")
		}
		buf.WriteString(fmt.Sprintf("%d: %s", i+1, err.Error()))
	}

	return errors.New(buf.String())
}

func isOnlyLowercaseLettersNumbersDotsOrHyphens(name string) error {
	for _, c := range name {
		if !isLowercaseLetterOrNumber(c) && c != '.' && c != '-' {
			return fmt.Errorf(
				"name must only contain letters, numbers, dots, or hyphens",
			)
		}
	}
	return nil
}

func isIPAddress(name string) error {
	if r := net.ParseIP(name); r != nil {
		return fmt.Errorf("name must not be formatted as an IP address")
	}
	return nil
}

func validateLength(name string) error {
	length := len(name)
	if length < 3 || length > 63 {
		return fmt.Errorf(
			"name must be between 3 and 63 characters long: length = %d",
			length,
		)
	}
	return nil
}

func beginOrEndWithLetterOrNumber(name string) error {
	// get the first rune
	rname := []rune(name)
	length := len(rname)
	if !isLowercaseLetterOrNumber(rname[0]) || !isLowercaseLetterOrNumber(rname[length-1]) {
		return fmt.Errorf("name must begin and end with a lowercase letter or number")
	}
	return nil
}

func hasAdjacentPeriods(name string) error {
	length := len(name)
	for i := 0; i < length-1; i++ {
		if name[i] == '.' && name[i+1] == '.' {
			return fmt.Errorf("bucket name must not contain two adjacent periods")
		}
	}
	return nil
}

func isLowercaseLetterOrNumber(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z':
		return true
	case r >= '0' && r <= '9':
		return true
	default:
		return false
	}
}
