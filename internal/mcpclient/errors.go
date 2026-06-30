package mcpclient

import "fmt"

// ConnectionError is returned when the client cannot reach the MCP server at the given address.
type ConnectionError struct {
	Addr string
}

// Error implements the error interface.
func (e *ConnectionError) Error() string {
	return fmt.Sprintf("mcpclient: cannot connect to MCP server at %s", e.Addr)
}

// ToolNotFoundError is returned when the requested tool does not exist on the MCP server.
type ToolNotFoundError struct {
	Name string
}

// Error implements the error interface.
func (e *ToolNotFoundError) Error() string {
	return fmt.Sprintf("mcpclient: tool %q not found on MCP server", e.Name)
}
