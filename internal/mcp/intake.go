package mcp

import (
	"context"
	"fmt"
	"time"
)

// HealthStatus represents the health state of an intake channel.
type HealthStatus string

const (
	HealthOK       HealthStatus = "ok"
	HealthDegraded HealthStatus = "degraded"
	HealthDown     HealthStatus = "down"
)

// HealthCheckResult holds the result of an intake channel health check.
type HealthCheckResult struct {
	Status  HealthStatus `json:"status"`
	Message string       `json:"message,omitempty"`
}

// IntakeChannel defines a source that can suggest proposals.
type IntakeChannel interface {
	Name() string
	Suggest(ctx context.Context, proposal *Proposal) error
	HealthCheck() HealthCheckResult
}

// LLMIntakeChannel processes proposals from MCP LLM clients.
type LLMIntakeChannel struct {
	store  *ProposalStore
	source string
}

// NewLLMIntakeChannel creates an intake channel for a specific MCP client.
func NewLLMIntakeChannel(store *ProposalStore, source string) *LLMIntakeChannel {
	return &LLMIntakeChannel{
		store:  store,
		source: source,
	}
}

// Name returns the channel identifier.
func (c *LLMIntakeChannel) Name() string {
	return c.source
}

// Suggest submits a proposal through this intake channel.
func (c *LLMIntakeChannel) Suggest(ctx context.Context, proposal *Proposal) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("suggest proposal: %w", ctx.Err())
	default:
	}

	proposal.Source = c.source
	proposal.CreatedAt = time.Now().UTC()
	proposal.ExpiresAt = proposal.CreatedAt.Add(DefaultExpirationDuration)

	if err := c.store.Create(proposal); err != nil {
		return fmt.Errorf("intake %s: create proposal: %w", c.source, err)
	}
	return nil
}

// HealthCheck reports the channel's health.
func (c *LLMIntakeChannel) HealthCheck() HealthCheckResult {
	if c.store == nil {
		return HealthCheckResult{Status: HealthDown, Message: "no proposal store configured"}
	}
	return HealthCheckResult{Status: HealthOK}
}
