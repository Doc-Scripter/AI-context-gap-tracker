package server

import (
	"net/http"
	"strconv"

	"github.com/cliffordotieno/ai-context-gap-tracker/internal/contexttracker"
	"github.com/cliffordotieno/ai-context-gap-tracker/internal/logicengine"
	"github.com/cliffordotieno/ai-context-gap-tracker/internal/promptrewriter"
	"github.com/cliffordotieno/ai-context-gap-tracker/internal/responseauditor"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// HTTPServer represents the HTTP server
type HTTPServer struct {
	router          *gin.Engine
	contextTracker  *contexttracker.ContextTracker
	logicEngine     *logicengine.LogicEngine
	responseAuditor *responseauditor.ResponseAuditor
	promptRewriter  *promptrewriter.PromptRewriter
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(router *gin.Engine, contextTracker *contexttracker.ContextTracker, logicEngine *logicengine.LogicEngine, responseAuditor *responseauditor.ResponseAuditor, promptRewriter *promptrewriter.PromptRewriter) *HTTPServer {
	return &HTTPServer{
		router:          router,
		contextTracker:  contextTracker,
		logicEngine:     logicEngine,
		responseAuditor: responseAuditor,
		promptRewriter:  promptRewriter,
	}
}

// SetupRoutes sets up the HTTP routes
func (s *HTTPServer) SetupRoutes() {
	api := s.router.Group("/api/v1")

	// Health check
	api.GET("/health", s.healthCheck)

	// Context tracking routes
	contextGroup := api.Group("/context")
	{
		contextGroup.POST("/track", s.trackContext)
		contextGroup.GET("/session/:sessionId", s.getSessionContext)
		contextGroup.GET("/session/:sessionId/turn/:turnNumber", s.getContext)
		contextGroup.GET("/session/:sessionId/memory", s.getMemoryGraph)
	}

	// Logic engine routes
	rulesGroup := api.Group("/rules")
	{
		rulesGroup.GET("", s.getRules)
		rulesGroup.POST("", s.createRule)
		rulesGroup.PUT("/:id", s.updateRule)
		rulesGroup.DELETE("/:id", s.deleteRule)
		rulesGroup.POST("/evaluate", s.evaluateRules)
		rulesGroup.POST("/initialize", s.initializeDefaultRules)
	}

	// Response auditor routes
	auditGroup := api.Group("/audit")
	{
		auditGroup.POST("/response", s.auditResponse)
		auditGroup.GET("/session/:sessionId/history", s.getAuditHistory)
	}

	// Prompt rewriter routes
	promptGroup := api.Group("/prompt")
	{
		promptGroup.POST("/rewrite", s.rewritePrompt)
		promptGroup.POST("/simple-rewrite", s.simpleRewrite)
	}

	// Pipeline route - combines all modules
	api.POST("/pipeline/process", s.processPipeline)
}

// Health check endpoint
func (s *HTTPServer) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "ai-context-gap-tracker",
	})
}

// Context tracking endpoints
func (s *HTTPServer) trackContext(c *gin.Context) {
	var request struct {
		SessionID  string `json:"session_id" binding:"required"`
		TurnNumber int    `json:"turn_number" binding:"required"`
		UserInput  string `json:"user_input" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context, err := s.contextTracker.TrackContext(c.Request.Context(), request.SessionID, request.TurnNumber, request.UserInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, context)
}

func (s *HTTPServer) getSessionContext(c *gin.Context) {
	sessionID := c.Param("sessionId")

	contexts, err := s.contextTracker.GetSessionContext(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contexts)
}

func (s *HTTPServer) getContext(c *gin.Context) {
	sessionID := c.Param("sessionId")
	turnNumberStr := c.Param("turnNumber")

	turnNumber, err := strconv.Atoi(turnNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid turn number"})
		return
	}

	context, err := s.contextTracker.GetContext(c.Request.Context(), sessionID, turnNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, context)
}

func (s *HTTPServer) getMemoryGraph(c *gin.Context) {
	sessionID := c.Param("sessionId")

	graph, err := s.contextTracker.GetMemoryGraph(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, graph)
}

// Logic engine endpoints
func (s *HTTPServer) getRules(c *gin.Context) {
	rules, err := s.logicEngine.GetActiveRules(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rules)
}

func (s *HTTPServer) createRule(c *gin.Context) {
	var rule logicengine.Rule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.logicEngine.CreateRule(c.Request.Context(), &rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

func (s *HTTPServer) updateRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	var rule logicengine.Rule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule.ID = id
	if err := s.logicEngine.UpdateRule(c.Request.Context(), &rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rule)
}

func (s *HTTPServer) deleteRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	if err := s.logicEngine.DeleteRule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rule deleted successfully"})
}

func (s *HTTPServer) evaluateRules(c *gin.Context) {
	var request logicengine.EvaluationContext
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := s.logicEngine.EvaluateRules(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (s *HTTPServer) initializeDefaultRules(c *gin.Context) {
	if err := s.logicEngine.InitializeDefaultRules(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "default rules initialized successfully"})
}

// Response auditor endpoints
func (s *HTTPServer) auditResponse(c *gin.Context) {
	var request struct {
		SessionID    string                 `json:"session_id" binding:"required"`
		TurnNumber   int                    `json:"turn_number" binding:"required"`
		ResponseText string                 `json:"response_text" binding:"required"`
		Context      map[string]interface{} `json:"context"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Context == nil {
		request.Context = make(map[string]interface{})
	}

	result, err := s.responseAuditor.AuditResponse(c.Request.Context(), request.SessionID, request.TurnNumber, request.ResponseText, request.Context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *HTTPServer) getAuditHistory(c *gin.Context) {
	sessionID := c.Param("sessionId")

	history, err := s.responseAuditor.GetAuditHistory(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// Prompt rewriter endpoints
func (s *HTTPServer) rewritePrompt(c *gin.Context) {
	var request promptrewriter.RewriteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := s.promptRewriter.RewritePrompt(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *HTTPServer) simpleRewrite(c *gin.Context) {
	var request struct {
		SessionID  string `json:"session_id" binding:"required"`
		TurnNumber int    `json:"turn_number" binding:"required"`
		UserInput  string `json:"user_input" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rewrittenPrompt, err := s.promptRewriter.SimpleRewrite(c.Request.Context(), request.SessionID, request.UserInput, request.TurnNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rewritten_prompt": rewrittenPrompt})
}

// Pipeline processing endpoint
func (s *HTTPServer) processPipeline(c *gin.Context) {
	var request struct {
		SessionID    string                 `json:"session_id" binding:"required"`
		TurnNumber   int                    `json:"turn_number" binding:"required"`
		UserInput    string                 `json:"user_input" binding:"required"`
		SystemPrompt string                 `json:"system_prompt"`
		Options      map[string]interface{} `json:"options"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// Step 1: Track context
	contextResult, err := s.contextTracker.TrackContext(ctx, request.SessionID, request.TurnNumber, request.UserInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "context tracking failed: " + err.Error()})
		return
	}

	// Step 2: Evaluate rules
	evalContext := &logicengine.EvaluationContext{
		SessionID:   request.SessionID,
		TurnNumber:  request.TurnNumber,
		UserInput:   request.UserInput,
		Entities:    contextResult.Entities,
		Topics:      contextResult.Topics,
		Timeline:    make([]interface{}, len(contextResult.Timeline)),
		Assertions:  make([]interface{}, len(contextResult.Assertions)),
		Ambiguities: make([]interface{}, len(contextResult.Ambiguities)),
	}

	// Convert context data to interface slices
	for i, item := range contextResult.Timeline {
		evalContext.Timeline[i] = item
	}
	for i, item := range contextResult.Assertions {
		evalContext.Assertions[i] = item
	}
	for i, item := range contextResult.Ambiguities {
		evalContext.Ambiguities[i] = item
	}

	ruleResults, err := s.logicEngine.EvaluateRules(ctx, evalContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rule evaluation failed: " + err.Error()})
		return
	}

	// Step 3: Rewrite prompt
	rewriteRequest := &promptrewriter.RewriteRequest{
		SessionID:    request.SessionID,
		TurnNumber:   request.TurnNumber,
		UserInput:    request.UserInput,
		SystemPrompt: request.SystemPrompt,
		Options:      promptrewriter.DefaultRewriteOptions(),
	}

	promptResult, err := s.promptRewriter.RewritePrompt(ctx, rewriteRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "prompt rewriting failed: " + err.Error()})
		return
	}

	// Return combined pipeline results
	response := gin.H{
		"session_id":     request.SessionID,
		"turn_number":    request.TurnNumber,
		"context":        contextResult,
		"rule_results":   ruleResults,
		"prompt_result":  promptResult,
		"pipeline_stage": "completed",
	}

	c.JSON(http.StatusOK, response)
}

// RegisterGRPCServices registers gRPC services (placeholder)
func RegisterGRPCServices(server *grpc.Server, contextTracker *contexttracker.ContextTracker, logicEngine *logicengine.LogicEngine, responseAuditor *responseauditor.ResponseAuditor, promptRewriter *promptrewriter.PromptRewriter) {
	// TODO: Implement gRPC services
	// This would require creating protobuf definitions and implementing the service handlers
	// For now, we'll use REST API only
}