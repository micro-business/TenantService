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

var _ = Describe("UpdateApplication method input parameters and dependency test", func() {
	var (
		mockCtrl           *gomock.Controller
		tenantDataService  *service.TenantDataService
		validTenantID      system.UUID
		validApplicationID system.UUID
		validApplication   contract.Application
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		tenantDataService = &service.TenantDataService{ClusterConfig: &gocql.ClusterConfig{}}

		validTenantID, _ = system.RandomUUID()
		validApplicationID, _ = system.RandomUUID()

		randomValue, _ := system.RandomUUID()
		validApplication = contract.Application{Name: randomValue.String()}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService.ClusterConfig = nil

			Ω(func() { tenantDataService.UpdateApplication(validTenantID, validApplicationID, validApplication) }).Should(Panic())
		})
	})
})

func TestUpdateApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateApplication method input parameters and dependency test")
}
