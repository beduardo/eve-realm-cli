package mcpclient

import (
	"errors"
	"strings"
	"testing"
)

func TestConnectionError_ErrorIncludesAddr(t *testing.T) {
	e := &ConnectionError{Addr: "localhost:50051"}
	if !strings.Contains(e.Error(), "localhost:50051") {
		t.Errorf("ConnectionError.Error() = %q, want it to contain address", e.Error())
	}
}

func TestToolNotFoundError_ErrorIncludesName(t *testing.T) {
	e := &ToolNotFoundError{Name: "ping"}
	if !strings.Contains(e.Error(), "ping") {
		t.Errorf("ToolNotFoundError.Error() = %q, want it to contain tool name", e.Error())
	}
}

func TestConnectionError_ErrorsAs(t *testing.T) {
	var err error = &ConnectionError{Addr: "localhost:50051"}

	var ce *ConnectionError
	if !errors.As(err, &ce) {
		t.Error("errors.As(connectionErr, &ConnectionError{}) returned false, want true")
	}

	var te *ToolNotFoundError
	if errors.As(err, &te) {
		t.Error("errors.As(connectionErr, &ToolNotFoundError{}) returned true, want false")
	}
}

func TestToolNotFoundError_ErrorsAs(t *testing.T) {
	var err error = &ToolNotFoundError{Name: "ping"}

	var te *ToolNotFoundError
	if !errors.As(err, &te) {
		t.Error("errors.As(toolNotFoundErr, &ToolNotFoundError{}) returned false, want true")
	}

	var ce *ConnectionError
	if errors.As(err, &ce) {
		t.Error("errors.As(toolNotFoundErr, &ConnectionError{}) returned true, want false")
	}
}
