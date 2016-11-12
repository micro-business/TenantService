package service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/micro-business/Micro-Business-Core/system"
	"github.com/micro-business/TenantService/business/domain"
	"github.com/micro-business/TenantService/business/service"
	"github.com/micro-business/TenantService/data/contract"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateTenant method input parameters and dependency test", func() {
	var (
		mockCtrl                          *gomock.Controller
		tenantService                     *service.TenantService
		mockTenantDataService             *MockTenantDataService
		validTenantID                     system.UUID
		validTenant                       domain.Tenant
		tenantWithEmptySecretKey          domain.Tenant
		tenantWithWhitespaceOnlySecretKey domain.Tenant
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		validTenantID, _ = system.RandomUUID()
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

			Ω(func() { tenantService.UpdateTenant(validTenantID, validTenant) }).Should(Panic())
		})
	})

	Describe("Input Parameters", func() {
		It("should panic when empty tenant unique identifier provided", func() {
			Ω(func() { tenantService.UpdateTenant(system.EmptyUUID, validTenant) }).Should(Panic())
		})

		It("should panic when tenant with empty secret key provided", func() {
			Ω(func() { tenantService.UpdateTenant(validTenantID, tenantWithEmptySecretKey) }).Should(Panic())
		})

		It("should panic when tenant with secret key contains whitespace characters only provided", func() {
			Ω(func() { tenantService.UpdateTenant(validTenantID, tenantWithWhitespaceOnlySecretKey) }).Should(Panic())
		})
	})
})

var _ = Describe("UpdateTenant method behaviour", func() {
	var (
		mockCtrl              *gomock.Controller
		tenantService         *service.TenantService
		mockTenantDataService *MockTenantDataService
		validTenantID         system.UUID
		validTenant           domain.Tenant
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		validTenantID, _ = system.RandomUUID()
		validTenant = domain.Tenant{SecretKey: "Secret Key"}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call tenant data service UpdateTenant function", func() {
		mappedTenant := contract.Tenant{SecretKey: validTenant.SecretKey}

		mockTenantDataService.EXPECT().UpdateTenant(validTenantID, mappedTenant)

		tenantService.UpdateTenant(validTenantID, validTenant)
	})

	Context("when tenant data service succeeds to update the existing tenant", func() {
		It("should return no error", func() {
			mappedTenant := contract.Tenant{SecretKey: validTenant.SecretKey}

			mockTenantDataService.
				EXPECT().
				UpdateTenant(validTenantID, mappedTenant).
				Return(nil)

			err := tenantService.UpdateTenant(validTenantID, validTenant)

			Expect(err).To(BeNil())
		})
	})

	Context("when tenant data service fails to update the existing tenant", func() {
		It("should return error returned by tenant data service", func() {
			mappedTenant := contract.Tenant{SecretKey: validTenant.SecretKey}

			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockTenantDataService.
				EXPECT().
				UpdateTenant(validTenantID, mappedTenant).
				Return(expectedError)

			err := tenantService.UpdateTenant(validTenantID, validTenant)

			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestUpdateTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateTenant method input parameters and dependency test")
	RunSpecs(t, "UpdateTenant method behaviour")
}
