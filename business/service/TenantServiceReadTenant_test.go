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
		validTenantID         system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		validTenantID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when tenant data service not provided", func() {
		It("should panic", func() {
			tenantService.TenantDataService = nil

			Ω(func() { tenantService.ReadTenant(validTenantID) }).Should(Panic())
		})
	})

	Describe("Input Parameters", func() {
		It("should panic when empty tenant unique identifier provided", func() {
			Ω(func() { tenantService.ReadTenant(system.EmptyUUID) }).Should(Panic())
		})
	})
})

var _ = Describe("ReadTenant method behaviour", func() {
	var (
		mockCtrl              *gomock.Controller
		tenantService         *service.TenantService
		mockTenantDataService *MockTenantDataService
		validTenantID         system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		validTenantID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call tenant data service ReadTenant function", func() {
		mockTenantDataService.EXPECT().ReadTenant(validTenantID)

		tenantService.ReadTenant(validTenantID)
	})

	Context("when tenant data service succeeds to read the requested tenant", func() {
		It("should return no error", func() {
			randomValue, _ := system.RandomUUID()
			expectedTenant := domain.Tenant{SecretKey: randomValue.String()}

			mockTenantDataService.
				EXPECT().
				ReadTenant(validTenantID).
				Return(contract.Tenant{SecretKey: expectedTenant.SecretKey}, nil)

			tenant, err := tenantService.ReadTenant(validTenantID)

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
				ReadTenant(validTenantID).
				Return(contract.Tenant{}, expectedError)

			expectedTenant, err := tenantService.ReadTenant(validTenantID)

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
