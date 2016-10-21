package service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/business/domain"
	"github.com/microbusinesses/TenantService/business/service"
	"github.com/microbusinesses/TenantService/data/contract"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadTenant method input parameters and dependency test", func() {
	var (
		mockCtrl              *gomock.Controller
		tenantService         *service.TenantService
		mockTenantDataService *MockTenantDataService
		tenantID              system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		tenantID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when tenant data service not provided", func() {
		It("should panic", func() {
			tenantService.TenantDataService = nil

			Ω(func() { tenantService.Read(tenantID) }).Should(Panic())
		})
	})

	Describe("Input Parameters", func() {
		It("should panic when empty tenant unique identifier provided", func() {
			Ω(func() { tenantService.Read(system.EmptyUUID) }).Should(Panic())
		})
	})
})

var _ = Describe("ReadTenant method behaviour", func() {
	var (
		mockCtrl              *gomock.Controller
		tenantService         *service.TenantService
		mockTenantDataService *MockTenantDataService
		tenantID              system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		tenantID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call tenant data service ReadTenant function", func() {
		mockTenantDataService.EXPECT().ReadTenant(tenantID)

		tenantService.Read(tenantID)
	})

	Context("when tenant data service succeeds to read the requested tenant", func() {
		It("should return no error", func() {
			randomValue, _ := system.RandomUUID()
			expectedTenant := domain.Tenant{SecretKey: randomValue.String()}

			mockTenantDataService.
				EXPECT().
				ReadTenant(tenantID).
				Return(contract.Tenant{SecretKey: expectedTenant.SecretKey}, nil)

			tenant, err := tenantService.Read(tenantID)

			Expect(tenant).To(Equal(expectedTenant))
			Expect(err).To(BeNil())
		})
	})

	Context("when tenant data service fails to read the requested tenant", func() {
		It("should return the error returned by tenant data service", func() {
			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockTenantDataService.
				EXPECT().
				ReadTenant(tenantID).
				Return(contract.Tenant{}, expectedError)

			expectedTenant, err := tenantService.Read(tenantID)

			Expect(expectedTenant).To(Equal(domain.Tenant{}))
			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestReadTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReadTenant method input parameters and dependency test")
	RunSpecs(t, "ReadTenant method behaviour")
}
