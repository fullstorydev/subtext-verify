// Package tunnelmcp implements a hand-rolled MCP stdio server that exposes
// tunnel-connect, tunnel-disconnect, and tunnel-status tools. It manages an
// in-process map of active tunnels backed by the tunnel.Client engine,
// mirroring the semantics of the TypeScript @fullstory/subtext-tunnel package.
package tunnelmcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/fullstorydev/subtext-verify/cli/internal/tunnel"
)

// Server is a stdio MCP server. Create one with New, then call Serve.
type Server struct {
	version string

	mu      sync.Mutex
	tunnels map[string]*runningTunnel // keyed by relay-assigned tunnelId
}

// runningTunnel holds the live state of one in-process tunnel connection.
type runningTunnel struct {
	cancel   context.CancelFunc
	relayURL string

	mu    sync.Mutex
	state tunnel.TunnelState
	info  tunnel.ReadyInfo
}

// New returns a Server. version is embedded in the MCP server-info handshake.
func New(version string) *Server {
	return &Server{
		version: version,
		tunnels: make(map[string]*runningTunnel),
	}
}

// Serve reads MCP JSON-RPC from r and writes responses to w until r returns
// EOF or ctx is cancelled. All tunnel goroutines are cancelled when Serve
// returns. Logs go to stderr.
func (s *Server) Serve(ctx context.Context, r io.Reader, w io.Writer) error {
	dec := json.NewDecoder(r)
	enc := json.NewEncoder(w)
	var writeMu sync.Mutex

	write := func(v any) {
		writeMu.Lock()
		defer writeMu.Unlock()
		if err := enc.Encode(v); err != nil && !errors.Is(err, io.ErrClosedPipe) {
			logf("encode error: %v", err)
		}
	}

	defer s.disconnectAll()

	for {
		var req rpcRequest
		if err := dec.Decode(&req); err != nil {
			if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
				return nil
			}
			if ctx.Err() != nil {
				return nil
			}
			// Malformed JSON; send parse error and continue.
			write(rpcErr(nil, -32700, "parse error"))
			continue
		}

		// Notifications (no id field) require no response.
		if req.ID == nil {
			continue
		}

		resp := s.handle(ctx, req)
		write(resp)
	}
}

// handle dispatches a single JSON-RPC request and returns the response object.
func (s *Server) handle(ctx context.Context, req rpcRequest) rpcResponse {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(ctx, req)
	default:
		return rpcErr(req.ID, -32601, "method not found: "+req.Method)
	}
}

// ---- MCP protocol handlers ----

func (s *Server) handleInitialize(req rpcRequest) rpcResponse {
	// Echo back the protocol version the client sent, or default to 2024-11-05.
	var params struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	if req.Params != nil {
		_ = json.Unmarshal(req.Params, &params)
	}
	if params.ProtocolVersion == "" {
		params.ProtocolVersion = "2024-11-05"
	}

	return rpcOK(req.ID, map[string]any{
		"protocolVersion": params.ProtocolVersion,
		"capabilities":    map[string]any{"tools": map[string]any{}},
		"serverInfo":      map[string]any{"name": "subtext_tunnel", "version": s.version},
	})
}

func (s *Server) handleToolsList(req rpcRequest) rpcResponse {
	return rpcOK(req.ID, map[string]any{"tools": mcpTools})
}

func (s *Server) handleToolsCall(ctx context.Context, req rpcRequest) rpcResponse {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if req.Params == nil {
		return rpcErr(req.ID, -32602, "missing params")
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return rpcErr(req.ID, -32602, "invalid params: "+err.Error())
	}

	switch params.Name {
	case "tunnel-connect":
		return s.toolConnect(ctx, req.ID, params.Arguments)
	case "tunnel-disconnect":
		return s.toolDisconnect(req.ID, params.Arguments)
	case "tunnel-status":
		return s.toolStatus(req.ID)
	default:
		return rpcErr(req.ID, -32602, "unknown tool: "+params.Name)
	}
}

// ---- Tool handlers ----

func (s *Server) toolConnect(ctx context.Context, id json.RawMessage, rawArgs json.RawMessage) rpcResponse {
	var args struct {
		RelayURL       string   `json:"relayUrl"`
		ConnectionID   string   `json:"connectionId"`
		AllowedOrigins []string `json:"allowedOrigins"`
	}
	if rawArgs != nil {
		if err := json.Unmarshal(rawArgs, &args); err != nil {
			return toolError(id, "invalid arguments: "+err.Error())
		}
	}
	if args.RelayURL == "" {
		return toolError(id, "relayUrl is required")
	}

	// Parse and canonicalize allowlist up front so the response can carry back
	// any rewrites and the caller learns what was actually registered.
	patterns, rewrites, err := tunnel.ParseOriginPatternsWithRewrites(args.AllowedOrigins)
	if err != nil {
		return toolError(id, err.Error())
	}

	opts := tunnel.ClientOptions{
		RelayURL:          args.RelayURL,
		ConnectionID:      args.ConnectionID,
		AllowedOrigins:    patterns,
		AllowedOriginsRaw: args.AllowedOrigins,
		Log:               func(f string, a ...any) { logf(f, a...) },
	}

	tunnelCtx, cancel := context.WithCancel(ctx)
	rt := &runningTunnel{
		cancel:   cancel,
		relayURL: args.RelayURL,
		state:    tunnel.StateConnecting,
	}

	client := tunnel.NewClient(opts)
	readyCh := make(chan tunnel.ReadyInfo, 1)

	onReady := func(info tunnel.ReadyInfo) {
		rt.mu.Lock()
		rt.info = info
		rt.state = tunnel.StateReady
		rt.mu.Unlock()

		select {
		case readyCh <- info:
			// First ready: register in the map keyed by relay-assigned tunnelId.
			s.mu.Lock()
			s.tunnels[info.TunnelID] = rt
			s.mu.Unlock()
		default:
		}
	}

	errCh := make(chan error, 1)
	go func() {
		err := client.Run(tunnelCtx, onReady)
		rt.mu.Lock()
		rt.state = tunnel.StateDisconnected
		info := rt.info
		rt.mu.Unlock()

		// Clean up from the map when the goroutine exits.
		if info.TunnelID != "" {
			s.mu.Lock()
			delete(s.tunnels, info.TunnelID)
			s.mu.Unlock()
		}
		errCh <- err
	}()

	// Wait up to 5 s for ready, a definitive error, or timeout.
	deadline := time.NewTimer(5 * time.Second)
	defer deadline.Stop()

	select {
	case info := <-readyCh:
		result := map[string]any{
			"state":        string(tunnel.StateReady),
			"tunnelId":     info.TunnelID,
			"connectionId": info.ConnectionID,
			"traceId":      info.TraceID,
			"relayUrl":     args.RelayURL,
		}
		if len(rewrites) > 0 {
			result["canonicalized"] = rewrites
		}
		return toolOK(id, result)

	case err := <-errCh:
		cancel()
		if errors.Is(err, tunnel.ErrNeedLiveTunnel) {
			return toolError(id, "resume token rejected; call live-tunnel to get a fresh relay URL")
		}
		if err != nil {
			return toolError(id, err.Error())
		}
		return toolError(id, "tunnel disconnected before becoming ready")

	case <-deadline.C:
		// Return current state even if not yet ready (mirrors TS behaviour).
		rt.mu.Lock()
		st := rt.state
		rt.mu.Unlock()
		result := map[string]any{
			"state":    string(st),
			"relayUrl": args.RelayURL,
		}
		if len(rewrites) > 0 {
			result["canonicalized"] = rewrites
		}
		return toolOK(id, result)

	case <-ctx.Done():
		cancel()
		return toolError(id, "server shutting down")
	}
}

func (s *Server) toolDisconnect(id json.RawMessage, rawArgs json.RawMessage) rpcResponse {
	var args struct {
		TunnelID string `json:"tunnelId"`
	}
	if rawArgs != nil {
		_ = json.Unmarshal(rawArgs, &args)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if args.TunnelID != "" {
		rt, ok := s.tunnels[args.TunnelID]
		if !ok {
			return toolError(id, fmt.Sprintf("no tunnel with id %s", args.TunnelID))
		}
		rt.cancel()
		delete(s.tunnels, args.TunnelID)
		return toolOK(id, map[string]any{"disconnected": args.TunnelID})
	}

	return toolOK(id, map[string]any{"disconnected": s.cancelAll()})
}

func (s *Server) toolStatus(id json.RawMessage) rpcResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	type entry struct {
		TunnelID string `json:"tunnelId"`
		State    string `json:"state"`
		TraceID  string `json:"traceId"`
	}
	tunnels := make([]entry, 0, len(s.tunnels))
	for tid, rt := range s.tunnels {
		rt.mu.Lock()
		st := rt.state
		traceID := rt.info.TraceID
		rt.mu.Unlock()
		tunnels = append(tunnels, entry{
			TunnelID: tid,
			State:    string(st),
			TraceID:  traceID,
		})
	}
	return toolOK(id, map[string]any{"tunnels": tunnels, "count": len(tunnels)})
}

// cancelAll cancels all tunnels and returns their IDs. s.mu must be held.
func (s *Server) cancelAll() []string {
	ids := make([]string, 0, len(s.tunnels))
	for tid, rt := range s.tunnels {
		rt.cancel()
		ids = append(ids, tid)
	}
	s.tunnels = make(map[string]*runningTunnel)
	return ids
}

func (s *Server) disconnectAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cancelAll()
}

// ---- Tool definitions ----

// mcpTools is the static tool list returned on every tools/list request.
var mcpTools = []map[string]any{
	{
		"name": "tunnel-connect",
		"description": "Connect a tunnel to the relay. Multiple tunnels can be active simultaneously. " +
			"Call live-tunnel on the subtext MCP server first to obtain the relayUrl.",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"relayUrl": map[string]any{
					"type":        "string",
					"description": "WebSocket URL of the relay (from live-tunnel).",
				},
				"connectionId": map[string]any{
					"type": "string",
					"description": "Connection ID to bind this tunnel to. Required for connection-first flow " +
						"(pass the connection_id from live-connect). Omit for tunnel-first flow " +
						"(the server mints one and returns it in the response).",
				},
				"allowedOrigins": map[string]any{
					"type": "array",
					"items": map[string]any{"type": "string"},
					"description": "Optional per-tunnel origin allowlist. Each entry is a bare " +
						"`host:port` (no scheme). DNS hosts implicitly match their subdomains, " +
						"so list the trunk you want to allow: `example.test:8043` covers " +
						"`app.example.test:8043`, `oauthtest.example.test:8043`, etc. " +
						"Hosts must be loopback-class (localhost, 127.x, ::1, *.test, *.localhost). " +
						"IP literals match exactly. The response includes a `canonicalized` field " +
						"listing any entries that were rewritten.",
				},
			},
			"required": []string{"relayUrl"},
		},
	},
	{
		"name":        "tunnel-disconnect",
		"description": "Disconnect a specific tunnel by its tunnelId. If no tunnelId is given, disconnects all tunnels.",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"tunnelId": map[string]any{
					"type":        "string",
					"description": "The tunnelId to disconnect (from tunnel-connect response). Omit to disconnect all.",
				},
			},
			"additionalProperties": false,
		},
	},
	{
		"name":        "tunnel-status",
		"description": "Returns the status of all active tunnels.",
		"inputSchema": map[string]any{
			"type":                 "object",
			"properties":           map[string]any{},
			"additionalProperties": false,
		},
	},
}

// ---- JSON-RPC helpers ----

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"` // number, string, or null; nil means notification
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcErrorObj    `json:"error,omitempty"`
}

type rpcErrorObj struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func rpcOK(id json.RawMessage, result any) rpcResponse {
	return rpcResponse{JSONRPC: "2.0", ID: id, Result: result}
}

func rpcErr(id json.RawMessage, code int, msg string) rpcResponse {
	return rpcResponse{JSONRPC: "2.0", ID: id, Error: &rpcErrorObj{Code: code, Message: msg}}
}

// toolOK wraps a successful tool result as per MCP spec.
func toolOK(id json.RawMessage, data any) rpcResponse {
	text, _ := json.Marshal(data)
	return rpcOK(id, map[string]any{
		"content": []map[string]any{{"type": "text", "text": string(text)}},
		"isError": false,
	})
}

// toolError wraps a tool-level error as per MCP spec (isError=true).
func toolError(id json.RawMessage, msg string) rpcResponse {
	text, _ := json.Marshal(map[string]string{"error": msg})
	return rpcOK(id, map[string]any{
		"content": []map[string]any{{"type": "text", "text": string(text)}},
		"isError": true,
	})
}

func logf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "[subtext-tunnel] "+format+"\n", args...)
}
