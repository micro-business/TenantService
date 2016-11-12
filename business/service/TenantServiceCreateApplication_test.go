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

var _ = Describe("CreateApplication method input parameters and dependency test", func() {
	var (
		mockCtrl                     *gomock.Controller
		tenantService                *service.TenantService
		mockTenantDataService        *MockTenantDataService
		validTenantID                system.UUID
		validApplication             domain.Application
		tenantWithEmptyName          domain.Application
		tenantWithWhitespaceOnlyName domain.Application
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		validTenantID, _ = system.RandomUUID()
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

			Ω(func() { tenantService.CreateApplication(validTenantID, validApplication) }).Should(Panic())
		})
	})

	Describe("Input Parameters", func() {
		It("should panic when tenant with empty name provided", func() {
			Ω(func() { tenantService.CreateApplication(validTenantID, tenantWithEmptyName) }).Should(Panic())
		})

		It("should panic when tenant with name contains whitespace characters only provided", func() {
			Ω(func() { tenantService.CreateApplication(validTenantID, tenantWithWhitespaceOnlyName) }).Should(Panic())
		})
	})
})

var _ = Describe("CreateApplication method behaviour", func() {
	var (
		mockCtrl              *gomock.Controller
		tenantService         *service.TenantService
		mockTenantDataService *MockTenantDataService
		validTenantID         system.UUID
		validApplication      domain.Application
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantDataService = NewMockTenantDataService(mockCtrl)

		tenantService = &service.TenantService{TenantDataService: mockTenantDataService}

		validTenantID, _ = system.RandomUUID()
		validApplication = domain.Application{Name: "Test Name"}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call tenant data service CreateApplication function", func() {
		mappedApplication := contract.Application{Name: validApplication.Name}

		mockTenantDataService.EXPECT().CreateApplication(validTenantID, mappedApplication)

		tenantService.CreateApplication(validTenantID, validApplication)
	})

	Context("when tenant data service succeeds to create the new application", func() {
		It("should return the returned application unique identifier by tenant data service and no error", func() {
			key, _ := system.RandomUUID()
			mappedApplication := contract.Application{Name: key.String()}

			expectedTenantID, _ := system.RandomUUID()
			mockTenantDataService.
				EXPECT().
				CreateApplication(validTenantID, mappedApplication).
				Return(expectedTenantID, nil)

			newApplicationID, err := tenantService.CreateApplication(validTenantID, domain.Application{Name: key.String()})

			Expect(expectedTenantID).To(Equal(newApplicationID))
			Expect(err).To(BeNil())
		})
	})

	Context("when tenant data service fails to create the new application", func() {
		It("should return application unique identifier as empty UUID and the returned error by tenant data service", func() {
			mappedApplication := contract.Application{Name: validApplication.Name}

			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockTenantDataService.
				EXPECT().
				CreateApplication(validTenantID, mappedApplication).
				Return(system.EmptyUUID, expectedError)

			newApplicationID, err := tenantService.CreateApplication(validTenantID, validApplication)

			Expect(newApplicationID).To(Equal(system.EmptyUUID))
			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestCreateApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateApplication method input parameters and dependency test")
	RunSpecs(t, "CreateApplication method behaviour")
}
