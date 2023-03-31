package terraclient

import (
	tfjson "github.com/hashicorp/terraform-json"
)

func parseTfPlan(plan *tfjson.Plan) *Plan {
	var drift Plan
	for _, change := range plan.ResourceChanges {
		for _, action := range change.Change.Actions {
			switch action {
			case tfjson.ActionCreate:
				drift.AddResources = append(drift.AddResources, change)
			case tfjson.ActionDelete:
				drift.DestroyResources = append(drift.DestroyResources, change)
			case tfjson.ActionUpdate:
				drift.ChangeResources = append(drift.ChangeResources, change)
			default:
				// We don't care about other actions for the summary
			}
		}
	}

	return &drift
}

type Plan struct {
	AddResources     []*tfjson.ResourceChange
	ChangeResources  []*tfjson.ResourceChange
	DestroyResources []*tfjson.ResourceChange
}

func (p *Plan) HasDrift() bool {
	if len(p.AddResources) == 0 && len(p.ChangeResources) == 0 && len(p.DestroyResources) == 0 {
		return false
	}

	return true
}
