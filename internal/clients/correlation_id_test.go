package clients

import (
	"testing"
)

func TestCorrelationRequestID(t *testing.T) {
	first := correlationRequestID()

	if first == "" {
		t.Fatal("no correlation request ID generated")
	}

	second := correlationRequestID()
	if first != second {
		t.Fatal("subsequent correlation request ID not the same as the first")
	}
}
