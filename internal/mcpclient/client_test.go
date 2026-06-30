package mcpclient

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"
	"time"

	mcpv1 "github.com/beduardo/eve-realm-cli/gen/proto/mcp/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// testConnectionTimeout is the deadline used when dialing an unreachable address.
// Kept short so that connection-error tests do not slow the suite.
const testConnectionTimeout = 2 * time.Second

// fakeMCPServer is an in-process gRPC server for testing.
// ListTools returns a fixed set of tools.
// InvokeTool returns a pong response for "ping", NOT_FOUND for anything else.
type fakeMCPServer struct {
	mcpv1.UnimplementedMCPServiceServer
}

func (f *fakeMCPServer) ListTools(_ context.Context, _ *mcpv1.ListToolsRequest) (*mcpv1.ListToolsResponse, error) {
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

func (f *fakeMCPServer) InvokeTool(_ context.Context, req *mcpv1.InvokeToolRequest) (*mcpv1.InvokeToolResponse, error) {
	if req.GetName() == "ping" {
		return &mcpv1.InvokeToolResponse{Output: `{"message":"pong"}`}, nil
	}
	return nil, status.Errorf(codes.NotFound, "tool %q not found", req.GetName())
}

// startBufconnServer starts an in-process gRPC server backed by a bufconn listener.
// It returns the listener and a cleanup function. The caller is responsible for
// calling cleanup when the test finishes.
func startBufconnServer(t *testing.T) (*bufconn.Listener, func()) {
	t.Helper()
	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	mcpv1.RegisterMCPServiceServer(srv, &fakeMCPServer{})
	go func() {
		if err := srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			// The server stopping on listener close is expected in tests.
			_ = err
		}
	}()
	cleanup := func() {
		srv.Stop()
		lis.Close()
	}
	return lis, cleanup
}

// dialBufconn returns a *grpc.ClientConn that dials the given bufconn listener.
func dialBufconn(ctx context.Context, lis *bufconn.Listener) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

// newTestClient returns a Client connected to the given bufconn listener.
// This helper exercises NewClientWithConn, which Step 5 will provide.
func newTestClient(t *testing.T, lis *bufconn.Listener) *Client {
	t.Helper()
	ctx := context.Background()
	conn, err := dialBufconn(ctx, lis)
	if err != nil {
		t.Fatalf("dialBufconn: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return NewClientWithConn(conn)
}

// ---------------------------------------------------------------------------
// SC-003 — ListTools returns tools from server
// ---------------------------------------------------------------------------

func TestListTools_ReturnsToolsFromServer(t *testing.T) {
	lis, cleanup := startBufconnServer(t)
	defer cleanup()

	client := newTestClient(t, lis)
	tools, err := client.ListTools(context.Background())
	if err != nil {
		t.Fatalf("ListTools: unexpected error: %v", err)
	}
	if len(tools) == 0 {
		t.Fatal("ListTools: expected at least one tool, got none")
	}
	for i, tool := range tools {
		if tool.Name == "" {
			t.Errorf("tools[%d].Name is empty", i)
		}
		if tool.Description == "" {
			t.Errorf("tools[%d].Description is empty", i)
		}
		if tool.InputSchema == "" {
			t.Errorf("tools[%d].InputSchema is empty", i)
		}
	}
}

// ---------------------------------------------------------------------------
// SC-004 — InvokeTool returns a valid JSON response
// ---------------------------------------------------------------------------

func TestInvokeTool_ReturnsPongResponse(t *testing.T) {
	lis, cleanup := startBufconnServer(t)
	defer cleanup()

	client := newTestClient(t, lis)
	output, err := client.InvokeTool(context.Background(), "ping", "{}")
	if err != nil {
		t.Fatalf("InvokeTool(ping): unexpected error: %v", err)
	}
	if output == "" {
		t.Fatal("InvokeTool(ping): expected non-empty JSON output, got empty string")
	}
}

// ---------------------------------------------------------------------------
// SC-005 — ConnectionError when server is unreachable
// ---------------------------------------------------------------------------

func TestListTools_ConnectionError(t *testing.T) {
	client := NewClient("localhost:59999")
	ctx, cancel := context.WithTimeout(context.Background(), testConnectionTimeout)
	defer cancel()

	_, err := client.ListTools(ctx)
	if err == nil {
		t.Fatal("ListTools: expected an error for unreachable address, got nil")
	}

	var ce *ConnectionError
	if !errors.As(err, &ce) {
		t.Errorf("ListTools: expected ConnectionError, got %T: %v", err, err)
	}
	if ce != nil && !strings.Contains(ce.Error(), "localhost:59999") {
		t.Errorf("ConnectionError.Error() = %q, want it to contain address %q", ce.Error(), "localhost:59999")
	}
}

func TestInvokeTool_ConnectionError(t *testing.T) {
	client := NewClient("localhost:59999")
	ctx, cancel := context.WithTimeout(context.Background(), testConnectionTimeout)
	defer cancel()

	_, err := client.InvokeTool(ctx, "ping", "{}")
	if err == nil {
		t.Fatal("InvokeTool: expected an error for unreachable address, got nil")
	}

	var ce *ConnectionError
	if !errors.As(err, &ce) {
		t.Errorf("InvokeTool: expected ConnectionError, got %T: %v", err, err)
	}
	if ce != nil && !strings.Contains(ce.Error(), "localhost:59999") {
		t.Errorf("ConnectionError.Error() = %q, want it to contain address %q", ce.Error(), "localhost:59999")
	}
}

// ---------------------------------------------------------------------------
// SC-006 — ToolNotFoundError for unknown tool (not a ConnectionError)
// ---------------------------------------------------------------------------

func TestInvokeTool_ToolNotFoundError(t *testing.T) {
	lis, cleanup := startBufconnServer(t)
	defer cleanup()

	client := newTestClient(t, lis)
	_, err := client.InvokeTool(context.Background(), "nonexistent", "{}")
	if err == nil {
		t.Fatal("InvokeTool(nonexistent): expected an error, got nil")
	}

	var tfe *ToolNotFoundError
	if !errors.As(err, &tfe) {
		t.Errorf("InvokeTool(nonexistent): expected ToolNotFoundError, got %T: %v", err, err)
	}

	var ce *ConnectionError
	if errors.As(err, &ce) {
		t.Error("InvokeTool(nonexistent): expected NOT a ConnectionError, but errors.As returned true")
	}

	if tfe != nil && !strings.Contains(tfe.Error(), "nonexistent") {
		t.Errorf("ToolNotFoundError.Error() = %q, want it to contain tool name %q", tfe.Error(), "nonexistent")
	}
}

// ---------------------------------------------------------------------------
// AC-6 — DefaultMCPAddr constant and NewClient default address
// ---------------------------------------------------------------------------

func TestNewClient_DefaultAddress(t *testing.T) {
	const wantAddr = "localhost:30051"
	if DefaultMCPAddr != wantAddr {
		t.Errorf("DefaultMCPAddr = %q, want %q", DefaultMCPAddr, wantAddr)
	}

	// NewClient("") must use DefaultMCPAddr.
	c := NewClient("")
	if c.addr != DefaultMCPAddr {
		t.Errorf("NewClient(\"\").addr = %q, want %q", c.addr, DefaultMCPAddr)
	}

	// NewClient with an explicit address must use that address.
	c2 := NewClient("localhost:9999")
	if c2.addr != "localhost:9999" {
		t.Errorf("NewClient(\"localhost:9999\").addr = %q, want \"localhost:9999\"", c2.addr)
	}
}
