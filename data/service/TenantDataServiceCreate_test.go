package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/contract"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create method input parameters and dependency test", func() {
	var (
		mockCtrl                                 *gomock.Controller
		tenantDataService                        *service.TenantDataService
		mockUUIDGeneratorService                 *MockUUIDGeneratorService
		validTenant                              contract.Tenant
		tenantWithEmptySecretKey                 contract.Tenant
		tenantWithSecretKeyContainWhitespaceOnly contract.Tenant
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = NewMockUUIDGeneratorService(mockCtrl)
		tenantDataService = &service.TenantDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: &gocql.ClusterConfig{}}

		randomValue, _ := system.RandomUUID()
		validTenant = contract.Tenant{SecretKey: randomValue.String()}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when UUID generator service not provided", func() {
		It("should panic", func() {
			tenantDataService.UUIDGeneratorService = nil

			Ω(func() { tenantDataService.Create(validTenant) }).Should(Panic())
		})
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService.ClusterConfig = nil

			Ω(func() { tenantDataService.Create(validTenant) }).Should(Panic())
		})
	})

	DescribeTable("Input Parameters",
		func(tenant contract.Tenant) {
			Ω(func() { tenantDataService.Create(tenant) }).Should(Panic())
		},
		Entry("should panic when tenant with empty secret key provided", tenantWithEmptySecretKey),
		Entry("should panic when tenant with secret key contains whitespace only provided", tenantWithSecretKeyContainWhitespaceOnly))
})

func TestCreate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create method input parameters and dependency test")
}
