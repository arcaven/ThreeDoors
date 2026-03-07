package mcp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ProposalStatus represents the lifecycle state of a proposal.
type ProposalStatus string

const (
	ProposalPending  ProposalStatus = "pending"
	ProposalApproved ProposalStatus = "approved"
	ProposalRejected ProposalStatus = "rejected"
	ProposalExpired  ProposalStatus = "expired"
	ProposalStale    ProposalStatus = "stale"
)

// ProposalType categorizes the kind of enrichment being proposed.
type ProposalType string

const (
	ProposalEnrichMetadata  ProposalType = "enrich-metadata"
	ProposalAddSubtasks     ProposalType = "add-subtasks"
	ProposalAddContext      ProposalType = "add-context"
	ProposalAddNote         ProposalType = "add-note"
	ProposalSuggestBlocker  ProposalType = "suggest-blocker"
	ProposalSuggestRelation ProposalType = "suggest-relationship"
	ProposalSuggestCategory ProposalType = "suggest-category"
	ProposalUpdateEffort    ProposalType = "update-effort"
)

// ValidProposalTypes is the set of allowed proposal types.
var ValidProposalTypes = map[ProposalType]bool{
	ProposalEnrichMetadata:  true,
	ProposalAddSubtasks:     true,
	ProposalAddContext:      true,
	ProposalAddNote:         true,
	ProposalSuggestBlocker:  true,
	ProposalSuggestRelation: true,
	ProposalSuggestCategory: true,
	ProposalUpdateEffort:    true,
}

// DefaultExpirationDuration is the default TTL for proposals.
const DefaultExpirationDuration = 7 * 24 * time.Hour

// MaxPendingPerTask is the cap on pending proposals per task.
const MaxPendingPerTask = 5

// DuplicateSimilarityThreshold is the text similarity threshold for dedup.
const DuplicateSimilarityThreshold = 0.85

// Proposal represents a suggested enrichment from an LLM or intake channel.
type Proposal struct {
	ID          string          `json:"id"`
	Type        ProposalType    `json:"type"`
	TaskID      string          `json:"task_id"`
	BaseVersion time.Time       `json:"base_version"`
	Payload     json.RawMessage `json:"payload"`
	Status      ProposalStatus  `json:"status"`
	Source      string          `json:"source"`
	Rationale   string          `json:"rationale"`
	CreatedAt   time.Time       `json:"created_at"`
	ReviewedAt  *time.Time      `json:"reviewed_at,omitempty"`
	ExpiresAt   time.Time       `json:"expires_at"`
}

// ProposalFilter controls which proposals are returned by List.
type ProposalFilter struct {
	TaskID string         `json:"task_id,omitempty"`
	Status ProposalStatus `json:"status,omitempty"`
	Type   ProposalType   `json:"type,omitempty"`
	Source string         `json:"source,omitempty"`
}

// NewProposal creates a proposal with defaults filled in.
func NewProposal(pType ProposalType, taskID string, baseVersion time.Time, payload json.RawMessage, source, rationale string) (*Proposal, error) {
	if !ValidProposalTypes[pType] {
		return nil, fmt.Errorf("invalid proposal type: %s", pType)
	}
	if taskID == "" {
		return nil, fmt.Errorf("task_id is required")
	}
	if len(payload) == 0 {
		return nil, fmt.Errorf("payload is required")
	}

	now := time.Now().UTC()
	return &Proposal{
		ID:          uuid.New().String(),
		Type:        pType,
		TaskID:      taskID,
		BaseVersion: baseVersion,
		Payload:     payload,
		Status:      ProposalPending,
		Source:      source,
		Rationale:   rationale,
		CreatedAt:   now,
		ExpiresAt:   now.Add(DefaultExpirationDuration),
	}, nil
}

// IsTerminal returns true if the proposal is in a final state.
func (p *Proposal) IsTerminal() bool {
	switch p.Status {
	case ProposalApproved, ProposalRejected, ProposalExpired:
		return true
	}
	return false
}
