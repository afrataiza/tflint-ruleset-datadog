package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-template/rules"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "datadog",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewMonitorPropertiesRule(),
				rules.NewMonitorNameRule(),
				rules.NewPriorityRangeRule(),
				rules.NewMessageFormatRule(),
				rules.NewMonitorTagsRule(),
			},
		},
	})
}
