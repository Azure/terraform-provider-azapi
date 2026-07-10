package provider

import "testing"

func TestGetOIDCTokenFilePathFromEnv(t *testing.T) {
	t.Run("prefers ARM_OIDC_TOKEN_FILE_PATH", func(t *testing.T) {
		t.Setenv("ARM_OIDC_TOKEN_FILE_PATH", "/tmp/arm-token")
		t.Setenv("AZURE_FEDERATED_TOKEN_FILE", "/tmp/azure-token")

		if got := getOIDCTokenFilePathFromEnv(false); got != "/tmp/arm-token" {
			t.Fatalf("expected ARM_OIDC_TOKEN_FILE_PATH to be used, got %q", got)
		}
	})

	t.Run("uses AZURE_FEDERATED_TOKEN_FILE only when AKS workload identity is enabled", func(t *testing.T) {
		t.Setenv("ARM_OIDC_TOKEN_FILE_PATH", "")
		t.Setenv("AZURE_FEDERATED_TOKEN_FILE", "/tmp/azure-token")

		if got := getOIDCTokenFilePathFromEnv(false); got != "" {
			t.Fatalf("expected no token file path when AKS workload identity is disabled, got %q", got)
		}

		if got := getOIDCTokenFilePathFromEnv(true); got != "/tmp/azure-token" {
			t.Fatalf("expected AZURE_FEDERATED_TOKEN_FILE to be used when AKS workload identity is enabled, got %q", got)
		}
	})
}
