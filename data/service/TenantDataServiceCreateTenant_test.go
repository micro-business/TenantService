package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesslimited/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateTenant method input parameters and dependency test", func() {
	var (
		mockCtrl                 *gomock.Controller
		tenantDataService        *service.TenantDataService
		mockUUIDGeneratorService *MockUUIDGeneratorService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = NewMockUUIDGeneratorService(mockCtrl)
		tenantDataService = &service.TenantDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: &gocql.ClusterConfig{}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when UUID generator service not provided", func() {
		It("should panic", func() {
			tenantDataService.UUIDGeneratorService = nil

			Ω(func() { tenantDataService.CreateTenant(createTenantInfo()) }).Should(Panic())
		})
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService.ClusterConfig = nil

			Ω(func() { tenantDataService.CreateTenant(createTenantInfo()) }).Should(Panic())
		})
	})
})

func TestCreateTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateTenant method input parameters and dependency test")
}
