package responseauditor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cliffordotieno/ai-context-gap-tracker/internal/database"
)

// ResponseAuditor audits and classifies responses
type ResponseAuditor struct {
	db *database.DB
}

// AuditResult represents the result of response auditing
type AuditResult struct {
	ID               int                    `json:"id"`
	SessionID        string                 `json:"session_id"`
	TurnNumber       int                    `json:"turn_number"`
	ResponseText     string                 `json:"response_text"`
	CertaintyLevel   string                 `json:"certainty_level"`
	Flags            map[string]interface{} `json:"flags"`
	Assumptions      []Assumption           `json:"assumptions"`
	Contradictions   []Contradiction        `json:"contradictions"`
	RetryCount       int                    `json:"retry_count"`
	Recommendations  []string               `json:"recommendations"`
	QualityScore     float64                `json:"quality_score"`
	CreatedAt        time.Time              `json:"created_at"`
}

// Assumption represents an assumption made in the response
type Assumption struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
	Source     string  `json:"source"`
	Critical   bool    `json:"critical"`
}

// Contradiction represents a contradiction found in the response
type Contradiction struct {
	Text        string  `json:"text"`
	ConflictsWith string `json:"conflicts_with"`
	Severity    string  `json:"severity"`
	Confidence  float64 `json:"confidence"`
}

// CertaintyLevel represents different levels of response certainty
type CertaintyLevel string

const (
	CertaintyHigh      CertaintyLevel = "high"
	CertaintyMedium    CertaintyLevel = "medium"
	CertaintyLow       CertaintyLevel = "low"
	CertaintyAssumed   CertaintyLevel = "assumed"
	CertaintyInferred  CertaintyLevel = "inferred"
	CertaintyAmbiguous CertaintyLevel = "ambiguous"
	CertaintyVerified  CertaintyLevel = "verified"
)

// New creates a new ResponseAuditor instance
func New(db *database.DB) *ResponseAuditor {
	return &ResponseAuditor{
		db: db,
	}
}

// AuditResponse audits a response and returns the audit result
func (ra *ResponseAuditor) AuditResponse(ctx context.Context, sessionID string, turnNumber int, responseText string, contextData map[string]interface{}) (*AuditResult, error) {
	auditResult := &AuditResult{
		SessionID:       sessionID,
		TurnNumber:      turnNumber,
		ResponseText:    responseText,
		Flags:           make(map[string]interface{}),
		Assumptions:     []Assumption{},
		Contradictions:  []Contradiction{},
		RetryCount:      0,
		Recommendations: []string{},
		QualityScore:    0.0,
		CreatedAt:       time.Now(),
	}

	// Classify certainty level
	auditResult.CertaintyLevel = ra.classifyCertaintyLevel(responseText)

	// Detect assumptions
	auditResult.Assumptions = ra.detectAssumptions(responseText, contextData)

	// Detect contradictions
	auditResult.Contradictions = ra.detectContradictions(responseText, contextData)

	// Set flags
	auditResult.Flags = ra.setFlags(responseText, auditResult)

	// Calculate quality score
	auditResult.QualityScore = ra.calculateQualityScore(auditResult)

	// Generate recommendations
	auditResult.Recommendations = ra.generateRecommendations(auditResult)

	// Store audit result
	if err := ra.storeAuditResult(ctx, auditResult); err != nil {
		return nil, fmt.Errorf("failed to store audit result: %w", err)
	}

	return auditResult, nil
}

// classifyCertaintyLevel determines the certainty level of a response
func (ra *ResponseAuditor) classifyCertaintyLevel(responseText string) string {
	text := strings.ToLower(responseText)

	// High certainty indicators
	highCertaintyKeywords := []string{
		"definitely", "certainly", "absolutely", "confirmed", "verified",
		"proven", "established", "documented", "factual", "precisely",
	}

	// Low certainty indicators
	lowCertaintyKeywords := []string{
		"maybe", "perhaps", "possibly", "might", "could be",
		"seems", "appears", "likely", "probably", "potentially",
	}

	// Assumption indicators
	assumptionKeywords := []string{
		"assuming", "suppose", "presuming", "let's say", "if we assume",
		"taking for granted", "based on the assumption", "presumably",
	}

	// Inference indicators
	inferenceKeywords := []string{
		"infer", "deduce", "conclude", "suggest", "imply",
		"based on", "from this we can", "it follows that",
	}

	// Ambiguity indicators
	ambiguityKeywords := []string{
		"unclear", "ambiguous", "uncertain", "vague", "confusing",
		"multiple interpretations", "could mean", "not sure",
	}

	// Check for different certainty levels
	for _, keyword := range highCertaintyKeywords {
		if strings.Contains(text, keyword) {
			return string(CertaintyHigh)
		}
	}

	for _, keyword := range assumptionKeywords {
		if strings.Contains(text, keyword) {
			return string(CertaintyAssumed)
		}
	}

	for _, keyword := range inferenceKeywords {
		if strings.Contains(text, keyword) {
			return string(CertaintyInferred)
		}
	}

	for _, keyword := range ambiguityKeywords {
		if strings.Contains(text, keyword) {
			return string(CertaintyAmbiguous)
		}
	}

	for _, keyword := range lowCertaintyKeywords {
		if strings.Contains(text, keyword) {
			return string(CertaintyLow)
		}
	}

	// Default to medium certainty
	return string(CertaintyMedium)
}

// detectAssumptions identifies assumptions in the response
func (ra *ResponseAuditor) detectAssumptions(responseText string, contextData map[string]interface{}) []Assumption {
	var assumptions []Assumption
	text := strings.ToLower(responseText)

	// Assumption patterns
	assumptionPatterns := []struct {
		keywords   []string
		confidence float64
		critical   bool
	}{
		{
			keywords:   []string{"assuming", "suppose", "presuming", "let's say"},
			confidence: 0.9,
			critical:   true,
		},
		{
			keywords:   []string{"if we assume", "based on the assumption", "presumably"},
			confidence: 0.8,
			critical:   true,
		},
		{
			keywords:   []string{"likely", "probably", "seems", "appears"},
			confidence: 0.6,
			critical:   false,
		},
		{
			keywords:   []string{"might", "could be", "possibly", "perhaps"},
			confidence: 0.4,
			critical:   false,
		},
	}

	for _, pattern := range assumptionPatterns {
		for _, keyword := range pattern.keywords {
			if strings.Contains(text, keyword) {
				assumption := Assumption{
					Text:       ra.extractAssumptionText(responseText, keyword),
					Confidence: pattern.confidence,
					Source:     "keyword_detection",
					Critical:   pattern.critical,
				}
				assumptions = append(assumptions, assumption)
			}
		}
	}

	return assumptions
}

// detectContradictions identifies contradictions in the response
func (ra *ResponseAuditor) detectContradictions(responseText string, contextData map[string]interface{}) []Contradiction {
	var contradictions []Contradiction
	text := strings.ToLower(responseText)

	// Contradiction patterns
	contradictoryPairs := []struct {
		words    []string
		severity string
	}{
		{words: []string{"yes", "no"}, severity: "high"},
		{words: []string{"always", "never"}, severity: "high"},
		{words: []string{"all", "none"}, severity: "high"},
		{words: []string{"increase", "decrease"}, severity: "medium"},
		{words: []string{"before", "after"}, severity: "medium"},
		{words: []string{"more", "less"}, severity: "low"},
	}

	for _, pair := range contradictoryPairs {
		if strings.Contains(text, pair.words[0]) && strings.Contains(text, pair.words[1]) {
			contradiction := Contradiction{
				Text:          fmt.Sprintf("Contains both '%s' and '%s'", pair.words[0], pair.words[1]),
				ConflictsWith: fmt.Sprintf("'%s' conflicts with '%s'", pair.words[0], pair.words[1]),
				Severity:      pair.severity,
				Confidence:    0.7,
			}
			contradictions = append(contradictions, contradiction)
		}
	}

	return contradictions
}

// setFlags sets various flags based on the response analysis
func (ra *ResponseAuditor) setFlags(responseText string, auditResult *AuditResult) map[string]interface{} {
	flags := make(map[string]interface{})
	text := strings.ToLower(responseText)

	// Length-based flags
	flags["response_length"] = len(responseText)
	flags["short_response"] = len(responseText) < 50
	flags["long_response"] = len(responseText) > 500

	// Content-based flags
	flags["contains_assumptions"] = len(auditResult.Assumptions) > 0
	flags["contains_contradictions"] = len(auditResult.Contradictions) > 0
	flags["high_certainty"] = auditResult.CertaintyLevel == string(CertaintyHigh)
	flags["low_certainty"] = auditResult.CertaintyLevel == string(CertaintyLow)

	// Question flags
	flags["contains_questions"] = strings.Contains(text, "?")
	flags["clarification_request"] = strings.Contains(text, "please clarify") || strings.Contains(text, "can you specify")

	// Hedge words
	hedgeWords := []string{"might", "could", "possibly", "perhaps", "maybe", "likely", "probably"}
	hedgeCount := 0
	for _, word := range hedgeWords {
		if strings.Contains(text, word) {
			hedgeCount++
		}
	}
	flags["hedge_words_count"] = hedgeCount
	flags["excessive_hedging"] = hedgeCount > 3

	// Confidence indicators
	flags["confidence_stated"] = strings.Contains(text, "confident") || strings.Contains(text, "certain")
	flags["uncertainty_stated"] = strings.Contains(text, "uncertain") || strings.Contains(text, "not sure")

	return flags
}

// calculateQualityScore calculates an overall quality score for the response
func (ra *ResponseAuditor) calculateQualityScore(auditResult *AuditResult) float64 {
	score := 1.0

	// Penalize for assumptions
	score -= float64(len(auditResult.Assumptions)) * 0.1

	// Penalize for contradictions
	score -= float64(len(auditResult.Contradictions)) * 0.2

	// Adjust based on certainty level
	switch auditResult.CertaintyLevel {
	case string(CertaintyHigh):
		score += 0.2
	case string(CertaintyVerified):
		score += 0.3
	case string(CertaintyLow):
		score -= 0.1
	case string(CertaintyAmbiguous):
		score -= 0.2
	}

	// Penalize for excessive hedging
	if hedgeCount, ok := auditResult.Flags["hedge_words_count"].(int); ok && hedgeCount > 3 {
		score -= 0.1
	}

	// Ensure score is between 0 and 1
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	return score
}

// generateRecommendations generates recommendations based on the audit result
func (ra *ResponseAuditor) generateRecommendations(auditResult *AuditResult) []string {
	var recommendations []string

	// Recommendations based on certainty level
	switch auditResult.CertaintyLevel {
	case string(CertaintyLow):
		recommendations = append(recommendations, "Consider providing more definitive information or seeking clarification")
	case string(CertaintyAssumed):
		recommendations = append(recommendations, "Verify assumptions before providing final response")
	case string(CertaintyAmbiguous):
		recommendations = append(recommendations, "Clarify ambiguous statements and provide clearer explanations")
	}

	// Recommendations based on assumptions
	if len(auditResult.Assumptions) > 0 {
		recommendations = append(recommendations, "Consider explicitly stating assumptions or seeking confirmation")
	}

	// Recommendations based on contradictions
	if len(auditResult.Contradictions) > 0 {
		recommendations = append(recommendations, "Resolve contradictions in the response")
	}

	// Recommendations based on quality score
	if auditResult.QualityScore < 0.6 {
		recommendations = append(recommendations, "Response quality is below threshold - consider revision")
	}

	// Recommendations based on flags
	if excessive, ok := auditResult.Flags["excessive_hedging"].(bool); ok && excessive {
		recommendations = append(recommendations, "Reduce excessive use of hedge words for clearer communication")
	}

	if short, ok := auditResult.Flags["short_response"].(bool); ok && short {
		recommendations = append(recommendations, "Consider providing more detailed response")
	}

	return recommendations
}

// storeAuditResult stores the audit result in the database
func (ra *ResponseAuditor) storeAuditResult(ctx context.Context, auditResult *AuditResult) error {
	flagsJSON, _ := json.Marshal(auditResult.Flags)
	assumptionsJSON, _ := json.Marshal(auditResult.Assumptions)
	contradictionsJSON, _ := json.Marshal(auditResult.Contradictions)

	query := `
		INSERT INTO audit_logs (session_id, turn_number, response_text, certainty_level, flags, assumptions, contradictions, retry_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	err := ra.db.QueryRow(query, auditResult.SessionID, auditResult.TurnNumber, auditResult.ResponseText,
		auditResult.CertaintyLevel, flagsJSON, assumptionsJSON, contradictionsJSON, auditResult.RetryCount).Scan(
		&auditResult.ID, &auditResult.CreatedAt)

	return err
}

// GetAuditHistory retrieves audit history for a session
func (ra *ResponseAuditor) GetAuditHistory(ctx context.Context, sessionID string) ([]*AuditResult, error) {
	query := `
		SELECT id, session_id, turn_number, response_text, certainty_level, flags, assumptions, contradictions, retry_count, created_at
		FROM audit_logs
		WHERE session_id = $1
		ORDER BY turn_number ASC
	`

	rows, err := ra.db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit history: %w", err)
	}
	defer rows.Close()

	var results []*AuditResult
	for rows.Next() {
		var result AuditResult
		var flagsJSON, assumptionsJSON, contradictionsJSON []byte

		err := rows.Scan(
			&result.ID,
			&result.SessionID,
			&result.TurnNumber,
			&result.ResponseText,
			&result.CertaintyLevel,
			&flagsJSON,
			&assumptionsJSON,
			&contradictionsJSON,
			&result.RetryCount,
			&result.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan audit result: %w", err)
		}

		// Parse JSON fields
		json.Unmarshal(flagsJSON, &result.Flags)
		json.Unmarshal(assumptionsJSON, &result.Assumptions)
		json.Unmarshal(contradictionsJSON, &result.Contradictions)

		results = append(results, &result)
	}

	return results, nil
}

// extractAssumptionText extracts the text around an assumption keyword
func (ra *ResponseAuditor) extractAssumptionText(responseText, keyword string) string {
	// Simple implementation - in a real scenario, this would use more sophisticated NLP
	index := strings.Index(strings.ToLower(responseText), keyword)
	if index == -1 {
		return keyword
	}

	start := index
	end := index + len(keyword) + 50
	if end > len(responseText) {
		end = len(responseText)
	}

	return responseText[start:end]
}

// ShouldRetry determines if a response should be retried based on audit results
func (ra *ResponseAuditor) ShouldRetry(auditResult *AuditResult) bool {
	// Retry if quality score is very low
	if auditResult.QualityScore < 0.3 {
		return true
	}

	// Retry if there are critical contradictions
	for _, contradiction := range auditResult.Contradictions {
		if contradiction.Severity == "high" {
			return true
		}
	}

	// Retry if there are critical assumptions
	for _, assumption := range auditResult.Assumptions {
		if assumption.Critical && assumption.Confidence > 0.8 {
			return true
		}
	}

	return false
}