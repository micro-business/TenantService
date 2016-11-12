package service_test

import (
	"testing"

	"github.com/microbusinesslimited/Micro-Business-Core/system"
	"github.com/microbusinesslimited/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateTenant method input parameters and dependency test", func() {
	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService := &service.TenantDataService{ClusterConfig: nil}
			tenantID, _ := system.RandomUUID()

			Ω(func() { tenantDataService.UpdateTenant(tenantID, createTenantInfo()) }).Should(Panic())
		})
	})
})

func TestUpdateTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateTenant method input parameters and dependency test")
}
