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

var _ = Describe("UpdateApplication method input parameters and dependency test", func() {
	var (
		mockCtrl                     *gomock.Controller
		tenantService                *service.TenantService
		mockTenantDataService        *MockTenantDataService
		validTenantID                system.UUID
		validApplicationID           system.UUID
		validApplication             domain.Application
		tenantWithEmptyName          domain.Application
		tenantWithWhitespaceOnlyName domain.Application
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		validTenantID, _ = system.RandomUUID()
		validApplicationID, _ = system.RandomUUID()
		validApplication = domain.Application{Name: "Test Name"}
		tenantWithEmptyName = domain.Application{Name: ""}
		tenantWithWhitespaceOnlyName = domain.Application{Name: "   "}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when tenant data service not provided", func() {
		It("should panic", func() {
			tenantService.TenantDataService = nil

			Ω(func() { tenantService.UpdateApplication(validTenantID, validApplicationID, validApplication) }).Should(Panic())
		})
	})

	Describe("Input Parameters", func() {
		It("should panic when empty tenant unique identifier provided", func() {
			Ω(func() { tenantService.UpdateApplication(system.EmptyUUID, validApplicationID, validApplication) }).Should(Panic())
		})

		It("should panic when empty application unique identifier provided", func() {
			Ω(func() { tenantService.UpdateApplication(validTenantID, system.EmptyUUID, validApplication) }).Should(Panic())
		})

		It("should panic when tenant with empty name provided", func() {
			Ω(func() { tenantService.UpdateApplication(validTenantID, validApplicationID, tenantWithEmptyName) }).Should(Panic())
		})

		It("should panic when tenant with name contains whitespace characters only provided", func() {
			Ω(func() {
				tenantService.UpdateApplication(validTenantID, validApplicationID, tenantWithWhitespaceOnlyName)
			}).Should(Panic())
		})
	})
})

var _ = Describe("UpdateApplication method behaviour", func() {
	var (
		mockCtrl              *gomock.Controller
		tenantService         *service.TenantService
		mockTenantDataService *MockTenantDataService
		validTenantID         system.UUID
		validApplicationID    system.UUID
		validApplication      domain.Application
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		validTenantID, _ = system.RandomUUID()
		validApplicationID, _ = system.RandomUUID()
		validApplication = domain.Application{Name: "Test Name"}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call tenant data service UpdateApplication function", func() {
		mappedApplication := contract.Application{Name: validApplication.Name}

		mockTenantDataService.EXPECT().UpdateApplication(validTenantID, validApplicationID, mappedApplication)

		tenantService.UpdateApplication(validTenantID, validApplicationID, validApplication)
	})

	Context("when tenant data service succeeds to update the existing application", func() {
		It("should return no error", func() {
			mappedApplication := contract.Application{Name: validApplication.Name}

			mockTenantDataService.
				EXPECT().
				UpdateApplication(validTenantID, validApplicationID, mappedApplication).
				Return(nil)

			err := tenantService.UpdateApplication(validTenantID, validApplicationID, validApplication)

			Expect(err).To(BeNil())
		})
	})

	Context("when tenant data service fails to update the existing application", func() {
		It("should return error returned by tenant data service", func() {
			mappedApplication := contract.Application{Name: validApplication.Name}

			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockTenantDataService.
				EXPECT().
				UpdateApplication(validTenantID, validApplicationID, mappedApplication).
				Return(expectedError)

			err := tenantService.UpdateApplication(validTenantID, validApplicationID, validApplication)

			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestUpdateApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateApplication method input parameters and dependency test")
	RunSpecs(t, "UpdateApplication method behaviour")
}
