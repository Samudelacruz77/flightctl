package proxy_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/flightctl/flightctl/test/harness/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	proxyDropInFile = "proxy.conf"
	proxyTimeout    = "5m"
	proxyPolling    = "250ms"
)

var _ = Describe("Agent HTTP proxy support", func() {
	It("routes agent traffic through the configured HTTPS_PROXY", Label("EDM-4246", "agent", "proxy"), func() {
		harness := e2e.GetWorkerHarness()

		By("discovering the API server endpoint from the VM agent config")
		apiHost, apiIP, apiPort, err := harness.GetAPIEndpointHostIPPortFromVM()
		Expect(err).ToNot(HaveOccurred(), "failed to read API endpoint from agent config")
		Expect(apiIP).ToNot(BeEmpty(), "API server IP should not be empty")
		Expect(apiPort).ToNot(BeEmpty(), "API server port should not be empty")
		GinkgoWriter.Printf("API endpoint: host=%s ip=%s port=%s\n", apiHost, apiIP, apiPort)

		By("discovering the VM default gateway to reach host services")
		gatewayIP, err := getVMDefaultGateway(harness)
		Expect(err).ToNot(HaveOccurred(), "failed to discover VM default gateway")
		Expect(gatewayIP).ToNot(BeEmpty(), "VM default gateway should not be empty")
		GinkgoWriter.Printf("VM default gateway: %s\n", gatewayIP)

		By("configuring HTTPS_PROXY via systemd drop-in")
		proxyURL := buildProxyURL(gatewayIP, squidSvc.Port)
		dropInContent := buildProxyDropIn(proxyURL)
		err = harness.CreateAgentDropIn(proxyDropInFile, dropInContent)
		Expect(err).ToNot(HaveOccurred(), "failed to create proxy drop-in")

		err = harness.VMDaemonReload()
		Expect(err).ToNot(HaveOccurred(), "failed to daemon-reload after drop-in")

		By("blocking direct traffic from VM to the API server")
		harness.BlockTrafficOnVM(apiIP, apiPort)
		DeferCleanup(func() {
			if harness.IsTrafficBlockedOnVM(apiIP, apiPort) {
				_ = harness.UnblockTrafficOnVM(apiIP, apiPort)
			}
		})

		By("starting the flightctl-agent with proxy configuration")
		err = harness.StartFlightCtlAgent()
		Expect(err).ToNot(HaveOccurred(), "failed to start flightctl-agent")

		By("enrolling the device through the proxy and waiting for online status")
		deviceID, _ := harness.EnrollAndWaitForOnlineStatus()
		Expect(deviceID).ToNot(BeEmpty(), "device should have enrolled through the proxy")
		GinkgoWriter.Printf("Device %s enrolled and online via proxy\n", deviceID)

		By("verifying Squid access logs contain CONNECT entries for the API server")
		connectTarget := buildConnectTarget(apiHost, apiIP, apiPort)
		Eventually(func() error {
			return verifySquidLogs(connectTarget)
		}, 30*time.Second, time.Second).Should(Succeed(),
			"squid access logs should contain CONNECT entry for %s", connectTarget)

		By("verifying VM network connections include the proxy endpoint")
		ssCmd := fmt.Sprintf("ss -tnp | grep %s:%s || true", gatewayIP, squidSvc.Port)
		Eventually(harness.VMCommandOutputFunc(ssCmd, false), proxyTimeout, proxyPolling).
			ShouldNot(BeEmpty(), "VM should have TCP connections to the proxy at %s:%s", gatewayIP, squidSvc.Port)
	})
})

func getVMDefaultGateway(h *e2e.Harness) (string, error) {
	if h == nil || h.VM == nil {
		return "", fmt.Errorf("harness or VM is nil")
	}
	out, err := h.VM.RunSSH([]string{"ip", "route", "show", "default"}, nil)
	if err != nil {
		return "", fmt.Errorf("failed to read default route: %w", err)
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) < 3 || fields[0] != "default" || fields[1] != "via" {
		return "", fmt.Errorf("unexpected default route output: %q", out.String())
	}
	return fields[2], nil
}

func buildProxyURL(gatewayIP, port string) string {
	return fmt.Sprintf("http://%s:%s", gatewayIP, port)
}

func buildProxyDropIn(proxyURL string) string {
	return fmt.Sprintf("[Service]\nEnvironment=\"HTTPS_PROXY=%s\"\nEnvironment=\"HTTP_PROXY=%s\"\n", proxyURL, proxyURL)
}

// buildConnectTarget returns the host:port string to look for in Squid CONNECT logs.
// HTTP CONNECT uses the hostname from the request URL, not the resolved IP.
func buildConnectTarget(apiHost, apiIP, apiPort string) string {
	host := apiHost
	if host == "" {
		host = apiIP
	}
	return fmt.Sprintf("%s:%s", host, apiPort)
}

func verifySquidLogs(connectTarget string) error {
	logs, err := squidSvc.AccessLogs()
	if err != nil {
		return fmt.Errorf("failed to read squid access logs: %w", err)
	}
	if !containsConnectEntry(logs, connectTarget) {
		return fmt.Errorf("squid access logs do not contain CONNECT entry for %s; logs:\n%s", connectTarget, logs)
	}
	return nil
}

func containsConnectEntry(logs, target string) bool {
	for _, line := range strings.Split(logs, "\n") {
		if strings.Contains(line, "CONNECT") && strings.Contains(line, target) {
			return true
		}
	}
	return false
}
