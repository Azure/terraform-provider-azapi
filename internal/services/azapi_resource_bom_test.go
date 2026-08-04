//go:build unix

package services_test

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGenericResource_utf8BOMResponse(t *testing.T) {
	startPeeringBOMMockProxy(t)

	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.peeringWithBOMResponse(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func (r GenericResource) peeringWithBOMResponse(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azapi" {
  # The Microsoft.Peering endpoints are fully mocked (see
  # testdata/peering_bom_mock.py), so skip the real resource-provider
  # registration call, which would otherwise reach live Azure.
  skip_provider_registration = true
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Peering/peerings@2025-05-01"
  name      = "acctest-peering-%[1]d"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-%[1]d"
  location  = "eastus"

  schema_validation_enabled = false

  body = {
    kind = "Exchange"
    sku = {
      name = "Basic_Exchange_Free"
    }
    properties = {
      peeringLocation = "peeringLocation0"
      exchange = {
        peerAsn = {
          id = "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Peering/peerAsns/myAsn1"
        }
        connections = [
          {
            connectionIdentifier = "CE495334-0E94-4E51-8164-8116D6CD284D"
            peeringDBFacilityId  = 99999
            bgpSession = {
              maxPrefixesAdvertisedV4 = 1000
              maxPrefixesAdvertisedV6 = 100
              peerSessionIPv4Address  = "192.168.2.1"
              peerSessionIPv6Address  = "fd00::1"
            }
          }
        ]
      }
    }
  }
}
`, data.RandomInteger)
}

func startPeeringBOMMockProxy(t *testing.T) {
	t.Helper()

	if os.Getenv("TF_ACC") == "" {
		return
	}

	mitm, err := exec.LookPath("mitmdump")
	if err != nil {
		t.Skip("mitmdump not found in PATH; it is required to mock the Microsoft.Peering/peerings API (see the azapi-acctest skill)")
	}

	host, port := "127.0.0.1", "7070"
	if proxy := firstNonEmptyEnv("HTTPS_PROXY", "HTTP_PROXY"); proxy != "" {
		if h, p := proxyHostPort(proxy); p != "" {
			host, port = h, p
		}
	} else {
		proxy := fmt.Sprintf("http://%s:%s", host, port)
		if err := os.Setenv("HTTP_PROXY", proxy); err != nil {
			t.Fatalf("setting HTTP_PROXY: %v", err)
		}
		if err := os.Setenv("HTTPS_PROXY", proxy); err != nil {
			t.Fatalf("setting HTTPS_PROXY: %v", err)
		}
	}

	addon, err := filepath.Abs(filepath.Join("testdata", "peering_bom_mock.py"))
	if err != nil {
		t.Fatalf("resolving mock addon path: %v", err)
	}

	addr := net.JoinHostPort(host, port)
	if c, derr := net.DialTimeout("tcp", addr, 3000*time.Millisecond); derr == nil {
		_ = c.Close()
		t.Skipf("%s is already in use; stop any other proxy on that port before running this reproduction test", addr)
	}

	logFile, err := os.CreateTemp("", "peering-bom-mitmdump-*.log")
	if err != nil {
		t.Fatalf("creating mitmdump log file: %v", err)
	}
	logPath := logFile.Name()
	readLog := func() string {
		b, _ := os.ReadFile(logPath)
		return string(b)
	}

	cmd := exec.Command(mitm, "--listen-host", host, "--listen-port", port, "-s", addon, "-q")
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		_ = logFile.Close()
		_ = os.Remove(logPath)
		t.Fatalf("starting mitmdump: %v", err)
	}
	if err := logFile.Close(); err != nil {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		_ = cmd.Wait()
		_ = os.Remove(logPath)
		t.Fatalf("closing mitmdump log file: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	stop := func() {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-done
		_ = os.Remove(logPath)
	}

	deadline := time.Now().Add(30 * time.Second)
	for {
		select {
		case werr := <-done:
			out := readLog()
			_ = os.Remove(logPath)
			t.Fatalf("mitmdump exited before becoming ready: %v\n%s", werr, out)
		default:
		}
		if c, derr := net.DialTimeout("tcp", addr, time.Second); derr == nil {
			_ = c.Close()
			t.Cleanup(stop)
			return
		}
		if time.Now().After(deadline) {
			stop()
			t.Fatalf("timed out waiting for mitmdump to listen on %s\n%s", addr, readLog())
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func firstNonEmptyEnv(keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func proxyHostPort(proxy string) (string, string) {
	if !strings.Contains(proxy, "://") {
		proxy = "http://" + proxy
	}
	u, err := url.Parse(proxy)
	if err != nil {
		return "", ""
	}
	host := u.Hostname()
	if host == "" {
		host = "127.0.0.1"
	}
	return host, u.Port()
}
