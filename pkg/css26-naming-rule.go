package pkg

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type Csc26NamingRule struct {
	tflint.DefaultRule
}

func NewCsc26NamingRule() *Csc26NamingRule {
	return &Csc26NamingRule{}
}

func (rule *Csc26NamingRule) Name() string {
	return "csc26_naming_convention"
}

func (rule *Csc26NamingRule) Enabled() bool {
	return true
}

func (rule *Csc26NamingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (rule *Csc26NamingRule) Link() string {
	return "https://confluence.office.ottonova.de/display/IT/CSC-26+-+Unified+naming+of+ressources"
}

func (rule *Csc26NamingRule) Check(runner tflint.Runner) error {
	// Define your module naming pattern here
	// Example: snake_case with specific prefixes
	pattern := regexp.MustCompile(`^[a-z][a-z0-9_]*[a-z0-9]$`)

	// Get all module calls
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "module",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "source"},
						{Name: "version"},
						{Name: "organizational_unit"},
						{Name: "environment"},
						{Name: "name"},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] linting %d blocks", len(content.Blocks))
	for _, block := range content.Blocks {
		if block.Type == "module" && len(block.Labels) > 0 {

			moduleName := block.Labels[0]
			log.Printf("[DEBUG] checking '%s'", moduleName)

			sourceAttr, ok := block.Body.Attributes["source"]
			var sourceVal string
			if ok {
				err := runner.EvaluateExpr(sourceAttr.Expr, &sourceVal, nil)
				if err != nil {
					return err
				}
			}

			// check if source is company module
			isCompanyModuleRegex := regexp.MustCompile(`^gitlab\.on\.ag`)
			if isCompanyModuleRegex.MatchString(sourceVal) {
				organizationalUnitAttr, ok := block.Body.Attributes["organizational_unit"]
				var organizationalUnit string
				if ok {
					err := runner.EvaluateExpr(organizationalUnitAttr.Expr, &organizationalUnit, nil)
					if err != nil {
						log.Printf("Error with ou")
						return err
					}
				}

				environmentAttr, ok := block.Body.Attributes["environment"]
				var environment string
				if ok {
					err := runner.EvaluateExpr(environmentAttr.Expr, &environment, nil)
					if err != nil {
						return err
					}
				}
				nameAttr, ok := block.Body.Attributes["name"]
				var name string
				if ok {
					err := runner.EvaluateExpr(nameAttr.Expr, &name, nil)
					if err != nil {
						return err
					}
					// check if module name is consistent with source and name attrs

					sourceSlice := strings.Split(sourceVal, "/")
					_ = sourceSlice[len(sourceSlice)-1]

				} else {
					if err := runner.EmitIssue(rule, fmt.Sprintf("Module '%s' must define a 'name' attribute", moduleName), block.DefRange); err != nil {
						return err
					}
				}

			}

		}
	}
	resources, err := runner.GetResourceContent("", nil, nil)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] linting %d resource blocks", len(resources.Blocks))
	for _, block := range resources.Blocks {
		if len(block.Labels) == 0 {
			continue
		}

		resourceType := block.Type
		resourceName := block.Labels[0]

		log.Printf("[DEBUG] checking resource '%s.%s'", resourceType, resourceName)

		if !pattern.MatchString(resourceName) {
			if err := runner.EmitIssue(
				rule,
				fmt.Sprintf("Resource name '%s' for type '%s' doesn't follow naming convention (expected: snake_case)", resourceName, resourceType),
				block.DefRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
