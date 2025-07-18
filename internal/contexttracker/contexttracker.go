package contexttracker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cliffordotieno/ai-context-gap-tracker/internal/database"
	"github.com/cliffordotieno/ai-context-gap-tracker/pkg/redis"
)

// ContextTracker manages conversational context
type ContextTracker struct {
	db    *database.DB
	redis *redis.Client
}

// Context represents a conversation context
type Context struct {
	ID          int                    `json:"id"`
	SessionID   string                 `json:"session_id"`
	TurnNumber  int                    `json:"turn_number"`
	UserInput   string                 `json:"user_input"`
	Entities    map[string]interface{} `json:"entities"`
	Topics      []string               `json:"topics"`
	Timeline    []TimelineEvent        `json:"timeline"`
	Assertions  []Assertion            `json:"assertions"`
	Ambiguities []Ambiguity            `json:"ambiguities"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// TimelineEvent represents a temporal event
type TimelineEvent struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Reference string    `json:"reference"`
}

// Assertion represents a factual claim
type Assertion struct {
	Claim      string  `json:"claim"`
	Confidence float64 `json:"confidence"`
	Source     string  `json:"source"`
}

// Ambiguity represents unclear information
type Ambiguity struct {
	Text        string   `json:"text"`
	Type        string   `json:"type"`
	Suggestions []string `json:"suggestions"`
}

// MemoryGraph represents the conversation memory structure
type MemoryGraph struct {
	SessionID string                 `json:"session_id"`
	Nodes     map[string]interface{} `json:"nodes"`
	Edges     []Edge                 `json:"edges"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// Edge represents a relationship between concepts
type Edge struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Type   string  `json:"type"`
	Weight float64 `json:"weight"`
}

// New creates a new ContextTracker instance
func New(db *database.DB, redisClient *redis.Client) *ContextTracker {
	return &ContextTracker{
		db:    db,
		redis: redisClient,
	}
}

// TrackContext stores and analyzes conversation context
func (ct *ContextTracker) TrackContext(ctx context.Context, sessionID string, turnNumber int, userInput string) (*Context, error) {
	// Create context object
	context := &Context{
		SessionID:   sessionID,
		TurnNumber:  turnNumber,
		UserInput:   userInput,
		Entities:    make(map[string]interface{}),
		Topics:      []string{},
		Timeline:    []TimelineEvent{},
		Assertions:  []Assertion{},
		Ambiguities: []Ambiguity{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Analyze entities (placeholder - would integrate with NLP service)
	context.Entities = ct.extractEntities(userInput)

	// Extract topics
	context.Topics = ct.extractTopics(userInput)

	// Identify timeline events
	context.Timeline = ct.extractTimelineEvents(userInput)

	// Extract assertions
	context.Assertions = ct.extractAssertions(userInput)

	// Identify ambiguities
	context.Ambiguities = ct.identifyAmbiguities(userInput)

	// Store in database
	if err := ct.storeContext(ctx, context); err != nil {
		return nil, fmt.Errorf("failed to store context: %w", err)
	}

	// Cache in Redis
	if err := ct.cacheContext(ctx, context); err != nil {
		log.Printf("Warning: failed to cache context: %v", err)
	}

	// Update memory graph
	if err := ct.updateMemoryGraph(ctx, sessionID, context); err != nil {
		log.Printf("Warning: failed to update memory graph: %v", err)
	}

	return context, nil
}

// GetContext retrieves context for a specific turn
func (ct *ContextTracker) GetContext(ctx context.Context, sessionID string, turnNumber int) (*Context, error) {
	// Try Redis first
	if cachedData, err := ct.redis.GetContext(ctx, sessionID, turnNumber); err == nil {
		var context Context
		if err := json.Unmarshal([]byte(cachedData), &context); err == nil {
			return &context, nil
		}
	}

	// Fallback to database
	query := `
		SELECT id, session_id, turn_number, user_input, entities, topics, timeline, assertions, ambiguities, created_at, updated_at
		FROM contexts
		WHERE session_id = $1 AND turn_number = $2
	`

	row := ct.db.QueryRow(query, sessionID, turnNumber)

	var context Context
	var entitiesJSON, topicsJSON, timelineJSON, assertionsJSON, ambiguitiesJSON []byte

	err := row.Scan(
		&context.ID,
		&context.SessionID,
		&context.TurnNumber,
		&context.UserInput,
		&entitiesJSON,
		&topicsJSON,
		&timelineJSON,
		&assertionsJSON,
		&ambiguitiesJSON,
		&context.CreatedAt,
		&context.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get context: %w", err)
	}

	// Parse JSON fields
	json.Unmarshal(entitiesJSON, &context.Entities)
	json.Unmarshal(topicsJSON, &context.Topics)
	json.Unmarshal(timelineJSON, &context.Timeline)
	json.Unmarshal(assertionsJSON, &context.Assertions)
	json.Unmarshal(ambiguitiesJSON, &context.Ambiguities)

	return &context, nil
}

// GetSessionContext retrieves all context for a session
func (ct *ContextTracker) GetSessionContext(ctx context.Context, sessionID string) ([]*Context, error) {
	query := `
		SELECT id, session_id, turn_number, user_input, entities, topics, timeline, assertions, ambiguities, created_at, updated_at
		FROM contexts
		WHERE session_id = $1
		ORDER BY turn_number ASC
	`

	rows, err := ct.db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session context: %w", err)
	}
	defer rows.Close()

	var contexts []*Context
	for rows.Next() {
		var context Context
		var entitiesJSON, topicsJSON, timelineJSON, assertionsJSON, ambiguitiesJSON []byte

		err := rows.Scan(
			&context.ID,
			&context.SessionID,
			&context.TurnNumber,
			&context.UserInput,
			&entitiesJSON,
			&topicsJSON,
			&timelineJSON,
			&assertionsJSON,
			&ambiguitiesJSON,
			&context.CreatedAt,
			&context.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan context: %w", err)
		}

		// Parse JSON fields
		json.Unmarshal(entitiesJSON, &context.Entities)
		json.Unmarshal(topicsJSON, &context.Topics)
		json.Unmarshal(timelineJSON, &context.Timeline)
		json.Unmarshal(assertionsJSON, &context.Assertions)
		json.Unmarshal(ambiguitiesJSON, &context.Ambiguities)

		contexts = append(contexts, &context)
	}

	return contexts, nil
}

// GetMemoryGraph retrieves the memory graph for a session
func (ct *ContextTracker) GetMemoryGraph(ctx context.Context, sessionID string) (*MemoryGraph, error) {
	// Try Redis first
	if cachedData, err := ct.redis.GetMemoryGraph(ctx, sessionID); err == nil {
		var graph MemoryGraph
		if err := json.Unmarshal([]byte(cachedData), &graph); err == nil {
			return &graph, nil
		}
	}

	// Fallback to database
	query := `
		SELECT context_graph FROM sessions WHERE id = $1
	`

	row := ct.db.QueryRow(query, sessionID)

	var graphJSON []byte
	err := row.Scan(&graphJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory graph: %w", err)
	}

	var graph MemoryGraph
	if err := json.Unmarshal(graphJSON, &graph); err != nil {
		return nil, fmt.Errorf("failed to unmarshal memory graph: %w", err)
	}

	return &graph, nil
}

// Helper methods (placeholder implementations)
func (ct *ContextTracker) extractEntities(input string) map[string]interface{} {
	// Placeholder - would integrate with NLP service
	return make(map[string]interface{})
}

func (ct *ContextTracker) extractTopics(input string) []string {
	// Placeholder - would integrate with NLP service
	return []string{}
}

func (ct *ContextTracker) extractTimelineEvents(input string) []TimelineEvent {
	// Placeholder - would integrate with NLP service
	return []TimelineEvent{}
}

func (ct *ContextTracker) extractAssertions(input string) []Assertion {
	// Placeholder - would integrate with NLP service
	return []Assertion{}
}

func (ct *ContextTracker) identifyAmbiguities(input string) []Ambiguity {
	// Placeholder - would integrate with NLP service
	return []Ambiguity{}
}

func (ct *ContextTracker) storeContext(ctx context.Context, context *Context) error {
	entitiesJSON, _ := json.Marshal(context.Entities)
	topicsJSON, _ := json.Marshal(context.Topics)
	timelineJSON, _ := json.Marshal(context.Timeline)
	assertionsJSON, _ := json.Marshal(context.Assertions)
	ambiguitiesJSON, _ := json.Marshal(context.Ambiguities)

	query := `
		INSERT INTO contexts (session_id, turn_number, user_input, entities, topics, timeline, assertions, ambiguities)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (session_id, turn_number) DO UPDATE SET
		user_input = $3, entities = $4, topics = $5, timeline = $6, assertions = $7, ambiguities = $8, updated_at = CURRENT_TIMESTAMP
	`

	_, err := ct.db.Exec(query, context.SessionID, context.TurnNumber, context.UserInput,
		entitiesJSON, topicsJSON, timelineJSON, assertionsJSON, ambiguitiesJSON)

	return err
}

func (ct *ContextTracker) cacheContext(ctx context.Context, context *Context) error {
	data, err := json.Marshal(context)
	if err != nil {
		return err
	}

	return ct.redis.SetContext(ctx, context.SessionID, context.TurnNumber, string(data))
}

func (ct *ContextTracker) updateMemoryGraph(ctx context.Context, sessionID string, context *Context) error {
	// Placeholder - would implement graph update logic
	graph := &MemoryGraph{
		SessionID: sessionID,
		Nodes:     make(map[string]interface{}),
		Edges:     []Edge{},
		UpdatedAt: time.Now(),
	}

	data, err := json.Marshal(graph)
	if err != nil {
		return err
	}

	return ct.redis.SetMemoryGraph(ctx, sessionID, string(data))
}