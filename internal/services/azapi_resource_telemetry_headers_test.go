package services_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGenericResource_telemetryHeadersNoPlanDiff(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.resourceGroupWithTelemetryHeaders(data, "", ""),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:             r.resourceGroupWithTelemetryHeaders(data, "module-added", "1.0.0"),
			PlanOnly:           true,
			ExpectNonEmptyPlan: false,
		},
		{
			Config:             r.resourceGroupWithTelemetryHeaders(data, "module-changed", "2.0.0"),
			PlanOnly:           true,
			ExpectNonEmptyPlan: false,
		},
		{
			Config:             r.resourceGroupWithTelemetryHeaders(data, "", ""),
			PlanOnly:           true,
			ExpectNonEmptyPlan: false,
		},
	})
}

func TestAccGenericResource_telemetryHeadersSentOnCreate(t *testing.T) {
	mitmLogPath := os.Getenv("AZAPI_ACCTEST_MITMPROXY_LOG")
	if mitmLogPath == "" {
		mitmLogPath = "/tmp/azapi-acctest/mitmproxy.log"
	}
	if _, err := exec.LookPath("mitmdump"); err != nil {
		t.Skip("mitmdump is required to verify telemetry headers in captured HTTP requests")
	}
	if _, err := os.Stat(mitmLogPath); err != nil {
		t.Skipf("mitmproxy capture file not found: %s", mitmLogPath)
	}

	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	moduleID := "module-create-telemetry"
	moduleVersion := "1.2.3"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.resourceGroupWithTelemetryHeaders(data, moduleID, moduleVersion),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				checkTelemetryHeadersCapturedInMitmLog(mitmLogPath, moduleID, moduleVersion),
			),
		},
	})
}

func checkTelemetryHeadersCapturedInMitmLog(mitmLogPath, moduleID, moduleVersion string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		out, err := exec.Command("mitmdump", "-nr", mitmLogPath, "--set", "flow_detail=4").CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to parse mitmproxy log: %w", err)
		}
		content := strings.ToLower(string(out))
		expectedModuleID := strings.ToLower("x-module-id: " + moduleID)
		expectedModuleVersion := strings.ToLower("x-module-version: " + moduleVersion)
		if !strings.Contains(content, expectedModuleID) {
			return fmt.Errorf("expected telemetry module id header %q in captured HTTP requests", expectedModuleID)
		}
		if !strings.Contains(content, expectedModuleVersion) {
			return fmt.Errorf("expected telemetry module version header %q in captured HTTP requests", expectedModuleVersion)
		}
		return nil
	}
}

func (r GenericResource) resourceGroupWithTelemetryHeaders(data acceptance.TestData, moduleID, moduleVersion string) string {
	telemetryHeaders := ""
	if moduleID != "" || moduleVersion != "" {
		telemetryHeaders = fmt.Sprintf(`
  telemetry_headers = {
    module_id      = %q
    module_version = %q
  }`, moduleID, moduleVersion)
	}

	return fmt.Sprintf(`
resource "azapi_resource" "test" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
%[3]s
}`, data.RandomInteger, data.LocationPrimary, telemetryHeaders)
}
