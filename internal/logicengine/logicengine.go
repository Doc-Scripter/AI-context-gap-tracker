package logicengine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cliffordotieno/ai-context-gap-tracker/internal/database"
)

// LogicEngine manages rule evaluation and logical consistency
type LogicEngine struct {
	db *database.DB
}

// Rule represents a logical rule
type Rule struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	RuleType    string                 `json:"rule_type"`
	Conditions  map[string]interface{} `json:"conditions"`
	Actions     map[string]interface{} `json:"actions"`
	Priority    int                    `json:"priority"`
	IsActive    bool                   `json:"is_active"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// RuleResult represents the result of rule evaluation
type RuleResult struct {
	RuleID      int         `json:"rule_id"`
	RuleName    string      `json:"rule_name"`
	Matched     bool        `json:"matched"`
	Confidence  float64     `json:"confidence"`
	Actions     []Action    `json:"actions"`
	Violations  []Violation `json:"violations"`
	Suggestions []string    `json:"suggestions"`
}

// Action represents an action to be taken
type Action struct {
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
}

// Violation represents a rule violation
type Violation struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Severity    string  `json:"severity"`
	Confidence  float64 `json:"confidence"`
}

// EvaluationContext contains context for rule evaluation
type EvaluationContext struct {
	SessionID   string                 `json:"session_id"`
	TurnNumber  int                    `json:"turn_number"`
	UserInput   string                 `json:"user_input"`
	Entities    map[string]interface{} `json:"entities"`
	Topics      []string               `json:"topics"`
	Timeline    []interface{}          `json:"timeline"`
	Assertions  []interface{}          `json:"assertions"`
	Ambiguities []interface{}          `json:"ambiguities"`
	History     []interface{}          `json:"history"`
}

// New creates a new LogicEngine instance
func New(db *database.DB) *LogicEngine {
	return &LogicEngine{
		db: db,
	}
}

// EvaluateRules evaluates all active rules against the given context
func (le *LogicEngine) EvaluateRules(ctx context.Context, evalContext *EvaluationContext) ([]*RuleResult, error) {
	// Get all active rules
	rules, err := le.GetActiveRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active rules: %w", err)
	}

	var results []*RuleResult

	// Evaluate each rule
	for _, rule := range rules {
		result, err := le.evaluateRule(ctx, rule, evalContext)
		if err != nil {
			log.Printf("Warning: failed to evaluate rule %s: %v", rule.Name, err)
			continue
		}

		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

// GetActiveRules retrieves all active rules from the database
func (le *LogicEngine) GetActiveRules(ctx context.Context) ([]*Rule, error) {
	query := `
		SELECT id, name, description, rule_type, conditions, actions, priority, is_active, created_at, updated_at
		FROM rules
		WHERE is_active = true
		ORDER BY priority DESC, created_at ASC
	`

	rows, err := le.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query rules: %w", err)
	}
	defer rows.Close()

	var rules []*Rule
	for rows.Next() {
		var rule Rule
		var conditionsJSON, actionsJSON []byte

		err := rows.Scan(
			&rule.ID,
			&rule.Name,
			&rule.Description,
			&rule.RuleType,
			&conditionsJSON,
			&actionsJSON,
			&rule.Priority,
			&rule.IsActive,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan rule: %w", err)
		}

		// Parse JSON fields
		json.Unmarshal(conditionsJSON, &rule.Conditions)
		json.Unmarshal(actionsJSON, &rule.Actions)

		rules = append(rules, &rule)
	}

	return rules, nil
}

// CreateRule creates a new rule in the database
func (le *LogicEngine) CreateRule(ctx context.Context, rule *Rule) error {
	conditionsJSON, _ := json.Marshal(rule.Conditions)
	actionsJSON, _ := json.Marshal(rule.Actions)

	query := `
		INSERT INTO rules (name, description, rule_type, conditions, actions, priority, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := le.db.QueryRow(query, rule.Name, rule.Description, rule.RuleType,
		conditionsJSON, actionsJSON, rule.Priority, rule.IsActive).Scan(
		&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)

	return err
}

// UpdateRule updates an existing rule
func (le *LogicEngine) UpdateRule(ctx context.Context, rule *Rule) error {
	conditionsJSON, _ := json.Marshal(rule.Conditions)
	actionsJSON, _ := json.Marshal(rule.Actions)

	query := `
		UPDATE rules
		SET name = $2, description = $3, rule_type = $4, conditions = $5, actions = $6, 
		    priority = $7, is_active = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := le.db.Exec(query, rule.ID, rule.Name, rule.Description, rule.RuleType,
		conditionsJSON, actionsJSON, rule.Priority, rule.IsActive)

	return err
}

// DeleteRule deletes a rule from the database
func (le *LogicEngine) DeleteRule(ctx context.Context, ruleID int) error {
	query := `DELETE FROM rules WHERE id = $1`
	_, err := le.db.Exec(query, ruleID)
	return err
}

// evaluateRule evaluates a single rule against the context
func (le *LogicEngine) evaluateRule(ctx context.Context, rule *Rule, evalContext *EvaluationContext) (*RuleResult, error) {
	// Evaluate based on rule type
	switch rule.RuleType {
	case "temporal_consistency":
		return le.evaluateTemporalConsistency(rule, evalContext)
	case "scope_agreement":
		return le.evaluateScopeAgreement(rule, evalContext)
	case "missing_information":
		return le.evaluateMissingInformation(rule, evalContext)
	case "contradiction_detection":
		return le.evaluateContradictionDetection(rule, evalContext)
	case "ambiguity_resolution":
		return le.evaluateAmbiguityResolution(rule, evalContext)
	default:
		return le.evaluateGenericRule(rule, evalContext)
	}
}

// evaluateTemporalConsistency checks for temporal consistency violations
func (le *LogicEngine) evaluateTemporalConsistency(rule *Rule, evalContext *EvaluationContext) (*RuleResult, error) {
	result := &RuleResult{
		RuleID:      rule.ID,
		RuleName:    rule.Name,
		Matched:     false,
		Confidence:  0.8,
		Actions:     []Action{},
		Violations:  []Violation{},
		Suggestions: []string{},
	}

	// Check for temporal keywords in user input
	temporalKeywords := []string{"yesterday", "tomorrow", "next week", "last month", "ago", "later"}
	userInput := strings.ToLower(evalContext.UserInput)

	for _, keyword := range temporalKeywords {
		if strings.Contains(userInput, keyword) {
			result.Matched = true
			result.Actions = append(result.Actions, Action{
				Type: "temporal_check",
				Parameters: map[string]interface{}{
					"keyword": keyword,
					"context": "temporal_reference_detected",
				},
			})
		}
	}

	// Check for timeline inconsistencies
	if len(evalContext.Timeline) > 1 {
		result.Suggestions = append(result.Suggestions, "Verify temporal sequence consistency")
	}

	return result, nil
}

// evaluateScopeAgreement checks for scope agreement violations
func (le *LogicEngine) evaluateScopeAgreement(rule *Rule, evalContext *EvaluationContext) (*RuleResult, error) {
	result := &RuleResult{
		RuleID:      rule.ID,
		RuleName:    rule.Name,
		Matched:     false,
		Confidence:  0.7,
		Actions:     []Action{},
		Violations:  []Violation{},
		Suggestions: []string{},
	}

	// Check for scope-related keywords
	scopeKeywords := []string{"all", "every", "some", "none", "most", "few"}
	userInput := strings.ToLower(evalContext.UserInput)

	for _, keyword := range scopeKeywords {
		if strings.Contains(userInput, keyword) {
			result.Matched = true
			result.Actions = append(result.Actions, Action{
				Type: "scope_clarification",
				Parameters: map[string]interface{}{
					"keyword": keyword,
					"context": "scope_quantifier_detected",
				},
			})
		}
	}

	return result, nil
}

// evaluateMissingInformation checks for missing information
func (le *LogicEngine) evaluateMissingInformation(rule *Rule, evalContext *EvaluationContext) (*RuleResult, error) {
	result := &RuleResult{
		RuleID:      rule.ID,
		RuleName:    rule.Name,
		Matched:     false,
		Confidence:  0.9,
		Actions:     []Action{},
		Violations:  []Violation{},
		Suggestions: []string{},
	}

	// Check for vague references
	vagueKeywords := []string{"it", "that", "this", "there", "place", "thing"}
	userInput := strings.ToLower(evalContext.UserInput)

	for _, keyword := range vagueKeywords {
		if strings.Contains(userInput, keyword) {
			result.Matched = true
			result.Violations = append(result.Violations, Violation{
				Type:        "vague_reference",
				Description: fmt.Sprintf("Vague reference detected: '%s'", keyword),
				Severity:    "medium",
				Confidence:  0.8,
			})
			result.Suggestions = append(result.Suggestions, fmt.Sprintf("Clarify what '%s' refers to", keyword))
		}
	}

	// Check for incomplete entities
	if len(evalContext.Entities) == 0 && len(evalContext.UserInput) > 10 {
		result.Matched = true
		result.Violations = append(result.Violations, Violation{
			Type:        "missing_entities",
			Description: "No entities detected in substantial input",
			Severity:    "low",
			Confidence:  0.6,
		})
	}

	return result, nil
}

// evaluateContradictionDetection checks for contradictions
func (le *LogicEngine) evaluateContradictionDetection(rule *Rule, evalContext *EvaluationContext) (*RuleResult, error) {
	result := &RuleResult{
		RuleID:      rule.ID,
		RuleName:    rule.Name,
		Matched:     false,
		Confidence:  0.8,
		Actions:     []Action{},
		Violations:  []Violation{},
		Suggestions: []string{},
	}

	// Check for contradictory keywords
	contradictoryPairs := [][]string{
		{"yes", "no"},
		{"always", "never"},
		{"all", "none"},
		{"before", "after"},
		{"increase", "decrease"},
	}

	userInput := strings.ToLower(evalContext.UserInput)

	for _, pair := range contradictoryPairs {
		if strings.Contains(userInput, pair[0]) && strings.Contains(userInput, pair[1]) {
			result.Matched = true
			result.Violations = append(result.Violations, Violation{
				Type:        "contradiction",
				Description: fmt.Sprintf("Potential contradiction detected: '%s' and '%s'", pair[0], pair[1]),
				Severity:    "high",
				Confidence:  0.7,
			})
			result.Suggestions = append(result.Suggestions, fmt.Sprintf("Clarify the relationship between '%s' and '%s'", pair[0], pair[1]))
		}
	}

	return result, nil
}

// evaluateAmbiguityResolution checks for ambiguities
func (le *LogicEngine) evaluateAmbiguityResolution(rule *Rule, evalContext *EvaluationContext) (*RuleResult, error) {
	result := &RuleResult{
		RuleID:      rule.ID,
		RuleName:    rule.Name,
		Matched:     false,
		Confidence:  0.9,
		Actions:     []Action{},
		Violations:  []Violation{},
		Suggestions: []string{},
	}

	// Check for ambiguous pronouns
	ambiguousPronouns := []string{"he", "she", "it", "they", "them", "this", "that"}
	userInput := strings.ToLower(evalContext.UserInput)

	for _, pronoun := range ambiguousPronouns {
		if strings.Contains(userInput, pronoun) {
			result.Matched = true
			result.Actions = append(result.Actions, Action{
				Type: "ambiguity_resolution",
				Parameters: map[string]interface{}{
					"pronoun": pronoun,
					"context": "ambiguous_pronoun_detected",
				},
			})
			result.Suggestions = append(result.Suggestions, fmt.Sprintf("Clarify what '%s' refers to", pronoun))
		}
	}

	// Check existing ambiguities
	if len(evalContext.Ambiguities) > 0 {
		result.Matched = true
		result.Actions = append(result.Actions, Action{
			Type: "resolve_ambiguities",
			Parameters: map[string]interface{}{
				"count": len(evalContext.Ambiguities),
			},
		})
	}

	return result, nil
}

// evaluateGenericRule evaluates a generic rule
func (le *LogicEngine) evaluateGenericRule(rule *Rule, evalContext *EvaluationContext) (*RuleResult, error) {
	result := &RuleResult{
		RuleID:      rule.ID,
		RuleName:    rule.Name,
		Matched:     false,
		Confidence:  0.5,
		Actions:     []Action{},
		Violations:  []Violation{},
		Suggestions: []string{},
	}

	// Basic generic evaluation
	if len(evalContext.UserInput) > 0 {
		result.Matched = true
		result.Actions = append(result.Actions, Action{
			Type: "generic_processing",
			Parameters: map[string]interface{}{
				"input_length": len(evalContext.UserInput),
			},
		})
	}

	return result, nil
}

// InitializeDefaultRules creates default rules in the database
func (le *LogicEngine) InitializeDefaultRules(ctx context.Context) error {
	defaultRules := []*Rule{
		{
			Name:        "Temporal Consistency Check",
			Description: "Checks for temporal consistency in user input",
			RuleType:    "temporal_consistency",
			Conditions:  map[string]interface{}{"enabled": true},
			Actions:     map[string]interface{}{"type": "temporal_check"},
			Priority:    100,
			IsActive:    true,
		},
		{
			Name:        "Missing Information Detection",
			Description: "Detects missing or vague information in user input",
			RuleType:    "missing_information",
			Conditions:  map[string]interface{}{"enabled": true},
			Actions:     map[string]interface{}{"type": "clarification_request"},
			Priority:    90,
			IsActive:    true,
		},
		{
			Name:        "Contradiction Detection",
			Description: "Detects contradictions in user input",
			RuleType:    "contradiction_detection",
			Conditions:  map[string]interface{}{"enabled": true},
			Actions:     map[string]interface{}{"type": "contradiction_alert"},
			Priority:    95,
			IsActive:    true,
		},
		{
			Name:        "Ambiguity Resolution",
			Description: "Identifies and resolves ambiguities",
			RuleType:    "ambiguity_resolution",
			Conditions:  map[string]interface{}{"enabled": true},
			Actions:     map[string]interface{}{"type": "ambiguity_clarification"},
			Priority:    85,
			IsActive:    true,
		},
	}

	for _, rule := range defaultRules {
		if err := le.CreateRule(ctx, rule); err != nil {
			log.Printf("Warning: failed to create default rule %s: %v", rule.Name, err)
		}
	}

	return nil
}