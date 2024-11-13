package rules

import (
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// MonitorTagsRule checks if the tags include "playbook-ops" and a valid "product:" tag
type MonitorTagsRule struct {
	tflint.DefaultRule
}

// NewMonitorTagsRule creates a new rule instance
func NewMonitorTagsRule() *MonitorTagsRule {
	return &MonitorTagsRule{}
}

// Name returns the rule name
func (r *MonitorTagsRule) Name() string {
	return "monitor_tags_format"
}

// Enabled returns whether the rule is enabled by default
func (r *MonitorTagsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *MonitorTagsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Check verifies if the monitor tags include the necessary tags
func (r *MonitorTagsRule) Check(runner tflint.Runner) error {
	// Regex para validar a tag "product:" com qualquer valor após o prefixo
	productTagPattern := `^product:([a-zA-Z0-9_-]+)$`
	productTagRegex := regexp.MustCompile(productTagPattern)

	// Obtém os recursos datadog_monitor do código terraform
	resources, err := runner.GetResourceContent("datadog_monitor", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "tags"},
		},
	}, nil)
	if err != nil {
		return err
	}

	// Valida as tags de cada monitor
	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes["tags"]
		if !exists {
			continue
		}

		var tags []string
		err := runner.EvaluateExpr(attribute.Expr, &tags, nil)
		if err != nil {
			return err
		}

		// Flags para verificar a presença das tags obrigatórias
		hasPlaybookOps := false
		hasValidProductTag := false

		// Verifica cada tag presente
		for _, tag := range tags {
			if tag == "playbook-ops" {
				hasPlaybookOps = true
			}
			if productTagRegex.MatchString(tag) {
				hasValidProductTag = true
			}
		}

		// Emite um aviso se as tags obrigatórias estiverem ausentes
		if !hasPlaybookOps || !hasValidProductTag {
			return runner.EmitIssue(
				r,
				"As tags 'playbook-ops' e 'product:[nome_do_produto]' são obrigatórias.",
				attribute.Expr.Range(),
			)
		}
	}

	return nil
}
