package service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/micro-business/Micro-Business-Core/system"
	"github.com/micro-business/TenantService/business/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteApplication method input parameters and dependency test", func() {
	var (
		mockCtrl              *gomock.Controller
		tenantService         *service.TenantService
		mockTenantDataService *MockTenantDataService
		validTenantID         system.UUID
		validApplicationID    system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		validTenantID, _ = system.RandomUUID()
		validApplicationID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when tenant data service not provided", func() {
		It("should panic", func() {
			tenantService.TenantDataService = nil

			Ω(func() { tenantService.DeleteApplication(validTenantID, validApplicationID) }).Should(Panic())
		})
	})

	Describe("Input Parameters", func() {
		It("should panic when empty tenant unique identifier provided", func() {
			Ω(func() { tenantService.DeleteApplication(system.EmptyUUID, validApplicationID) }).Should(Panic())
		})

		It("should panic when empty application unique identifier provided", func() {
			Ω(func() { tenantService.DeleteApplication(validTenantID, system.EmptyUUID) }).Should(Panic())
		})
	})
})

var _ = Describe("DeleteApplication method behaviour", func() {
	var (
		mockCtrl              *gomock.Controller
		tenantService         *service.TenantService
		mockTenantDataService *MockTenantDataService
		validTenantID         system.UUID
		validApplicationID    system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		validTenantID, _ = system.RandomUUID()
		validApplicationID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call tenant data service DeleteApplication function", func() {
		mockTenantDataService.EXPECT().DeleteApplication(validTenantID, validApplicationID)

		tenantService.DeleteApplication(validTenantID, validApplicationID)
	})

	Context("when tenant data service succeeds to delete the existing application", func() {
		It("should return no error", func() {
			mockTenantDataService.
				EXPECT().
				DeleteApplication(validTenantID, validApplicationID).
				Return(nil)

			err := tenantService.DeleteApplication(validTenantID, validApplicationID)

			Expect(err).To(BeNil())
		})
	})

	Context("when tenant data service fails to delete the existing application", func() {
		It("should return error returned by tenant data service", func() {
			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockTenantDataService.
				EXPECT().
				DeleteApplication(validTenantID, validApplicationID).
				Return(expectedError)

			err := tenantService.DeleteApplication(validTenantID, validApplicationID)

			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestDeleteApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DeleteApplication method input parameters and dependency test")
	RunSpecs(t, "DeleteApplication method behaviour")
}
