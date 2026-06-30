package config

import (
	"os"
	"path/filepath"
	"testing"
)

// ---------------------------------------------------------------------------
// TestLoadHostConfig_YAMLRoundTrip — YAML file is read and field is populated
// ---------------------------------------------------------------------------

func TestLoadHostConfig_YAMLRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "eve-realm.yaml")

	content := "mcp_server_addr: localhost:30051\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := LoadHostConfig(path)
	if err != nil {
		t.Fatalf("LoadHostConfig: unexpected error: %v", err)
	}
	if cfg.MCPServerAddr != "localhost:30051" {
		t.Errorf("MCPServerAddr = %q, want %q", cfg.MCPServerAddr, "localhost:30051")
	}
}

// ---------------------------------------------------------------------------
// TestLoadHostConfig_MissingFile — nonexistent path returns zero-value + nil error
// ---------------------------------------------------------------------------

func TestLoadHostConfig_MissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "does-not-exist.yaml")

	cfg, err := LoadHostConfig(path)
	if err != nil {
		t.Fatalf("LoadHostConfig: expected nil error for missing file, got: %v", err)
	}
	if cfg.MCPServerAddr != "" {
		t.Errorf("MCPServerAddr = %q, want empty string for missing file", cfg.MCPServerAddr)
	}
}

// ---------------------------------------------------------------------------
// TestResolve_EnvVarOverridesYAML — env var wins over the YAML value
// ---------------------------------------------------------------------------

func TestResolve_EnvVarOverridesYAML(t *testing.T) {
	const envKey = "EVE_REALM_MCP_ADDR"
	const envVal = "remote-host:9090"
	const yamlVal = "localhost:30051"

	t.Setenv(envKey, envVal)

	got := Resolve(envKey, yamlVal)
	if got != envVal {
		t.Errorf("Resolve(%q, %q) = %q, want %q", envKey, yamlVal, got, envVal)
	}
}

// ---------------------------------------------------------------------------
// TestResolve_YAMLFallbackWhenEnvEmpty — YAML value used when env var is unset
// ---------------------------------------------------------------------------

func TestResolve_YAMLFallbackWhenEnvEmpty(t *testing.T) {
	const envKey = "EVE_REALM_MCP_ADDR_UNSET_FOR_TEST"
	const yamlVal = "localhost:30051"

	// Ensure the env var is absent.
	t.Setenv(envKey, "")

	got := Resolve(envKey, yamlVal)
	if got != yamlVal {
		t.Errorf("Resolve(%q, %q) = %q, want %q", envKey, yamlVal, got, yamlVal)
	}
}
