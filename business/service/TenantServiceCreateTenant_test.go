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

var _ = Describe("CreateTenant method input parameters and dependency test", func() {
	var (
		mockCtrl                          *gomock.Controller
		tenantService                     *service.TenantService
		mockTenantDataService             *MockTenantDataService
		validTenant                       domain.Tenant
		tenantWithEmptySecretKey          domain.Tenant
		tenantWithWhitespaceOnlySecretKey domain.Tenant
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		validTenant = domain.Tenant{SecretKey: "Secret Key"}
		tenantWithEmptySecretKey = domain.Tenant{SecretKey: ""}
		tenantWithWhitespaceOnlySecretKey = domain.Tenant{SecretKey: "   "}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when tenant data service not provided", func() {
		It("should panic", func() {
			tenantService.TenantDataService = nil

			Ω(func() { tenantService.CreateTenant(validTenant) }).Should(Panic())
		})
	})

	Describe("Input Parameters", func() {
		It("should panic when tenant with empty secret key provided", func() {
			Ω(func() { tenantService.CreateTenant(tenantWithEmptySecretKey) }).Should(Panic())
		})

		It("should panic when tenant with secret key contains whitespace characters only provided", func() {
			Ω(func() { tenantService.CreateTenant(tenantWithWhitespaceOnlySecretKey) }).Should(Panic())
		})
	})
})

var _ = Describe("CreateTenant method behaviour", func() {
	var (
		mockCtrl              *gomock.Controller
		tenantService         *service.TenantService
		mockTenantDataService *MockTenantDataService
		validTenant           domain.Tenant
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		validTenant = domain.Tenant{SecretKey: "Secret Key"}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call tenant data service CreateTenant function", func() {
		mappedTenant := contract.Tenant{SecretKey: validTenant.SecretKey}

		mockTenantDataService.EXPECT().CreateTenant(mappedTenant)

		tenantService.CreateTenant(validTenant)
	})

	Context("when tenant data service succeeds to create the new tenant", func() {
		It("should return the returned tenant unique identifier by tenant data service and no error", func() {
			key, _ := system.RandomUUID()
			mappedTenant := contract.Tenant{SecretKey: key.String()}

			expectedTenantID, _ := system.RandomUUID()
			mockTenantDataService.
				EXPECT().
				CreateTenant(mappedTenant).
				Return(expectedTenantID, nil)

			newTenantID, err := tenantService.CreateTenant(domain.Tenant{SecretKey: key.String()})

			Expect(expectedTenantID).To(Equal(newTenantID))
			Expect(err).To(BeNil())
		})
	})

	Context("when tenant data service fails to create the new tenant", func() {
		It("should return tenant unique identifier as empty UUID and the returned error by tenant data service", func() {
			mappedTenant := contract.Tenant{SecretKey: validTenant.SecretKey}

			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockTenantDataService.
				EXPECT().
				CreateTenant(mappedTenant).
				Return(system.EmptyUUID, expectedError)

			newTenantID, err := tenantService.CreateTenant(validTenant)

			Expect(newTenantID).To(Equal(system.EmptyUUID))
			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestCreateTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateTenant method input parameters and dependency test")
	RunSpecs(t, "CreateTenant method behaviour")
}
