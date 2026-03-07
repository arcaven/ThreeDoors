package mcp

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimitConfig holds rate limiting parameters.
type RateLimitConfig struct {
	GlobalPerMin    int
	ProposalsPerMin int
	QueriesPerMin   int
	MaxPending      int
	BurstAllowance  int
}

// DefaultRateLimitConfig returns the standard rate limits.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		GlobalPerMin:    100,
		ProposalsPerMin: 20,
		QueriesPerMin:   60,
		MaxPending:      5,
		BurstAllowance:  10,
	}
}

// RateLimiter enforces per-connection rate limits on MCP requests.
type RateLimiter struct {
	global    *rate.Limiter
	proposals *rate.Limiter
	queries   *rate.Limiter
	config    RateLimitConfig
}

// NewRateLimiter creates a rate limiter middleware with the given config.
func NewRateLimiter(cfg RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		global:    rate.NewLimiter(rate.Limit(float64(cfg.GlobalPerMin)/60.0), cfg.BurstAllowance),
		proposals: rate.NewLimiter(rate.Limit(float64(cfg.ProposalsPerMin)/60.0), cfg.BurstAllowance),
		queries:   rate.NewLimiter(rate.Limit(float64(cfg.QueriesPerMin)/60.0), cfg.BurstAllowance),
		config:    cfg,
	}
}

// Middleware returns the middleware func for this rate limiter.
func (rl *RateLimiter) Middleware() Middleware {
	return func(next Handler) Handler {
		return func(req *Request) *Response {
			if !rl.global.Allow() {
				retryAfter := math.Ceil(60.0 / float64(rl.config.GlobalPerMin))
				return NewErrorResponseWithData(req.ID, CodeRateLimited,
					"rate limit exceeded",
					map[string]any{"retry_after_seconds": retryAfter})
			}

			if isProposalMethod(req.Method) {
				if !rl.proposals.Allow() {
					retryAfter := math.Ceil(60.0 / float64(rl.config.ProposalsPerMin))
					return NewErrorResponseWithData(req.ID, CodeRateLimited,
						"proposal rate limit exceeded",
						map[string]any{"retry_after_seconds": retryAfter})
				}
			}

			if isQueryMethod(req.Method) {
				if !rl.queries.Allow() {
					retryAfter := math.Ceil(60.0 / float64(rl.config.QueriesPerMin))
					return NewErrorResponseWithData(req.ID, CodeRateLimited,
						"query rate limit exceeded",
						map[string]any{"retry_after_seconds": retryAfter})
				}
			}

			return next(req)
		}
	}
}

func isProposalMethod(method string) bool {
	return strings.HasPrefix(method, "tools/call") && strings.Contains(method, "propos")
}

func isQueryMethod(method string) bool {
	switch method {
	case "resources/list", "resources/read", "tools/list", "prompts/list":
		return true
	}
	return strings.HasPrefix(method, "tools/call")
}

// AuditEntry represents a single audit log record.
type AuditEntry struct {
	Timestamp  time.Time `json:"ts"`
	RequestID  string    `json:"req_id"`
	Tool       string    `json:"tool"`
	Args       any       `json:"args,omitempty"`
	Result     string    `json:"result"`
	Source     string    `json:"source,omitempty"`
	DurationMS int64     `json:"duration_ms"`
	Error      string    `json:"error,omitempty"`
	PrevHash   string    `json:"prev_hash"`
}

// AuditLogger logs every MCP request to a JSONL file with hash chaining.
type AuditLogger struct {
	dir       string
	retention time.Duration
	mu        sync.Mutex
	prevHash  string
	nowFunc   func() time.Time
}

// NewAuditLogger creates an audit logger writing to the given directory.
func NewAuditLogger(dir string) *AuditLogger {
	return &AuditLogger{
		dir:       dir,
		retention: 30 * 24 * time.Hour,
		nowFunc:   func() time.Time { return time.Now().UTC() },
	}
}

// Middleware returns the middleware func for audit logging.
func (al *AuditLogger) Middleware() Middleware {
	return func(next Handler) Handler {
		return func(req *Request) *Response {
			start := al.nowFunc()
			resp := next(req)

			result := "ok"
			var errMsg string
			if resp != nil && resp.Error != nil {
				if resp.Error.Code == CodeRateLimited {
					result = "rate_limited"
				} else {
					result = "error"
				}
				errMsg = resp.Error.Message
			}

			entry := AuditEntry{
				Timestamp:  start,
				RequestID:  formatRequestID(req.ID),
				Tool:       req.Method,
				Result:     result,
				DurationMS: al.nowFunc().Sub(start).Milliseconds(),
				Error:      errMsg,
			}

			if req.Params != nil {
				var args any
				if err := json.Unmarshal(req.Params, &args); err == nil {
					entry.Args = args
				}
			}

			al.writeEntry(entry)
			return resp
		}
	}
}

func formatRequestID(id json.RawMessage) string {
	if id == nil {
		return ""
	}
	return string(id)
}

func (al *AuditLogger) writeEntry(entry AuditEntry) {
	al.mu.Lock()
	defer al.mu.Unlock()

	entry.PrevHash = al.prevHash

	data, err := json.Marshal(entry)
	if err != nil {
		return
	}

	hash := sha256.Sum256(data)
	al.prevHash = hex.EncodeToString(hash[:])

	filename := fmt.Sprintf("mcp-audit-%s.jsonl", entry.Timestamp.Format("2006-01-02"))
	logPath := filepath.Join(al.dir, filename)

	if err := os.MkdirAll(al.dir, 0o700); err != nil {
		return
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	data = append(data, '\n')
	_, _ = f.Write(data)
}

// RotateLogs removes audit log files older than the retention period.
func (al *AuditLogger) RotateLogs() error {
	al.mu.Lock()
	defer al.mu.Unlock()

	entries, err := os.ReadDir(al.dir)
	if err != nil {
		return fmt.Errorf("read audit dir: %w", err)
	}

	cutoff := al.nowFunc().Add(-al.retention)
	for _, e := range entries {
		if e.IsDir() || !strings.HasPrefix(e.Name(), "mcp-audit-") || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}

		dateStr := strings.TrimPrefix(e.Name(), "mcp-audit-")
		dateStr = strings.TrimSuffix(dateStr, ".jsonl")
		fileDate, parseErr := time.Parse("2006-01-02", dateStr)
		if parseErr != nil {
			continue
		}

		if fileDate.Before(cutoff) {
			path := filepath.Join(al.dir, e.Name())
			if removeErr := os.Remove(path); removeErr != nil {
				return fmt.Errorf("remove old audit log %s: %w", e.Name(), removeErr)
			}
		}
	}
	return nil
}

var uuidV4Re = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// SchemaValidator validates MCP request inputs.
func SchemaValidator() Middleware {
	return func(next Handler) Handler {
		return func(req *Request) *Response {
			if req.Params == nil {
				return next(req)
			}

			var params map[string]any
			if err := json.Unmarshal(req.Params, &params); err != nil {
				return next(req)
			}

			if err := validateParams(params); err != nil {
				return NewErrorResponse(req.ID, CodeInvalidParams, err.Error())
			}

			return next(req)
		}
	}
}

func validateParams(params map[string]any) error {
	for key, val := range params {
		switch {
		case isTaskIDField(key):
			s, ok := val.(string)
			if !ok {
				return fmt.Errorf("field %q: expected string", key)
			}
			if !uuidV4Re.MatchString(strings.ToLower(s)) {
				return fmt.Errorf("field %q: invalid UUID v4 format", key)
			}

		case isTextField(key):
			s, ok := val.(string)
			if !ok {
				continue
			}
			if len(s) > 500 {
				return fmt.Errorf("field %q: exceeds 500 character limit (%d chars)", key, len(s))
			}

		case key == "status":
			s, ok := val.(string)
			if !ok {
				return fmt.Errorf("field %q: expected string", key)
			}
			if !isValidStatus(s) {
				return fmt.Errorf("field %q: invalid status %q", key, s)
			}

		case isTimestampField(key):
			s, ok := val.(string)
			if !ok {
				continue
			}
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return fmt.Errorf("field %q: invalid timestamp format", key)
			}
			if t.After(time.Now().UTC().Add(1 * time.Minute)) {
				return fmt.Errorf("field %q: timestamp is in the future", key)
			}
		}
	}
	return nil
}

func isTaskIDField(key string) bool {
	return key == "task_id" || key == "taskId" || key == "id"
}

func isTextField(key string) bool {
	return key == "text" || key == "description" || key == "comment" || key == "reason" || key == "title" || key == "name"
}

func isTimestampField(key string) bool {
	return key == "timestamp" || key == "created_at" || key == "updated_at" || key == "due_date"
}

func isValidStatus(s string) bool {
	switch s {
	case "todo", "blocked", "in-progress", "in-review", "complete", "deferred", "archived":
		return true
	}
	return false
}

var readOnlyBlockedKeywords = []string{
	"save", "delete", "remove", "update", "create", "write",
}

// ReadOnlyEnforcer blocks requests that would directly call write operations
// outside the proposal workflow.
func ReadOnlyEnforcer() Middleware {
	return func(next Handler) Handler {
		return func(req *Request) *Response {
			if isDirectWriteAttempt(req) {
				return NewErrorResponse(req.ID, CodeReadOnly,
					"direct writes are not allowed — use the proposal workflow")
			}
			return next(req)
		}
	}
}

func isDirectWriteAttempt(req *Request) bool {
	if !strings.HasPrefix(req.Method, "tools/call") {
		return false
	}

	if req.Params == nil {
		return false
	}

	var params map[string]any
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return false
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return false
	}

	lower := strings.ToLower(toolName)
	// Allow proposal-related tools
	if strings.Contains(lower, "propos") {
		return false
	}

	for _, keyword := range readOnlyBlockedKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}

// NewErrorResponseWithData creates an error JSON-RPC response with additional data.
func NewErrorResponseWithData(id json.RawMessage, code int, message string, data any) *Response {
	return &Response{
		JSONRPC: jsonRPCVersion,
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

// ApplyDefaultMiddleware wires the standard middleware chain in the specified order:
// ReadOnlyEnforcer -> RateLimiter -> AuditLogger -> SchemaValidator -> coreHandler
func ApplyDefaultMiddleware(server *MCPServer, auditDir string) *AuditLogger {
	al := NewAuditLogger(auditDir)
	rl := NewRateLimiter(DefaultRateLimitConfig())

	// Applied in reverse order — first-added is outermost.
	// Chain: ReadOnlyEnforcer -> RateLimiter -> AuditLogger -> SchemaValidator -> handler
	server.Use(SchemaValidator())
	server.Use(al.Middleware())
	server.Use(rl.Middleware())
	server.Use(ReadOnlyEnforcer())

	return al
}
