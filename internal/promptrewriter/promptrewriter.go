package promptrewriter

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cliffordotieno/ai-context-gap-tracker/internal/contexttracker"
	"github.com/cliffordotieno/ai-context-gap-tracker/internal/logicengine"
)

// PromptRewriter enhances prompts with context and clarity information
type PromptRewriter struct {
	contextTracker *contexttracker.ContextTracker
	logicEngine    *logicengine.LogicEngine
}

// RewriteRequest represents a request to rewrite a prompt
type RewriteRequest struct {
	SessionID    string                 `json:"session_id"`
	TurnNumber   int                    `json:"turn_number"`
	UserInput    string                 `json:"user_input"`
	Context      map[string]interface{} `json:"context"`
	SystemPrompt string                 `json:"system_prompt"`
	Options      RewriteOptions         `json:"options"`
}

// RewriteOptions configures the rewrite behavior
type RewriteOptions struct {
	IncludeContext       bool `json:"include_context"`
	IncludeAmbiguities   bool `json:"include_ambiguities"`
	IncludeAssumptions   bool `json:"include_assumptions"`
	IncludeHistory       bool `json:"include_history"`
	AddDisambiguation    bool `json:"add_disambiguation"`
	AddClarityFlags      bool `json:"add_clarity_flags"`
	MaxContextLength     int  `json:"max_context_length"`
	MaxHistoryTurns      int  `json:"max_history_turns"`
	OptimizeForClarity   bool `json:"optimize_for_clarity"`
	OptimizeForAccuracy  bool `json:"optimize_for_accuracy"`
}

// RewriteResult represents the result of prompt rewriting
type RewriteResult struct {
	OriginalPrompt   string                 `json:"original_prompt"`
	RewrittenPrompt  string                 `json:"rewritten_prompt"`
	Context          map[string]interface{} `json:"context"`
	Ambiguities      []string               `json:"ambiguities"`
	Assumptions      []string               `json:"assumptions"`
	ClarityFlags     []string               `json:"clarity_flags"`
	DisambiguationFlags []string            `json:"disambiguation_flags"`
	QualityScore     float64                `json:"quality_score"`
	Recommendations  []string               `json:"recommendations"`
	ProcessingTime   time.Duration          `json:"processing_time"`
}

// ClarityFlag represents a clarity flag in the prompt
type ClarityFlag struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Severity    string  `json:"severity"`
	Confidence  float64 `json:"confidence"`
}

// DisambiguationFlag represents a disambiguation flag
type DisambiguationFlag struct {
	Type         string   `json:"type"`
	AmbiguousItem string   `json:"ambiguous_item"`
	Suggestions  []string `json:"suggestions"`
	Confidence   float64  `json:"confidence"`
}

// DefaultRewriteOptions returns default rewrite options
func DefaultRewriteOptions() RewriteOptions {
	return RewriteOptions{
		IncludeContext:       true,
		IncludeAmbiguities:   true,
		IncludeAssumptions:   true,
		IncludeHistory:       true,
		AddDisambiguation:    true,
		AddClarityFlags:      true,
		MaxContextLength:     2000,
		MaxHistoryTurns:      5,
		OptimizeForClarity:   true,
		OptimizeForAccuracy:  true,
	}
}

// New creates a new PromptRewriter instance
func New(contextTracker *contexttracker.ContextTracker, logicEngine *logicengine.LogicEngine) *PromptRewriter {
	return &PromptRewriter{
		contextTracker: contextTracker,
		logicEngine:    logicEngine,
	}
}

// RewritePrompt rewrites a prompt with enhanced context and clarity information
func (pr *PromptRewriter) RewritePrompt(ctx context.Context, request *RewriteRequest) (*RewriteResult, error) {
	startTime := time.Now()

	result := &RewriteResult{
		OriginalPrompt:      request.UserInput,
		Context:             make(map[string]interface{}),
		Ambiguities:         []string{},
		Assumptions:         []string{},
		ClarityFlags:        []string{},
		DisambiguationFlags: []string{},
		QualityScore:        0.0,
		Recommendations:     []string{},
	}

	// Get context information
	contextInfo, err := pr.gatherContextInformation(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to gather context information: %w", err)
	}

	// Get rule evaluation results
	ruleResults, err := pr.evaluateRules(ctx, request, contextInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate rules: %w", err)
	}

	// Build the rewritten prompt
	rewrittenPrompt := pr.buildRewrittenPrompt(request, contextInfo, ruleResults)

	// Extract clarity and disambiguation information
	result.Context = contextInfo
	result.Ambiguities = pr.extractAmbiguities(contextInfo, ruleResults)
	result.Assumptions = pr.extractAssumptions(contextInfo, ruleResults)
	result.ClarityFlags = pr.extractClarityFlags(ruleResults)
	result.DisambiguationFlags = pr.extractDisambiguationFlags(ruleResults)

	// Calculate quality score
	result.QualityScore = pr.calculateQualityScore(result)

	// Generate recommendations
	result.Recommendations = pr.generateRecommendations(result)

	result.RewrittenPrompt = rewrittenPrompt
	result.ProcessingTime = time.Since(startTime)

	return result, nil
}

// gatherContextInformation gathers relevant context information
func (pr *PromptRewriter) gatherContextInformation(ctx context.Context, request *RewriteRequest) (map[string]interface{}, error) {
	contextInfo := make(map[string]interface{})

	// Get current context
	currentContext, err := pr.contextTracker.GetContext(ctx, request.SessionID, request.TurnNumber)
	if err == nil && currentContext != nil {
		contextInfo["current_turn"] = currentContext
		contextInfo["entities"] = currentContext.Entities
		contextInfo["topics"] = currentContext.Topics
		contextInfo["timeline"] = currentContext.Timeline
		contextInfo["assertions"] = currentContext.Assertions
		contextInfo["ambiguities"] = currentContext.Ambiguities
	}

	// Get session history if requested
	if request.Options.IncludeHistory {
		sessionContext, err := pr.contextTracker.GetSessionContext(ctx, request.SessionID)
		if err == nil {
			// Limit history to specified number of turns
			maxTurns := request.Options.MaxHistoryTurns
			if maxTurns > 0 && len(sessionContext) > maxTurns {
				sessionContext = sessionContext[len(sessionContext)-maxTurns:]
			}
			contextInfo["history"] = sessionContext
		}
	}

	// Get memory graph
	memoryGraph, err := pr.contextTracker.GetMemoryGraph(ctx, request.SessionID)
	if err == nil && memoryGraph != nil {
		contextInfo["memory_graph"] = memoryGraph
	}

	return contextInfo, nil
}

// evaluateRules evaluates logic rules against the context
func (pr *PromptRewriter) evaluateRules(ctx context.Context, request *RewriteRequest, contextInfo map[string]interface{}) ([]*logicengine.RuleResult, error) {
	// Create evaluation context
	evalContext := &logicengine.EvaluationContext{
		SessionID:  request.SessionID,
		TurnNumber: request.TurnNumber,
		UserInput:  request.UserInput,
		Entities:   make(map[string]interface{}),
		Topics:     []string{},
		Timeline:   []interface{}{},
		Assertions: []interface{}{},
		Ambiguities: []interface{}{},
		History:    []interface{}{},
	}

	// Populate evaluation context from gathered context info
	if entities, ok := contextInfo["entities"].(map[string]interface{}); ok {
		evalContext.Entities = entities
	}
	if topics, ok := contextInfo["topics"].([]string); ok {
		evalContext.Topics = topics
	}
	if timeline, ok := contextInfo["timeline"].([]interface{}); ok {
		evalContext.Timeline = timeline
	}
	if assertions, ok := contextInfo["assertions"].([]interface{}); ok {
		evalContext.Assertions = assertions
	}
	if ambiguities, ok := contextInfo["ambiguities"].([]interface{}); ok {
		evalContext.Ambiguities = ambiguities
	}
	if history, ok := contextInfo["history"].([]interface{}); ok {
		evalContext.History = history
	}

	// Evaluate rules
	return pr.logicEngine.EvaluateRules(ctx, evalContext)
}

// buildRewrittenPrompt constructs the enhanced prompt
func (pr *PromptRewriter) buildRewrittenPrompt(request *RewriteRequest, contextInfo map[string]interface{}, ruleResults []*logicengine.RuleResult) string {
	var promptBuilder strings.Builder

	// Start with system prompt if provided
	if request.SystemPrompt != "" {
		promptBuilder.WriteString(request.SystemPrompt)
		promptBuilder.WriteString("\n\n")
	}

	// Add context section
	if request.Options.IncludeContext {
		promptBuilder.WriteString("CONTEXT INFORMATION:\n")
		pr.addContextSection(&promptBuilder, contextInfo, request.Options)
		promptBuilder.WriteString("\n")
	}

	// Add disambiguation flags
	if request.Options.AddDisambiguation {
		disambiguationFlags := pr.extractDisambiguationFlags(ruleResults)
		if len(disambiguationFlags) > 0 {
			promptBuilder.WriteString("DISAMBIGUATION REQUIRED:\n")
			for _, flag := range disambiguationFlags {
				promptBuilder.WriteString(fmt.Sprintf("- %s\n", flag))
			}
			promptBuilder.WriteString("\n")
		}
	}

	// Add clarity flags
	if request.Options.AddClarityFlags {
		clarityFlags := pr.extractClarityFlags(ruleResults)
		if len(clarityFlags) > 0 {
			promptBuilder.WriteString("CLARITY CONSIDERATIONS:\n")
			for _, flag := range clarityFlags {
				promptBuilder.WriteString(fmt.Sprintf("- %s\n", flag))
			}
			promptBuilder.WriteString("\n")
		}
	}

	// Add ambiguities section
	if request.Options.IncludeAmbiguities {
		ambiguities := pr.extractAmbiguities(contextInfo, ruleResults)
		if len(ambiguities) > 0 {
			promptBuilder.WriteString("IDENTIFIED AMBIGUITIES:\n")
			for _, ambiguity := range ambiguities {
				promptBuilder.WriteString(fmt.Sprintf("- %s\n", ambiguity))
			}
			promptBuilder.WriteString("\n")
		}
	}

	// Add assumptions section
	if request.Options.IncludeAssumptions {
		assumptions := pr.extractAssumptions(contextInfo, ruleResults)
		if len(assumptions) > 0 {
			promptBuilder.WriteString("CURRENT ASSUMPTIONS:\n")
			for _, assumption := range assumptions {
				promptBuilder.WriteString(fmt.Sprintf("- %s\n", assumption))
			}
			promptBuilder.WriteString("\n")
		}
	}

	// Add user input
	promptBuilder.WriteString("USER INPUT:\n")
	promptBuilder.WriteString(request.UserInput)
	promptBuilder.WriteString("\n\n")

	// Add optimization instructions
	if request.Options.OptimizeForClarity {
		promptBuilder.WriteString("RESPONSE REQUIREMENTS:\n")
		promptBuilder.WriteString("- Provide clear, unambiguous responses\n")
		promptBuilder.WriteString("- State assumptions explicitly\n")
		promptBuilder.WriteString("- Ask for clarification when needed\n")
		
		if request.Options.OptimizeForAccuracy {
			promptBuilder.WriteString("- Verify information before stating facts\n")
			promptBuilder.WriteString("- Indicate confidence levels\n")
		}
		
		promptBuilder.WriteString("\n")
	}

	return promptBuilder.String()
}

// addContextSection adds context information to the prompt
func (pr *PromptRewriter) addContextSection(builder *strings.Builder, contextInfo map[string]interface{}, options RewriteOptions) {
	// Add entities
	if entities, ok := contextInfo["entities"].(map[string]interface{}); ok && len(entities) > 0 {
		builder.WriteString("Entities: ")
		entitiesJSON, _ := json.Marshal(entities)
		builder.WriteString(string(entitiesJSON))
		builder.WriteString("\n")
	}

	// Add topics
	if topics, ok := contextInfo["topics"].([]string); ok && len(topics) > 0 {
		builder.WriteString("Topics: ")
		topicsJSON, _ := json.Marshal(topics)
		builder.WriteString(string(topicsJSON))
		builder.WriteString("\n")
	}

	// Add timeline
	if timeline, ok := contextInfo["timeline"].([]interface{}); ok && len(timeline) > 0 {
		builder.WriteString("Timeline: ")
		timelineJSON, _ := json.Marshal(timeline)
		builder.WriteString(string(timelineJSON))
		builder.WriteString("\n")
	}

	// Add history (limited)
	if options.IncludeHistory {
		if history, ok := contextInfo["history"].([]interface{}); ok && len(history) > 0 {
			builder.WriteString("Recent History: ")
			historyJSON, _ := json.Marshal(history)
			historyStr := string(historyJSON)
			if len(historyStr) > options.MaxContextLength {
				historyStr = historyStr[:options.MaxContextLength] + "..."
			}
			builder.WriteString(historyStr)
			builder.WriteString("\n")
		}
	}
}

// extractAmbiguities extracts ambiguity information
func (pr *PromptRewriter) extractAmbiguities(contextInfo map[string]interface{}, ruleResults []*logicengine.RuleResult) []string {
	var ambiguities []string

	// From context info
	if ambiguitiesData, ok := contextInfo["ambiguities"].([]interface{}); ok {
		for _, ambiguity := range ambiguitiesData {
			if ambiguityStr, ok := ambiguity.(string); ok {
				ambiguities = append(ambiguities, ambiguityStr)
			}
		}
	}

	// From rule results
	for _, result := range ruleResults {
		if result.Matched {
			for _, suggestion := range result.Suggestions {
				if strings.Contains(strings.ToLower(suggestion), "ambiguous") ||
					strings.Contains(strings.ToLower(suggestion), "clarify") {
					ambiguities = append(ambiguities, suggestion)
				}
			}
		}
	}

	return ambiguities
}

// extractAssumptions extracts assumption information
func (pr *PromptRewriter) extractAssumptions(contextInfo map[string]interface{}, ruleResults []*logicengine.RuleResult) []string {
	var assumptions []string

	// From context info
	if assertionsData, ok := contextInfo["assertions"].([]interface{}); ok {
		for _, assertion := range assertionsData {
			if assertionStr, ok := assertion.(string); ok {
				assumptions = append(assumptions, assertionStr)
			}
		}
	}

	// From rule results
	for _, result := range ruleResults {
		if result.Matched {
			for _, suggestion := range result.Suggestions {
				if strings.Contains(strings.ToLower(suggestion), "assume") ||
					strings.Contains(strings.ToLower(suggestion), "presuming") {
					assumptions = append(assumptions, suggestion)
				}
			}
		}
	}

	return assumptions
}

// extractClarityFlags extracts clarity flags from rule results
func (pr *PromptRewriter) extractClarityFlags(ruleResults []*logicengine.RuleResult) []string {
	var flags []string

	for _, result := range ruleResults {
		if result.Matched {
			for _, action := range result.Actions {
				if action.Type == "clarification_request" || action.Type == "scope_clarification" {
					flags = append(flags, fmt.Sprintf("Clarification needed: %s", result.RuleName))
				}
			}
		}
	}

	return flags
}

// extractDisambiguationFlags extracts disambiguation flags from rule results
func (pr *PromptRewriter) extractDisambiguationFlags(ruleResults []*logicengine.RuleResult) []string {
	var flags []string

	for _, result := range ruleResults {
		if result.Matched {
			for _, action := range result.Actions {
				if action.Type == "ambiguity_resolution" {
					flags = append(flags, fmt.Sprintf("Disambiguation needed: %s", result.RuleName))
				}
			}
		}
	}

	return flags
}

// calculateQualityScore calculates the quality score of the rewritten prompt
func (pr *PromptRewriter) calculateQualityScore(result *RewriteResult) float64 {
	score := 1.0

	// Reduce score for many ambiguities
	score -= float64(len(result.Ambiguities)) * 0.1

	// Reduce score for many assumptions
	score -= float64(len(result.Assumptions)) * 0.05

	// Increase score for clarity flags (shows awareness)
	score += float64(len(result.ClarityFlags)) * 0.1

	// Increase score for disambiguation flags (shows awareness)
	score += float64(len(result.DisambiguationFlags)) * 0.1

	// Ensure score is between 0 and 1
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	return score
}

// generateRecommendations generates recommendations for the prompt
func (pr *PromptRewriter) generateRecommendations(result *RewriteResult) []string {
	var recommendations []string

	if len(result.Ambiguities) > 3 {
		recommendations = append(recommendations, "Consider resolving ambiguities before processing")
	}

	if len(result.Assumptions) > 2 {
		recommendations = append(recommendations, "Verify assumptions with user")
	}

	if result.QualityScore < 0.6 {
		recommendations = append(recommendations, "Prompt quality is below threshold - consider additional context")
	}

	if len(result.ClarityFlags) == 0 && len(result.DisambiguationFlags) == 0 {
		recommendations = append(recommendations, "Prompt appears clear and unambiguous")
	}

	return recommendations
}

// SimpleRewrite provides a simple prompt rewrite with minimal context
func (pr *PromptRewriter) SimpleRewrite(ctx context.Context, sessionID, userInput string, turnNumber int) (string, error) {
	request := &RewriteRequest{
		SessionID:  sessionID,
		TurnNumber: turnNumber,
		UserInput:  userInput,
		Options: RewriteOptions{
			IncludeContext:      true,
			IncludeAmbiguities:  true,
			AddClarityFlags:     true,
			MaxContextLength:    500,
			MaxHistoryTurns:     2,
			OptimizeForClarity:  true,
		},
	}

	result, err := pr.RewritePrompt(ctx, request)
	if err != nil {
		return "", err
	}

	return result.RewrittenPrompt, nil
}