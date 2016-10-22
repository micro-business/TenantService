package service_test

import (
	"testing"

	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadTenant method input parameters and dependency test", func() {
	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService := &service.TenantDataService{ClusterConfig: nil}
			tenantID, _ := system.RandomUUID()

			Ω(func() { tenantDataService.ReadTenant(tenantID) }).Should(Panic())
		})
	})
})

func TestReadTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReadTenant method input parameters and dependency test")
}
