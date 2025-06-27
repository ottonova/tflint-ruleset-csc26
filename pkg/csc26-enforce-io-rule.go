package pkg

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type Csc26EnforceIORule struct {
	tflint.DefaultRule
	requiredInputs  []string
	requiredOutputs []string
}

func NewCsc26EnforceIORule() *Csc26EnforceIORule {
	return &Csc26EnforceIORule{
		requiredInputs:  []string{"name", "environment", "organizational_unit"},
		requiredOutputs: []string{"name"},
	}
}

func (r *Csc26EnforceIORule) Name() string              { return "csc26_enforce_io" }
func (r *Csc26EnforceIORule) Enabled() bool             { return true }
func (r *Csc26EnforceIORule) Severity() tflint.Severity { return tflint.ERROR }
func (r *Csc26EnforceIORule) Link() string {
	return "https://confluence.office.ottonova.de/display/IT/CSC-26+-+Unified+naming+of+ressources"
}

func (r *Csc26EnforceIORule) Check(runner tflint.Runner) error {
	// 1) Collect all variable names
	varVars, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
			},
		},
	}, nil)
	if err != nil {
		return err
	}
	foundVars := map[string]bool{}
	for _, block := range varVars.Blocks {
		if len(block.Labels) > 0 {
			foundVars[block.Labels[0]] = true
		}
	}

	// 2) Collect all output names
	outVars, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "output",
				LabelNames: []string{"name"},
			},
		},
	}, nil)
	if err != nil {
		return err
	}
	foundOut := map[string]bool{}
	for _, block := range outVars.Blocks {
		if len(block.Labels) > 0 {
			foundOut[block.Labels[0]] = true
		}
	}
	var exampleRange hcl.Range
	if len(varVars.Blocks) > 0 {
		exampleRange = varVars.Blocks[0].DefRange
	}
	// 3) Check required inputs
	for _, want := range r.requiredInputs {
		if !foundVars[want] {

			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("Module is missing required input variable: %q", want),
				exampleRange,
			); err != nil {
				return err
			}
		}
	}

	// 4) Check required outputs
	var outExampleRange hcl.Range
	if len(outVars.Blocks) > 0 {
		outExampleRange = outVars.Blocks[0].DefRange
	}

	for _, want := range r.requiredOutputs {
		if !foundOut[want] {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("Module is missing required output: %q", want),
				outExampleRange,
			); err != nil {
				return err
			}
		}
	}
	return nil
}
