package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Example request structures
type TrackContextRequest struct {
	SessionID  string `json:"session_id"`
	TurnNumber int    `json:"turn_number"`
	UserInput  string `json:"user_input"`
}

type RewriteRequest struct {
	SessionID  string `json:"session_id"`
	TurnNumber int    `json:"turn_number"`
	UserInput  string `json:"user_input"`
}

type PipelineRequest struct {
	SessionID    string `json:"session_id"`
	TurnNumber   int    `json:"turn_number"`
	UserInput    string `json:"user_input"`
	SystemPrompt string `json:"system_prompt"`
}

const (
	baseURL = "http://localhost:8080/api/v1"
)

func main() {
	fmt.Println("AI Context Gap Tracker - Example Usage")
	fmt.Println("=====================================")

	// Wait for services to start
	time.Sleep(2 * time.Second)

	// Example 1: Track context for a conversation
	fmt.Println("\n1. Context Tracking Example")
	sessionID := "demo-session-001"
	
	// First turn
	trackContext(sessionID, 1, "I want to get a visa for a place I mentioned earlier.")
	
	// Second turn  
	trackContext(sessionID, 2, "Actually, I changed my mind about Canada. I want to go to France instead.")

	// Example 2: Demonstrate rule evaluation
	fmt.Println("\n2. Rule Evaluation Example")
	evaluateRules(sessionID, 2, "I want to go there soon, but I'm not sure about the requirements.")

	// Example 3: Show prompt rewriting
	fmt.Println("\n3. Prompt Rewriting Example")
	rewritePrompt(sessionID, 3, "Can you help me with that thing we discussed?")

	// Example 4: Complete pipeline processing
	fmt.Println("\n4. Complete Pipeline Example")
	processPipeline(sessionID, 4, "I need information about the process for that European country.", "You are a helpful assistant for travel planning.")

	// Example 5: Response auditing
	fmt.Println("\n5. Response Auditing Example")
	auditResponse(sessionID, 4, "Based on your previous mention of France, I assume you're asking about French visa requirements. The process typically involves submitting documents to the French consulate.")

	// Example 6: Get session context
	fmt.Println("\n6. Session Context Retrieval")
	getSessionContext(sessionID)

	fmt.Println("\n✅ Example usage completed successfully!")
}

func trackContext(sessionID string, turnNumber int, userInput string) {
	fmt.Printf("📝 Tracking context for turn %d: %s\n", turnNumber, userInput)
	
	request := TrackContextRequest{
		SessionID:  sessionID,
		TurnNumber: turnNumber,
		UserInput:  userInput,
	}

	response, err := makeRequest("POST", "/context/track", request)
	if err != nil {
		log.Printf("❌ Error tracking context: %v", err)
		return
	}

	fmt.Printf("✅ Context tracked successfully\n")
	printResponse(response)
}

func evaluateRules(sessionID string, turnNumber int, userInput string) {
	fmt.Printf("🔍 Evaluating rules for: %s\n", userInput)
	
	request := map[string]interface{}{
		"session_id":  sessionID,
		"turn_number": turnNumber,
		"user_input":  userInput,
		"entities":    make(map[string]interface{}),
		"topics":      []string{},
		"timeline":    []interface{}{},
		"assertions":  []interface{}{},
		"ambiguities": []interface{}{},
		"history":     []interface{}{},
	}

	response, err := makeRequest("POST", "/rules/evaluate", request)
	if err != nil {
		log.Printf("❌ Error evaluating rules: %v", err)
		return
	}

	fmt.Printf("✅ Rules evaluated successfully\n")
	printResponse(response)
}

func rewritePrompt(sessionID string, turnNumber int, userInput string) {
	fmt.Printf("✏️ Rewriting prompt: %s\n", userInput)
	
	request := RewriteRequest{
		SessionID:  sessionID,
		TurnNumber: turnNumber,
		UserInput:  userInput,
	}

	response, err := makeRequest("POST", "/prompt/simple-rewrite", request)
	if err != nil {
		log.Printf("❌ Error rewriting prompt: %v", err)
		return
	}

	fmt.Printf("✅ Prompt rewritten successfully\n")
	printResponse(response)
}

func processPipeline(sessionID string, turnNumber int, userInput, systemPrompt string) {
	fmt.Printf("🔄 Processing pipeline for: %s\n", userInput)
	
	request := PipelineRequest{
		SessionID:    sessionID,
		TurnNumber:   turnNumber,
		UserInput:    userInput,
		SystemPrompt: systemPrompt,
	}

	response, err := makeRequest("POST", "/pipeline/process", request)
	if err != nil {
		log.Printf("❌ Error processing pipeline: %v", err)
		return
	}

	fmt.Printf("✅ Pipeline processed successfully\n")
	printResponse(response)
}

func auditResponse(sessionID string, turnNumber int, responseText string) {
	fmt.Printf("🔍 Auditing response: %s\n", responseText[:50]+"...")
	
	request := map[string]interface{}{
		"session_id":    sessionID,
		"turn_number":   turnNumber,
		"response_text": responseText,
		"context":       make(map[string]interface{}),
	}

	response, err := makeRequest("POST", "/audit/response", request)
	if err != nil {
		log.Printf("❌ Error auditing response: %v", err)
		return
	}

	fmt.Printf("✅ Response audited successfully\n")
	printResponse(response)
}

func getSessionContext(sessionID string) {
	fmt.Printf("📋 Getting session context for: %s\n", sessionID)
	
	resp, err := http.Get(fmt.Sprintf("%s/context/session/%s", baseURL, sessionID))
	if err != nil {
		log.Printf("❌ Error getting session context: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Error reading response: %v", err)
		return
	}

	fmt.Printf("✅ Session context retrieved successfully\n")
	printJSONResponse(body)
}

func makeRequest(method, endpoint string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	url := baseURL + endpoint
	
	var req *http.Request
	if method == "POST" {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func printResponse(response []byte) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, response, "", "  "); err != nil {
		fmt.Printf("Raw response: %s\n", string(response))
	} else {
		fmt.Printf("%s\n", prettyJSON.String())
	}
}

func printJSONResponse(response []byte) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, response, "", "  "); err != nil {
		fmt.Printf("Raw response: %s\n", string(response))
	} else {
		fmt.Printf("%s\n", prettyJSON.String())
	}
}