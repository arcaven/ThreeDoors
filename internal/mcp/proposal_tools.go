package mcp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// proposalToolDefinitions returns tool definitions for proposal-related MCP tools.
func proposalToolDefinitions() []ToolItem {
	return []ToolItem{
		{
			Name:        "propose_enrichment",
			Description: "Propose an enrichment to an existing task. The proposal enters a review queue — it does not modify the task directly.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"task_id":   map[string]any{"type": "string", "description": "The task ID to enrich"},
					"type":      map[string]any{"type": "string", "description": "Proposal type: enrich-metadata, add-subtasks, add-context, add-note, suggest-blocker, suggest-category, update-effort"},
					"payload":   map[string]any{"type": "object", "description": "Type-specific enrichment data"},
					"rationale": map[string]any{"type": "string", "description": "Why this enrichment is suggested"},
				},
				"required": []string{"task_id", "type", "payload", "rationale"},
			},
		},
		{
			Name:        "suggest_task",
			Description: "Suggest a new task to be created. The suggestion enters a review queue for user approval.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"text":      map[string]any{"type": "string", "description": "Task description text"},
					"context":   map[string]any{"type": "string", "description": "Additional context for the task"},
					"effort":    map[string]any{"type": "string", "description": "Effort level: quick-win, medium, deep-work"},
					"type":      map[string]any{"type": "string", "description": "Task type: creative, administrative, technical, physical"},
					"rationale": map[string]any{"type": "string", "description": "Why this task is suggested"},
				},
				"required": []string{"text", "rationale"},
			},
		},
		{
			Name:        "suggest_relationship",
			Description: "Suggest a relationship between two tasks (e.g., blocks, depends-on, related-to).",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"from_id":       map[string]any{"type": "string", "description": "Source task ID"},
					"to_id":         map[string]any{"type": "string", "description": "Target task ID"},
					"relation_type": map[string]any{"type": "string", "description": "Relationship type: blocks, depends-on, related-to"},
					"rationale":     map[string]any{"type": "string", "description": "Why this relationship is suggested"},
				},
				"required": []string{"from_id", "to_id", "relation_type", "rationale"},
			},
		},
	}
}

// toolProposeEnrichment handles the propose_enrichment tool call.
func (s *MCPServer) toolProposeEnrichment(req *Request, args json.RawMessage) *Response {
	var params struct {
		TaskID    string          `json:"task_id"`
		Type      ProposalType    `json:"type"`
		Payload   json.RawMessage `json:"payload"`
		Rationale string          `json:"rationale"`
	}
	if args != nil {
		if err := json.Unmarshal(args, &params); err != nil {
			return NewErrorResponse(req.ID, CodeInvalidParams, fmt.Sprintf("invalid arguments: %v", err))
		}
	}

	if params.TaskID == "" {
		return NewErrorResponse(req.ID, CodeInvalidParams, "task_id is required")
	}
	if params.Rationale == "" {
		return NewErrorResponse(req.ID, CodeInvalidParams, "rationale is required")
	}

	// Verify task exists.
	task := s.pool.GetTask(params.TaskID)
	if task == nil {
		return s.toolError(req, fmt.Sprintf("task not found: %s", params.TaskID))
	}

	proposal, err := NewProposal(params.Type, params.TaskID, task.UpdatedAt, params.Payload, s.clientSource(), params.Rationale)
	if err != nil {
		return s.toolError(req, fmt.Sprintf("invalid proposal: %v", err))
	}

	if err := s.proposalStore.Create(proposal); err != nil {
		return s.toolError(req, fmt.Sprintf("create proposal: %v", err))
	}

	return s.toolJSON(req, map[string]any{
		"proposal_id": proposal.ID,
		"status":      proposal.Status,
		"expires_at":  proposal.ExpiresAt,
		"message":     "proposal submitted for review",
	})
}

// toolSuggestTask handles the suggest_task tool call.
func (s *MCPServer) toolSuggestTask(req *Request, args json.RawMessage) *Response {
	var params struct {
		Text      string `json:"text"`
		Context   string `json:"context"`
		Effort    string `json:"effort"`
		Type      string `json:"type"`
		Rationale string `json:"rationale"`
	}
	if args != nil {
		if err := json.Unmarshal(args, &params); err != nil {
			return NewErrorResponse(req.ID, CodeInvalidParams, fmt.Sprintf("invalid arguments: %v", err))
		}
	}

	if params.Text == "" {
		return NewErrorResponse(req.ID, CodeInvalidParams, "text is required")
	}
	if params.Rationale == "" {
		return NewErrorResponse(req.ID, CodeInvalidParams, "rationale is required")
	}

	payload, err := json.Marshal(map[string]string{
		"text":    params.Text,
		"context": params.Context,
		"effort":  params.Effort,
		"type":    params.Type,
	})
	if err != nil {
		return s.toolError(req, fmt.Sprintf("marshal payload: %v", err))
	}

	// suggest_task uses a synthetic task ID since the task doesn't exist yet.
	syntheticID := "new:" + uuid.New().String()
	now := time.Now().UTC()

	proposal, err := NewProposal(ProposalAddContext, syntheticID, now, payload, s.clientSource(), params.Rationale)
	if err != nil {
		return s.toolError(req, fmt.Sprintf("invalid proposal: %v", err))
	}

	if err := s.proposalStore.Create(proposal); err != nil {
		return s.toolError(req, fmt.Sprintf("create proposal: %v", err))
	}

	return s.toolJSON(req, map[string]any{
		"proposal_id": proposal.ID,
		"status":      proposal.Status,
		"expires_at":  proposal.ExpiresAt,
		"message":     "task suggestion submitted for review",
	})
}

// toolSuggestRelationship handles the suggest_relationship tool call.
func (s *MCPServer) toolSuggestRelationship(req *Request, args json.RawMessage) *Response {
	var params struct {
		FromID       string `json:"from_id"`
		ToID         string `json:"to_id"`
		RelationType string `json:"relation_type"`
		Rationale    string `json:"rationale"`
	}
	if args != nil {
		if err := json.Unmarshal(args, &params); err != nil {
			return NewErrorResponse(req.ID, CodeInvalidParams, fmt.Sprintf("invalid arguments: %v", err))
		}
	}

	if params.FromID == "" || params.ToID == "" {
		return NewErrorResponse(req.ID, CodeInvalidParams, "from_id and to_id are required")
	}
	if params.RelationType == "" {
		return NewErrorResponse(req.ID, CodeInvalidParams, "relation_type is required")
	}
	if params.Rationale == "" {
		return NewErrorResponse(req.ID, CodeInvalidParams, "rationale is required")
	}

	// Verify both tasks exist.
	fromTask := s.pool.GetTask(params.FromID)
	if fromTask == nil {
		return s.toolError(req, fmt.Sprintf("from task not found: %s", params.FromID))
	}
	toTask := s.pool.GetTask(params.ToID)
	if toTask == nil {
		return s.toolError(req, fmt.Sprintf("to task not found: %s", params.ToID))
	}

	payload, err := json.Marshal(map[string]string{
		"from_id":       params.FromID,
		"to_id":         params.ToID,
		"relation_type": params.RelationType,
	})
	if err != nil {
		return s.toolError(req, fmt.Sprintf("marshal payload: %v", err))
	}

	proposal, err := NewProposal(ProposalSuggestRelation, params.FromID, fromTask.UpdatedAt, payload, s.clientSource(), params.Rationale)
	if err != nil {
		return s.toolError(req, fmt.Sprintf("invalid proposal: %v", err))
	}

	if err := s.proposalStore.Create(proposal); err != nil {
		return s.toolError(req, fmt.Sprintf("create proposal: %v", err))
	}

	return s.toolJSON(req, map[string]any{
		"proposal_id": proposal.ID,
		"status":      proposal.Status,
		"expires_at":  proposal.ExpiresAt,
		"message":     "relationship suggestion submitted for review",
	})
}

// clientSource returns the MCP client source identifier from initialization.
func (s *MCPServer) clientSource() string {
	if s.clientInfo.Name != "" {
		return fmt.Sprintf("mcp:%s", s.clientInfo.Name)
	}
	return "mcp:unknown"
}
