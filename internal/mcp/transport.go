package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Transport defines the interface for MCP server transports.
type Transport interface {
	// Serve starts serving requests. It blocks until the context is cancelled
	// or an unrecoverable error occurs.
	Serve(ctx context.Context, server *MCPServer) error
}

// StdioTransport implements MCP over stdin/stdout using newline-delimited JSON-RPC.
type StdioTransport struct {
	reader io.Reader
	writer io.Writer
}

// NewStdioTransport creates a transport that reads from r and writes to w.
func NewStdioTransport(r io.Reader, w io.Writer) *StdioTransport {
	return &StdioTransport{reader: r, writer: w}
}

// Serve reads JSON-RPC messages from stdin line-by-line and writes responses to stdout.
func (t *StdioTransport) Serve(ctx context.Context, server *MCPServer) error {
	scanner := bufio.NewScanner(t.reader)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024) // 1MB max message

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("read stdin: %w", err)
			}
			return nil // EOF
		}

		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		resp, err := server.HandleRequest(line)
		if err != nil {
			return fmt.Errorf("handle request: %w", err)
		}
		if resp == nil {
			continue // notification — no response
		}

		resp = append(resp, '\n')
		if _, err := t.writer.Write(resp); err != nil {
			return fmt.Errorf("write response: %w", err)
		}
	}
}

// SSETransport implements MCP over HTTP with Server-Sent Events.
type SSETransport struct {
	addr string

	mu      sync.Mutex
	clients map[string]chan []byte
}

// NewSSETransport creates an SSE transport that listens on the given address.
func NewSSETransport(addr string) *SSETransport {
	return &SSETransport{
		addr:    addr,
		clients: make(map[string]chan []byte),
	}
}

// Serve starts the HTTP server with SSE endpoints.
func (t *SSETransport) Serve(ctx context.Context, server *MCPServer) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/sse", t.handleSSE(server))
	mux.HandleFunc("/message", t.handleMessage(server))

	srv := &http.Server{
		Addr:              t.addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown SSE server: %w", err)
		}
		return ctx.Err()
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("SSE server: %w", err)
		}
		return nil
	}
}

func (t *SSETransport) handleSSE(server *MCPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		clientID := uuid.New().String()
		ch := make(chan []byte, 64)

		t.mu.Lock()
		t.clients[clientID] = ch
		t.mu.Unlock()

		defer func() {
			t.mu.Lock()
			delete(t.clients, clientID)
			t.mu.Unlock()
			close(ch)
		}()

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Send the endpoint event so the client knows where to POST messages.
		if _, err := fmt.Fprintf(w, "event: endpoint\ndata: /message?sessionId=%s\n\n", clientID); err != nil {
			return
		}
		flusher.Flush()

		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-ch:
				if !ok {
					return
				}
				if _, err := fmt.Fprintf(w, "event: message\ndata: %s\n\n", data); err != nil {
					return
				}
				flusher.Flush()
			}
		}
	}
}

func (t *SSETransport) handleMessage(server *MCPServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sessionID := r.URL.Query().Get("sessionId")
		if sessionID == "" {
			http.Error(w, "missing sessionId", http.StatusBadRequest)
			return
		}

		t.mu.Lock()
		ch, ok := t.clients[sessionID]
		t.mu.Unlock()
		if !ok {
			http.Error(w, "unknown session", http.StatusNotFound)
			return
		}

		body, err := io.ReadAll(io.LimitReader(r.Body, 1024*1024))
		if err != nil {
			http.Error(w, "read body failed", http.StatusBadRequest)
			return
		}
		defer func() { _ = r.Body.Close() }()

		resp, err := server.HandleRequest(body)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if resp != nil {
			ch <- resp
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

// TransportFromFlags selects the appropriate transport based on CLI flags.
func TransportFromFlags(transportType, addr string, stdin io.Reader, stdout io.Writer) Transport {
	switch transportType {
	case "sse":
		return NewSSETransport(addr)
	default:
		return NewStdioTransport(stdin, stdout)
	}
}

// SSEEvent represents a server-sent event for testing.
type SSEEvent struct {
	Event string
	Data  json.RawMessage
}
