package service_test

import (
	"testing"

	"github.com/micro-business/Micro-Business-Core/system"
	"github.com/micro-business/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateApplication method input parameters and dependency test", func() {
	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService := &service.TenantDataService{ClusterConfig: nil}
			validTenantID, _ := system.RandomUUID()
			validApplicationID, _ := system.RandomUUID()

			Ω(func() {
				tenantDataService.UpdateApplication(validTenantID, validApplicationID, createApplicationInfo())
			}).Should(Panic())
		})
	})
})

func TestUpdateApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateApplication method input parameters and dependency test")
}
