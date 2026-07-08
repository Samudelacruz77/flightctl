package proxy_test

import (
	"context"
	"testing"

	"github.com/flightctl/flightctl/test/e2e/infra/auxiliary"
	"github.com/flightctl/flightctl/test/e2e/infra/setup"
	"github.com/flightctl/flightctl/test/harness/e2e"
	"github.com/flightctl/flightctl/test/login"
	testutil "github.com/flightctl/flightctl/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProxy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Proxy E2E Suite")
}

var squidSvc *auxiliary.Squid

var _ = BeforeSuite(func() {
	ctx := context.Background()
	auxiliary.Get(ctx)

	squidServices, err := auxiliary.StartServices(ctx, []auxiliary.Service{auxiliary.ServiceSquid})
	Expect(err).ToNot(HaveOccurred())
	squidSvc = squidServices.Squid
	Expect(squidSvc).ToNot(BeNil())
	GinkgoWriter.Printf("Squid proxy started: %s (port %s)\n", squidSvc.URL, squidSvc.Port)

	Expect(setup.EnsureDefaultProviders(nil)).To(Succeed())
	e2e.SetupWorkerHarnessOrAbort()
})

var _ = BeforeEach(func() {
	workerID := GinkgoParallelProcess()
	harness := e2e.GetWorkerHarness()
	suiteCtx := e2e.GetWorkerContext()

	_, err := login.LoginToAPIWithToken(harness)
	Expect(err).ToNot(HaveOccurred())

	GinkgoWriter.Printf("[BeforeEach] Worker %d: Setting up test with VM from pool (proxy)\n", workerID)

	ctx := testutil.StartSpecTracerForGinkgo(suiteCtx)
	harness.SetTestContext(ctx)

	err = harness.SetupVMFromPool(workerID)
	Expect(err).ToNot(HaveOccurred())

	GinkgoWriter.Printf("[BeforeEach] Worker %d: VM ready, agent NOT started (proxy config pending)\n", workerID)
})

var _ = AfterEach(func() {
	workerID := GinkgoParallelProcess()
	GinkgoWriter.Printf("[AfterEach] Worker %d: Cleaning up test resources\n", workerID)

	harness := e2e.GetWorkerHarness()
	suiteCtx := e2e.GetWorkerContext()

	harness.PrintAgentLogsIfFailed()
	harness.CaptureDeploymentLogsIfFailed()

	_, err := login.LoginToAPIWithToken(harness)
	Expect(err).ToNot(HaveOccurred())

	err = harness.CleanUpAllTestResources()
	Expect(err).ToNot(HaveOccurred())

	harness.SetTestContext(suiteCtx)
	GinkgoWriter.Printf("[AfterEach] Worker %d: Test cleanup completed\n", workerID)
})
