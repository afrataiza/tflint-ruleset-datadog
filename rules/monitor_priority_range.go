package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// PriorityRangeRule verifica se a chave `priority` está entre 1 e 5.
type PriorityRangeRule struct {
	tflint.DefaultRule
}

// NewPriorityRangeRule cria uma nova regra para `priority`
func NewPriorityRangeRule() *PriorityRangeRule {
	return &PriorityRangeRule{}
}

// Name retorna o nome da regra
func (r *PriorityRangeRule) Name() string {
	return "monitor_priority_range"
}

// Enabled define se a regra é habilitada por padrão
func (r *PriorityRangeRule) Enabled() bool {
	return true
}

// Severity define a severidade da regra
func (r *PriorityRangeRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Check valida se `priority` está no intervalo entre 1 e 5
func (r *PriorityRangeRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("datadog_monitor", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "priority"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes["priority"]
		if !exists {
			continue
		}

		var priority int
		err := runner.EvaluateExpr(attribute.Expr, &priority, nil)
		if err != nil {
			return err
		}

		if priority < 1 || priority > 5 {
			return runner.EmitIssue(
				r,
				fmt.Sprintf("O valor da prioridade deve estar entre 1 e 5, mas é %d", priority),
				attribute.Expr.Range(),
			)
		}
	}

	return nil
}
