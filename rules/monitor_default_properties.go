package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// MonitorPropertiesRule valida valores de propriedades específicas do monitor
type MonitorPropertiesRule struct {
	tflint.DefaultRule
}

// NewMonitorPropertiesRule cria uma nova instância da regra
func NewMonitorPropertiesRule() *MonitorPropertiesRule {
	return &MonitorPropertiesRule{}
}

// Name retorna o nome da regra
func (r *MonitorPropertiesRule) Name() string {
	return "monitor_properties"
}

// Enabled indica se a regra está habilitada por padrão
func (r *MonitorPropertiesRule) Enabled() bool {
	return true
}

// Severity define a severidade da regra
func (r *MonitorPropertiesRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Check valida as propriedades renotify_interval, renotify_occurrences e timeout_h
func (r *MonitorPropertiesRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("datadog_monitor", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "renotify_interval"},
			{Name: "renotify_occurrences"},
			{Name: "timeout_h"},
		},
	}, nil)
	if err != nil {
		return err
	}

	checkAttribute := func(attrName string, expectedValue int, resource *hclext.Block) error {
		attribute, exists := resource.Body.Attributes[attrName]
		if !exists {
			return nil
		}

		return runner.EvaluateExpr(attribute.Expr, func(value int) error {
			if value != expectedValue {
				return runner.EmitIssue(
					r,
					fmt.Sprintf("Propriedade %s deve ser %d", attrName, expectedValue),
					attribute.Expr.Range(),
				)
			}
			return nil
		}, nil)
	}

	for _, resource := range resources.Blocks {
		if err := checkAttribute("renotify_interval", 60, resource); err != nil {
			return err
		}
		if err := checkAttribute("renotify_occurrences", 72, resource); err != nil {
			return err
		}
		if err := checkAttribute("timeout_h", 1, resource); err != nil {
			return err
		}
	}

	return nil
}
