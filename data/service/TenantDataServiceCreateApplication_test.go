package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/micro-business/Micro-Business-Core/system"
	"github.com/micro-business/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateApplication method input parameters and dependency test", func() {
	var (
		mockCtrl                 *gomock.Controller
		tenantDataService        *service.TenantDataService
		mockUUIDGeneratorService *MockUUIDGeneratorService
		validTenantID            system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = NewMockUUIDGeneratorService(mockCtrl)
		tenantDataService = &service.TenantDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: &gocql.ClusterConfig{}}

		validTenantID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when UUID generator service not provided", func() {
		It("should panic", func() {
			tenantDataService.UUIDGeneratorService = nil

			Ω(func() { tenantDataService.CreateApplication(validTenantID, createApplicationInfo()) }).Should(Panic())
		})
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService.ClusterConfig = nil

			Ω(func() { tenantDataService.CreateApplication(validTenantID, createApplicationInfo()) }).Should(Panic())
		})
	})
})

func TestCreateApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateApplication method input parameters and dependency test")
}
