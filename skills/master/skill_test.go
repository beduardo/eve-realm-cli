package master

import (
	"context"
	"encoding/json"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/beduardo/eve-realm-cli/internal/mcpclient"
	mcpv1 "github.com/beduardo/eve-realm-cli/gen/proto/mcp/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// ---------------------------------------------------------------------------
// Mock for MCPClient — defined here for tests.
// MCPClient interface is declared in skill.go (production code).
// ---------------------------------------------------------------------------

// mockMCPClient implements MCPClient with configurable function fields so that
// each test case can supply its own behaviour without global state.
type mockMCPClient struct {
	listToolsFn  func(ctx context.Context) ([]mcpclient.Tool, error)
	invokeToolFn func(ctx context.Context, name, input string) (string, error)
}

func (m *mockMCPClient) ListTools(ctx context.Context) ([]mcpclient.Tool, error) {
	return m.listToolsFn(ctx)
}

func (m *mockMCPClient) InvokeTool(ctx context.Context, name, input string) (string, error) {
	return m.invokeToolFn(ctx, name, input)
}

// ---------------------------------------------------------------------------
// In-process gRPC server for SC-00B integration test
// ---------------------------------------------------------------------------

const intTestBufSize = 1024 * 1024

// integrationFakeMCPServer returns realistic responses for the end-to-end test.
// InvokeTool("ping") returns a JSON payload containing "pong" and an RFC 3339
// timestamp. Any other tool name returns NOT_FOUND.
type integrationFakeMCPServer struct {
	mcpv1.UnimplementedMCPServiceServer
}

func (s *integrationFakeMCPServer) ListTools(_ context.Context, _ *mcpv1.ListToolsRequest) (*mcpv1.ListToolsResponse, error) {
	return &mcpv1.ListToolsResponse{
		Tools: []*mcpv1.ToolDescriptor{
			{
				Name:        "ping",
				Description: "Sends a ping and receives a pong",
				InputSchema: `{"type":"object","properties":{}}`,
			},
		},
	}, nil
}

func (s *integrationFakeMCPServer) InvokeTool(_ context.Context, req *mcpv1.InvokeToolRequest) (*mcpv1.InvokeToolResponse, error) {
	if req.GetName() == "ping" {
		ts := time.Now().UTC().Format(time.RFC3339)
		payload := `{"message":"pong","timestamp":"` + ts + `"}`
		return &mcpv1.InvokeToolResponse{Output: payload}, nil
	}
	return nil, status.Errorf(codes.NotFound, "tool %q not found", req.GetName())
}

// startIntegrationServer starts an in-process gRPC server backed by a bufconn listener.
// Returns the listener and a cleanup function.
func startIntegrationServer(t *testing.T) (*bufconn.Listener, func()) {
	t.Helper()
	lis := bufconn.Listen(intTestBufSize)
	srv := grpc.NewServer()
	mcpv1.RegisterMCPServiceServer(srv, &integrationFakeMCPServer{})
	go func() {
		if err := srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			_ = err
		}
	}()
	cleanup := func() {
		srv.Stop()
		lis.Close()
	}
	return lis, cleanup
}

// dialIntegrationBufconn returns a *grpc.ClientConn that routes through the
// given bufconn listener, bypassing the real network.
func dialIntegrationBufconn(ctx context.Context, lis *bufconn.Listener) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

// ---------------------------------------------------------------------------
// SC-007 — Discovery mode lists all available tools
// ---------------------------------------------------------------------------

func TestDiscoveryMode_ListsAllTools(t *testing.T) {
	tools := []mcpclient.Tool{
		{Name: "ping", Description: "Sends a ping and receives a pong", InputSchema: `{"type":"object","properties":{}}`},
		{Name: "echo", Description: "Echoes the input back", InputSchema: `{"type":"object","properties":{"text":{"type":"string"}}}`},
	}

	mock := &mockMCPClient{
		listToolsFn: func(_ context.Context) ([]mcpclient.Tool, error) {
			return tools, nil
		},
		invokeToolFn: func(_ context.Context, _, _ string) (string, error) {
			t.Fatal("InvokeTool should not be called in discovery mode")
			return "", nil
		},
	}

	// Discovery mode: no arguments.
	output := Run(mock, []string{})

	for _, tool := range tools {
		if !strings.Contains(output, tool.Name) {
			t.Errorf("output missing tool name %q\noutput: %s", tool.Name, output)
		}
		if !strings.Contains(output, tool.Description) {
			t.Errorf("output missing description for tool %q\noutput: %s", tool.Name, output)
		}
		if !strings.Contains(output, tool.InputSchema) {
			t.Errorf("output missing input schema for tool %q\noutput: %s", tool.Name, output)
		}
	}
}

// ---------------------------------------------------------------------------
// SC-008 — Invocation mode returns raw tool response
// ---------------------------------------------------------------------------

func TestInvocationMode_ReturnsPongResponse(t *testing.T) {
	const wantOutput = `{"message":"pong"}`

	mock := &mockMCPClient{
		listToolsFn: func(_ context.Context) ([]mcpclient.Tool, error) {
			t.Fatal("ListTools should not be called in invocation mode")
			return nil, nil
		},
		invokeToolFn: func(_ context.Context, name, input string) (string, error) {
			if name != "ping" {
				t.Errorf("InvokeTool: got name %q, want %q", name, "ping")
			}
			return wantOutput, nil
		},
	}

	output := Run(mock, []string{"ping", "{}"})

	if output != wantOutput {
		t.Errorf("Run(ping) = %q, want %q", output, wantOutput)
	}
}

// ---------------------------------------------------------------------------
// SC-009 — Skill handles unreachable MCP Server gracefully
// ---------------------------------------------------------------------------

func TestDiscoveryMode_ConnectionError(t *testing.T) {
	const addr = "localhost:30051"

	mock := &mockMCPClient{
		listToolsFn: func(_ context.Context) ([]mcpclient.Tool, error) {
			return nil, &mcpclient.ConnectionError{Addr: addr}
		},
		invokeToolFn: func(_ context.Context, _, _ string) (string, error) {
			t.Fatal("InvokeTool should not be called in discovery mode")
			return "", nil
		},
	}

	output := Run(mock, []string{})

	if !strings.Contains(output, addr) {
		t.Errorf("output missing server address %q\noutput: %s", addr, output)
	}
	for _, forbidden := range []string{"rpc error", "Unavailable", "desc ="} {
		if strings.Contains(output, forbidden) {
			t.Errorf("output must not expose gRPC internals: found %q\noutput: %s", forbidden, output)
		}
	}
}

func TestInvocationMode_ConnectionError(t *testing.T) {
	const addr = "localhost:30051"

	mock := &mockMCPClient{
		listToolsFn: func(_ context.Context) ([]mcpclient.Tool, error) {
			t.Fatal("ListTools should not be called in invocation mode when error is on InvokeTool")
			return nil, nil
		},
		invokeToolFn: func(_ context.Context, _, _ string) (string, error) {
			return "", &mcpclient.ConnectionError{Addr: addr}
		},
	}

	output := Run(mock, []string{"ping", "{}"})

	if !strings.Contains(output, addr) {
		t.Errorf("output missing server address %q\noutput: %s", addr, output)
	}
	for _, forbidden := range []string{"rpc error", "Unavailable", "desc ="} {
		if strings.Contains(output, forbidden) {
			t.Errorf("output must not expose gRPC internals: found %q\noutput: %s", forbidden, output)
		}
	}
}

// ---------------------------------------------------------------------------
// SC-00A — Skill suggests alternatives when tool not found
// ---------------------------------------------------------------------------

func TestInvocationMode_ToolNotFound(t *testing.T) {
	availableTools := []mcpclient.Tool{
		{Name: "ping", Description: "Sends a ping and receives a pong", InputSchema: `{}`},
		{Name: "echo", Description: "Echoes the input back", InputSchema: `{}`},
	}

	mock := &mockMCPClient{
		listToolsFn: func(_ context.Context) ([]mcpclient.Tool, error) {
			return availableTools, nil
		},
		invokeToolFn: func(_ context.Context, name, _ string) (string, error) {
			return "", &mcpclient.ToolNotFoundError{Name: name}
		},
	}

	output := Run(mock, []string{"nonexistent", "{}"})

	if !strings.Contains(output, "nonexistent") {
		t.Errorf("output missing the requested tool name %q\noutput: %s", "nonexistent", output)
	}
	for _, available := range []string{"ping", "echo"} {
		if !strings.Contains(output, available) {
			t.Errorf("output missing available tool name %q\noutput: %s", available, output)
		}
	}
}

// ---------------------------------------------------------------------------
// SC-00B — End-to-end ping invocation via master skill (bufconn integration)
// ---------------------------------------------------------------------------

func TestEndToEnd_PingViaMasterSkill(t *testing.T) {
	lis, cleanup := startIntegrationServer(t)
	defer cleanup()

	ctx := context.Background()
	conn, err := dialIntegrationBufconn(ctx, lis)
	if err != nil {
		t.Fatalf("dialIntegrationBufconn: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	// mcpclient.Client satisfies MCPClient — this is the compile-time assertion
	// that the real client implements the interface defined in skill.go.
	client := mcpclient.NewClientWithConn(conn)

	before := time.Now().UTC()
	output := Run(client, []string{"ping", "{}"})
	after := time.Now().UTC()

	// The output must be valid JSON.
	var payload struct {
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
	}
	if err := json.Unmarshal([]byte(output), &payload); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, output)
	}

	// message must be "pong".
	if payload.Message != "pong" {
		t.Errorf("payload.message = %q, want %q", payload.Message, "pong")
	}

	// timestamp must parse as RFC 3339.
	ts, err := time.Parse(time.RFC3339, payload.Timestamp)
	if err != nil {
		t.Fatalf("payload.timestamp %q is not valid RFC 3339: %v", payload.Timestamp, err)
	}

	// timestamp must be within the test window (1-minute tolerance for clock skew).
	const tolerance = time.Minute
	if ts.Before(before.Add(-tolerance)) || ts.After(after.Add(tolerance)) {
		t.Errorf("timestamp %v is outside the expected window [%v, %v]", ts, before.Add(-tolerance), after.Add(tolerance))
	}
}
