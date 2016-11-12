package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/micro-business/Micro-Business-Core/system"
	"github.com/micro-business/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteTenant method input parameters and dependency test", func() {
	var (
		tenantDataService *service.TenantDataService
		validTenantID     system.UUID
	)

	BeforeEach(func() {
		tenantDataService = &service.TenantDataService{ClusterConfig: &gocql.ClusterConfig{}}

		validTenantID, _ = system.RandomUUID()
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService.ClusterConfig = nil

			Ω(func() { tenantDataService.DeleteTenant(validTenantID) }).Should(Panic())
		})
	})
})

func TestDeleteTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DeleteTenant method input parameters and dependency test")
}
