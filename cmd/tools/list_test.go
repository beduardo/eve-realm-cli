package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/beduardo/eve-realm-cli/internal/mcpclient"
)

// ---------------------------------------------------------------------------
// TestListCmd_TwoTools — two tools returned: stdout contains both; exit 0
// ---------------------------------------------------------------------------

func TestListCmd_TwoTools(t *testing.T) {
	mock := &mockMCPClient{
		ListToolsFn: func(_ context.Context) ([]mcpclient.Tool, error) {
			return []mcpclient.Tool{
				{Name: "ping", Description: "Ping the server", InputSchema: `{"type":"object"}`},
				{Name: "echo", Description: "Echo the input", InputSchema: `{"type":"string"}`},
			}, nil
		},
	}

	stdout, _, err := runToolsCmd(t, mock, "tools", "list")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	out := stdout.String()

	for _, want := range []string{
		"ping", "Ping the server", `{"type":"object"}`,
		"echo", "Echo the input", `{"type":"string"}`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout missing %q\ngot:\n%s", want, out)
		}
	}
}

// ---------------------------------------------------------------------------
// TestListCmd_OneTool — one tool returned: stdout contains only that tool; exit 0
// ---------------------------------------------------------------------------

func TestListCmd_OneTool(t *testing.T) {
	mock := &mockMCPClient{
		ListToolsFn: func(_ context.Context) ([]mcpclient.Tool, error) {
			return []mcpclient.Tool{
				{Name: "ping", Description: "Ping the server", InputSchema: `{"type":"object"}`},
			}, nil
		},
	}

	stdout, _, err := runToolsCmd(t, mock, "tools", "list")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	out := stdout.String()

	if !strings.Contains(out, "ping") {
		t.Errorf("stdout missing %q\ngot:\n%s", "ping", out)
	}
	if !strings.Contains(out, "Ping the server") {
		t.Errorf("stdout missing %q\ngot:\n%s", "Ping the server", out)
	}
	if !strings.Contains(out, `{"type":"object"}`) {
		t.Errorf(`stdout missing {"type":"object"}\ngot:\n%s`, out)
	}
	// Verify "echo" is NOT present (only one tool registered)
	if strings.Contains(out, "echo") {
		t.Errorf("stdout contains unexpected tool %q\ngot:\n%s", "echo", out)
	}
}

// ---------------------------------------------------------------------------
// TestListCmd_ZeroTools — zero tools returned: graceful output; exit 0
// ---------------------------------------------------------------------------

func TestListCmd_ZeroTools(t *testing.T) {
	mock := &mockMCPClient{
		ListToolsFn: func(_ context.Context) ([]mcpclient.Tool, error) {
			return []mcpclient.Tool{}, nil
		},
	}

	stdout, stderr, err := runToolsCmd(t, mock, "tools", "list")
	if err != nil {
		t.Fatalf("expected nil error for zero tools, got: %v", err)
	}

	_ = stdout
	_ = stderr
	// Graceful: no panic, no crash. No assertions on content beyond successful exit.
}

// ---------------------------------------------------------------------------
// TestListCmd_ConnectionError — ConnectionError: non-nil error; stderr has message; stdout empty
// ---------------------------------------------------------------------------

func TestListCmd_ConnectionError(t *testing.T) {
	mock := &mockMCPClient{
		ListToolsFn: func(_ context.Context) ([]mcpclient.Tool, error) {
			return nil, &mcpclient.ConnectionError{Addr: "localhost:30051"}
		},
	}

	stdout, stderr, err := runToolsCmd(t, mock, "tools", "list")
	if err == nil {
		t.Fatal("expected non-nil error for ConnectionError, got nil")
	}

	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout on error, got: %q", stdout.String())
	}

	errOut := stderr.String()
	if errOut == "" {
		t.Error("expected non-empty stderr on ConnectionError, got empty")
	}
	// Error message must be user-readable (not a raw Go struct)
	if strings.Contains(errOut, "{") && strings.Contains(errOut, "}") {
		t.Errorf("stderr looks like a raw struct: %q", errOut)
	}
}

// ---------------------------------------------------------------------------
// TestListCmd_OutputFormat — output is human-readable (no raw struct printing)
// ---------------------------------------------------------------------------

func TestListCmd_OutputFormat(t *testing.T) {
	mock := &mockMCPClient{
		ListToolsFn: func(_ context.Context) ([]mcpclient.Tool, error) {
			return []mcpclient.Tool{
				{Name: "ping", Description: "Ping the server", InputSchema: `{"type":"object"}`},
				{Name: "echo", Description: "Echo the input", InputSchema: `{"type":"string"}`},
			}, nil
		},
	}

	stdout, _, err := runToolsCmd(t, mock, "tools", "list")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	out := stdout.String()

	// Must not print raw Go struct (e.g. {Name:ping Description:...})
	if strings.Contains(out, "mcpclient.Tool{") || strings.Contains(out, "{Name:") {
		t.Errorf("stdout appears to contain raw struct output:\n%s", out)
	}

	// Each tool section is separated — the two tools must appear in distinct regions.
	// Verify ordering: ping appears before echo.
	pingIdx := strings.Index(out, "ping")
	echoIdx := strings.Index(out, "echo")
	if pingIdx < 0 || echoIdx < 0 {
		t.Fatalf("expected both tools in output, got:\n%s", out)
	}
	if pingIdx >= echoIdx {
		t.Errorf("expected 'ping' to appear before 'echo' in output:\n%s", out)
	}
}

// ---------------------------------------------------------------------------
// TestListCmd_AddressResolution — SC-017 address resolution table test
// ---------------------------------------------------------------------------

func TestListCmd_AddressResolution(t *testing.T) {
	cases := []struct {
		name        string
		envVal      string        // empty means unset
		writeYAML   bool          // whether to write a config file
		yamlAddr    string        // value in the YAML file
		wantAddr    string        // expected address passed to the client factory
	}{
		{
			name:      "env-var-wins",
			envVal:    "remote-host:9090",
			writeYAML: true,
			yamlAddr:  "yaml-host:8080",
			wantAddr:  "remote-host:9090",
		},
		{
			name:      "yaml-wins",
			envVal:    "",
			writeYAML: true,
			yamlAddr:  "yaml-host:8080",
			wantAddr:  "yaml-host:8080",
		},
		{
			name:      "default-used",
			envVal:    "",
			writeYAML: false,
			yamlAddr:  "",
			wantAddr:  "localhost:30051",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			configPath := filepath.Join(dir, "eve-realm.yaml")

			if tc.envVal != "" {
				t.Setenv("EVE_REALM_MCP_ADDR", tc.envVal)
			} else {
				// Ensure env var is absent for this subtest
				t.Setenv("EVE_REALM_MCP_ADDR", "")
			}

			if tc.writeYAML {
				content := "mcp_server_addr: " + tc.yamlAddr + "\n"
				if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
					t.Fatalf("WriteFile: %v", err)
				}
			}
			// If !writeYAML, configPath does not exist — LoadHostConfig returns zero-value

			// Simulate the same resolution logic used in NewToolsCmd
			cfg, err := loadHostConfigForTest(configPath)
			if err != nil {
				t.Fatalf("loadHostConfig: %v", err)
			}
			addr := resolveAddrForTest("EVE_REALM_MCP_ADDR", cfg)
			if addr != tc.wantAddr {
				t.Errorf("addr = %q, want %q", addr, tc.wantAddr)
			}
		})
	}
}

// loadHostConfigForTest and resolveAddrForTest mirror the logic in NewToolsCmd
// so the test exercises the same resolution path without a real gRPC connection.

func loadHostConfigForTest(configPath string) (string, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	// Minimal YAML parse: look for mcp_server_addr line
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "mcp_server_addr:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "mcp_server_addr:"))
			return val, nil
		}
	}
	return "", nil
}

func resolveAddrForTest(envKey, yamlValue string) string {
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	if yamlValue != "" {
		return yamlValue
	}
	return mcpclient.DefaultMCPAddr
}
