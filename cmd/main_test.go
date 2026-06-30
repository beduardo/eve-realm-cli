package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionDefaults(t *testing.T) {
	if Version == "" {
		t.Fatal("Version must not be empty")
	}
	if GitHash == "" {
		t.Fatal("GitHash must not be empty")
	}
	if BuildDate == "" {
		t.Fatal("BuildDate must not be empty")
	}
}

func TestRootCommand_SilenceErrors(t *testing.T) {
	cmd := newRootCmd()
	if !cmd.SilenceErrors {
		t.Error("root command must have SilenceErrors = true")
	}
}

func TestVersionCommand_Output(t *testing.T) {
	Version = "1.2.3"
	GitHash = "abc1234"
	BuildDate = "2026-06-30"

	cmd := newRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "1.2.3") {
		t.Errorf("version output missing version string, got: %q", got)
	}
	if !strings.Contains(got, "abc1234") {
		t.Errorf("version output missing git hash, got: %q", got)
	}
	if !strings.Contains(got, "2026-06-30") {
		t.Errorf("version output missing build date, got: %q", got)
	}
}
