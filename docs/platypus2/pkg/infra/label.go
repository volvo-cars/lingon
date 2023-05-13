// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package infra

import "github.com/volvo-cars/lingon/pkg/terra"

const (
	TagManagedBy      = "ManagedBy"
	TagManagedByValue = "Lingon"
	// TagName human-readable resource name. Note that the AWS Console UI displays the case-sensitive "Name" tag.
	TagName = "Name"
	// TagAppID is a tag specifying the application identifier, application using the resource.
	TagAppID = "app-id"
	// TagAppRole is a tag specifying the resource's technical function, e.g. webserver, database, etc.
	TagAppRole = "app-role"
	// TagPurpose  is a tag specifying the resource's business purpose, e.g. "frontend ui", "payment processor", etc.
	TagPurpose = "purpose"
	// TagEnv is a tag specifying the environment.
	TagEnv = "environment"
	// TagProject is a tag specifying the project.
	TagProject = "project"
	// TagOwner is a tag specifying the person of contact.
	TagOwner = "owner"
	// TagCostCenter is a tag specifying the cost center that will receive the bill.
	TagCostCenter = "cost-center"
	// TagAutomationExclude is a tag specifying if the resource should be excluded from automation.
	// Value: true/false
	TagAutomationExclude = "automation-exclude"
	// TagPII is a tag specifying if the resource contains Personally Identifiable Information.
	// Value: true/false
	TagPII = "pii"
)

func Stags(ss ...string) terra.MapValue[terra.StringValue] {
	sv := make(map[string]terra.StringValue, 0)
	for i := 0; i < len(ss); i += 2 {
		if i+1 >= len(ss) {
			panic("odd number of strings")
		}
		sv[ss[i]] = S(ss[i+1])
	}
	sv[TagManagedBy] = S(TagManagedByValue)

	return terra.Map(sv)
}

func Ttags(m map[string]string) terra.MapValue[terra.StringValue] {
	sv := make(map[string]terra.StringValue, 0)
	for k, v := range m {
		sv[k] = S(v)
	}
	sv[TagManagedBy] = S(TagManagedByValue)
	return terra.Map(sv)
}

func MergeTags(m ...map[string]string) terra.MapValue[terra.StringValue] {
	sv := make(map[string]terra.StringValue, 0)
	for _, mm := range m {
		for k, v := range mm {
			sv[k] = S(v)
		}
	}
	sv[TagManagedBy] = S(TagManagedByValue)
	return terra.Map(sv)
}

func MergeSTags(
	m map[string]string,
	ss ...string,
) terra.MapValue[terra.StringValue] {
	sv := make(map[string]terra.StringValue, 0)
	for k, v := range m {
		sv[k] = S(v)
	}
	for i := 0; i < len(ss); i += 2 {
		if i+1 >= len(ss) {
			sv[ss[i]] = S("")
			break
		}
		sv[ss[i]] = S(ss[i+1])
	}
	sv[TagManagedBy] = S(TagManagedByValue)
	return terra.Map(sv)
}
