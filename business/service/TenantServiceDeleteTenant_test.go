package service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/business/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteTenant method input parameters and dependency test", func() {
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

			Ω(func() { tenantService.DeleteTenant(validTenantID) }).Should(Panic())
		})
	})

	Describe("Input Parameters", func() {
		It("should panic when empty tenant unique identifier provided", func() {
			Ω(func() { tenantService.DeleteTenant(system.EmptyUUID) }).Should(Panic())
		})
	})
})

var _ = Describe("DeleteTenant method behaviour", func() {
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

	It("should call tenant data service DeleteTenant function", func() {
		mockTenantDataService.EXPECT().DeleteTenant(validTenantID)

		tenantService.DeleteTenant(validTenantID)
	})

	Context("when tenant data service succeeds to delete the existing tenant", func() {
		It("should return no error", func() {
			mockTenantDataService.
				EXPECT().
				DeleteTenant(validTenantID).
				Return(nil)

			err := tenantService.DeleteTenant(validTenantID)

			Expect(err).To(BeNil())
		})
	})

	Context("when tenant data service fails to delete the existing tenant", func() {
		It("should return error returned by tenant data service", func() {
			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockTenantDataService.
				EXPECT().
				DeleteTenant(validTenantID).
				Return(expectedError)

			err := tenantService.DeleteTenant(validTenantID)

			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestDeleteTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DeleteTenant method input parameters and dependency test")
	RunSpecs(t, "DeleteTenant method behaviour")
}
