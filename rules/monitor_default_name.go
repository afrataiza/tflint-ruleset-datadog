package rules

import (
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// MonitorNameRule valida se o nome do monitor segue o padrão esperado.
type MonitorNameRule struct {
	tflint.DefaultRule
}

// NewMonitorNameRule cria uma nova instância da regra
func NewMonitorNameRule() *MonitorNameRule {
	return &MonitorNameRule{}
}

// Name retorna o nome da regra
func (r *MonitorNameRule) Name() string {
	return "monitor_name_format"
}

// Enabled indica se a regra está habilitada por padrão
func (r *MonitorNameRule) Enabled() bool {
	return true
}

// Severity define a severidade da regra
func (r *MonitorNameRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Check valida o formato do nome do monitor
func (r *MonitorNameRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("datadog_monitor", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "name"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes["name"]
		if !exists {
			continue
		}

		err := runner.EvaluateExpr(attribute.Expr, func(name string) error {
			// Regex para validar o padrão do nome
			pattern := `^\[P[0-4]\]\[[A-Z0-9]+\]\[[A-Z0-9]+\]\[[A-Z0-9]+\]\[(PRODUCTION|STAGING)\] .+`
			matched, err := regexp.MatchString(pattern, name)
			if err != nil {
				return err
			}
			if !matched {
				return runner.EmitIssue(
					r,
					fmt.Sprintf("Nome do monitor '%s' não segue o padrão [PRIORIDADE][PRODUTO][DOMÍNIO][TIME][AMBIENTE] Título do alerta", name),
					attribute.Expr.Range(),
				)
			}
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
