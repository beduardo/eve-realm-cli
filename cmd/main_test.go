package main

import "testing"

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
