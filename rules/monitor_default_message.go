package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// MessageFormatRule valida se a propriedade message segue o formato esperado
type MessageFormatRule struct {
	tflint.DefaultRule
}

// NewMessageFormatRule cria uma nova instância da regra
func NewMessageFormatRule() *MessageFormatRule {
	return &MessageFormatRule{}
}

// Name retorna o nome da regra
func (r *MessageFormatRule) Name() string {
	return "monitor_message_format"
}

// Enabled retorna se a regra está habilitada por padrão
func (r *MessageFormatRule) Enabled() bool {
	return true
}

// Severity retorna a gravidade da regra
func (r *MessageFormatRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Check verifica se a propriedade message segue o formato esperado
func (r *MessageFormatRule) Check(runner tflint.Runner) error {
	// Expressão regular para validar os elementos obrigatórios de 'message' de forma flexível
	requiredElementsPattern := `(?s).*{{#is_alert}}.*@opsgenie-.*{{/is_alert}}.*{{#is_recovery}}.*@opsgenie-.*{{/is_recovery}}.*`

	// Obtém os recursos datadog_monitor do código terraform
	resources, err := runner.GetResourceContent("datadog_monitor", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "message"},
		},
	}, nil)
	if err != nil {
		return err
	}

	// Valida o conteúdo do 'message' de cada monitor
	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes["message"]
		if !exists {
			continue
		}

		var message string
		err := runner.EvaluateExpr(attribute.Expr, &message, nil)
		if err != nil {
			return err
		}

		// Verifica se o 'message' contém {{#is_alert}} e {{/is_alert}}, {{#is_recovery}} e {{/is_recovery}},
		// e @opsgenie- em qualquer lugar do texto, e permite flexibilidade no restante do conteúdo
		if !regexp.MustCompile(requiredElementsPattern).MatchString(message) {
			return runner.EmitIssue(
				r,
				"O conteúdo de 'message' não segue o formato esperado. Certifique-se de incluir todos os elementos necessários. Consulte a documentação: https://oraculo.rdstation.com.br/estrutura/enablers/sre/verticals/observability/alertas-fora-do-horario.",
				attribute.Expr.Range(),
			)
		}

		// Valida que cada seção obrigatória tenha pelo menos um título
		requiredSections := []string{
			"Impacto no Negócio",
			"Descrição técnica do problema",
			"Links úteis",
			"Possíveis causas",
			"Acionáveis",
			"Integração AlertManager",
			"Integração para recuperação do AlertManager",
		}

		for _, section := range requiredSections {
			if !strings.Contains(message, section) {
				return runner.EmitIssue(
					r,
					fmt.Sprintf("A seção '%s' está ausente ou com formato incorreto no conteúdo de 'message'.", section),
					attribute.Expr.Range(),
				)
			}
		}
	}

	return nil
}
