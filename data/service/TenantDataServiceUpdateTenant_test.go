package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/contract"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateTenant method input parameters and dependency test", func() {
	var (
		mockCtrl          *gomock.Controller
		tenantDataService *service.TenantDataService
		validTenantID     system.UUID
		validTenant       contract.Tenant
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		tenantDataService = &service.TenantDataService{ClusterConfig: &gocql.ClusterConfig{}}
		validTenantID, _ = system.RandomUUID()

		randomValue, _ := system.RandomUUID()
		validTenant = contract.Tenant{SecretKey: randomValue.String()}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService.ClusterConfig = nil

			Ω(func() { tenantDataService.UpdateTenant(validTenantID, validTenant) }).Should(Panic())
		})
	})
})

func TestUpdateTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateTenant method input parameters and dependency test")
}
