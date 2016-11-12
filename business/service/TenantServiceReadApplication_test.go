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

var _ = Describe("ReadApplication method input parameters and dependency test", func() {
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

			Ω(func() { tenantService.ReadApplication(validTenantID, validApplicationID) }).Should(Panic())
		})
	})

	Describe("Input Parameters", func() {
		It("should panic when empty tenant unique identifier provided", func() {
			Ω(func() { tenantService.ReadApplication(system.EmptyUUID, validApplicationID) }).Should(Panic())
		})

		It("should panic when empty application unique identifier provided", func() {
			Ω(func() { tenantService.ReadApplication(validTenantID, system.EmptyUUID) }).Should(Panic())
		})
	})
})

var _ = Describe("ReadApplication method behaviour", func() {
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

	It("should call tenant data service ReadApplication function", func() {
		mockTenantDataService.EXPECT().ReadApplication(validTenantID, validApplicationID)

		tenantService.ReadApplication(validTenantID, validApplicationID)
	})

	Context("when tenant data service succeeds to read the requested application", func() {
		It("should return no error", func() {
			randomValue, _ := system.RandomUUID()
			expectedApplication := domain.Application{Name: randomValue.String()}

			mockTenantDataService.
				EXPECT().
				ReadApplication(validTenantID, validApplicationID).
				Return(contract.Application{Name: expectedApplication.Name}, nil)

			application, err := tenantService.ReadApplication(validTenantID, validApplicationID)

			Expect(application).To(Equal(expectedApplication))
			Expect(err).To(BeNil())
		})
	})

	Context("when tenant data service fails to read the requested application", func() {
		It("should return the error returned by tenant data service", func() {
			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockTenantDataService.
				EXPECT().
				ReadApplication(validTenantID, validApplicationID).
				Return(contract.Application{}, expectedError)

			application, err := tenantService.ReadApplication(validTenantID, validApplicationID)

			Expect(application).To(Equal(domain.Application{}))
			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestReadApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReadApplication method input parameters and dependency test")
	RunSpecs(t, "ReadApplication method behaviour")
}
