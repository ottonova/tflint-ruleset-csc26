package main

import (
	"log"

	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"gitlab.on.ag/shared/software/tflint-ruleset-csc26/pkg"
)

// Plugin version
var Version string = "0.1.0"

func main() {
	log.Printf("[DEBUG] Starting CSC-26 lint")
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "csc26",
			Version: Version,
			Rules: []tflint.Rule{
				pkg.NewCsc26NamingRule(),
				pkg.NewCsc26EnforceIORule(),
			},
		},
	})
}
