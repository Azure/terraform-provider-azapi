package customization

import "testing"

func TestExtractStringField(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		payload := map[string]interface{}{"id": "asst_123"}
		got, err := extractStringField(payload, "id")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if got != "asst_123" {
			t.Fatalf("expected %q, got %q", "asst_123", got)
		}
	})

	t.Run("missing", func(t *testing.T) {
		payload := map[string]interface{}{"not_id": "x"}
		_, err := extractStringField(payload, "id")
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("wrongType", func(t *testing.T) {
		payload := map[string]interface{}{"id": 123}
		_, err := extractStringField(payload, "id")
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}
