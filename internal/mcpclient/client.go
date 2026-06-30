package mcpclient

import (
	"context"

	mcpv1 "github.com/beduardo/eve-realm-cli/gen/proto/mcp/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// DefaultMCPAddr is the address used when NewClient is called with an empty string.
const DefaultMCPAddr = "localhost:30051"

// Tool describes a single tool exposed by the MCP server.
type Tool struct {
	Name        string
	Description string
	InputSchema string
}

// Client is a gRPC client for the MCP server.
type Client struct {
	addr string
	conn *grpc.ClientConn
}

// NewClient creates a Client that dials addr on first use.
// If addr is empty, DefaultMCPAddr is used.
func NewClient(addr string) *Client {
	if addr == "" {
		addr = DefaultMCPAddr
	}
	conn, _ := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return &Client{addr: addr, conn: conn}
}

// NewClientWithConn creates a Client using a pre-established gRPC connection.
// This is intended for testing with bufconn.
func NewClientWithConn(conn *grpc.ClientConn) *Client {
	return &Client{addr: conn.Target(), conn: conn}
}

// ListTools calls the MCPService.ListTools RPC and maps the response to a []Tool slice.
// A gRPC Unavailable error is wrapped into ConnectionError.
func (c *Client) ListTools(ctx context.Context) ([]Tool, error) {
	stub := mcpv1.NewMCPServiceClient(c.conn)
	resp, err := stub.ListTools(ctx, &mcpv1.ListToolsRequest{})
	if err != nil {
		return nil, c.mapError(err, "")
	}
	tools := make([]Tool, 0, len(resp.GetTools()))
	for _, td := range resp.GetTools() {
		tools = append(tools, Tool{
			Name:        td.GetName(),
			Description: td.GetDescription(),
			InputSchema: td.GetInputSchema(),
		})
	}
	return tools, nil
}

// InvokeTool calls the MCPService.InvokeTool RPC and returns the output field.
// A gRPC Unavailable error is wrapped into ConnectionError.
// A gRPC NotFound error is wrapped into ToolNotFoundError.
func (c *Client) InvokeTool(ctx context.Context, name, input string) (string, error) {
	stub := mcpv1.NewMCPServiceClient(c.conn)
	resp, err := stub.InvokeTool(ctx, &mcpv1.InvokeToolRequest{
		Name:  name,
		Input: input,
	})
	if err != nil {
		return "", c.mapError(err, name)
	}
	return resp.GetOutput(), nil
}

// mapError converts gRPC status errors into the typed errors defined in errors.go.
// toolName is only relevant for ToolNotFoundError; pass empty string for ListTools.
func (c *Client) mapError(err error, toolName string) error {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.NotFound:
			return &ToolNotFoundError{Name: toolName}
		case codes.Unavailable, codes.DeadlineExceeded:
			return &ConnectionError{Addr: c.addr}
		}
	}
	// Context deadline exceeded (not a gRPC status) or any other dial/transport error.
	return &ConnectionError{Addr: c.addr}
}
