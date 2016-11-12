package service_test

import (
	"testing"

	"github.com/micro-business/Micro-Business-Core/system"
	"github.com/micro-business/TenantService/data/service"
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
