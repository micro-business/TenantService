package service_test

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microbusinesslimited/Micro-Business-Core/system"
	"github.com/microbusinesslimited/TenantService/business/domain"
	"github.com/microbusinesslimited/TenantService/business/service"
	"github.com/microbusinesslimited/TenantService/data/contract"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadAllApplications method input parameters and dependency test", func() {
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

			Ω(func() { tenantService.ReadAllApplications(validTenantID) }).Should(Panic())
		})
	})

	Describe("Input Parameters", func() {
		It("should panic when empty tenant unique identifier provided", func() {
			Ω(func() { tenantService.ReadAllApplications(system.EmptyUUID) }).Should(Panic())
		})
	})
})

var _ = Describe("ReadAllApplications method behaviour", func() {
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

	It("should call tenant data service ReadAllApplications function", func() {
		mockTenantDataService.EXPECT().ReadAllApplications(validTenantID)

		tenantService.ReadAllApplications(validTenantID)
	})

	Context("when tenant data service succeeds to read the requested applications for tenant without any application registered", func() {
		It("should return no error and empty list of applications", func() {
			expectedDomainApplications := make(map[system.UUID]domain.Application)
			expectedApplications := make(map[system.UUID]contract.Application)

			mockTenantDataService.
				EXPECT().
				ReadAllApplications(validTenantID).
				Return(expectedApplications, nil)

			applications, err := tenantService.ReadAllApplications(validTenantID)

			Expect(applications).To(Equal(expectedDomainApplications))
			Expect(err).To(BeNil())
		})
	})

	Context("when tenant data service succeeds to read the requested applications", func() {
		It("should return no error", func() {
			expectedDomainApplications := make(map[system.UUID]domain.Application)
			expectedApplications := make(map[system.UUID]contract.Application)

			for idx := 0; idx < rand.Intn(10)+1; idx++ {
				applicationID, _ := system.RandomUUID()
				randomValue, _ := system.RandomUUID()

				expectedDomainApplications[applicationID] = domain.Application{Name: randomValue.String()}
				expectedApplications[applicationID] = contract.Application{Name: randomValue.String()}
			}

			mockTenantDataService.
				EXPECT().
				ReadAllApplications(validTenantID).
				Return(expectedApplications, nil)

			applications, err := tenantService.ReadAllApplications(validTenantID)

			Expect(applications).To(Equal(expectedDomainApplications))
			Expect(err).To(BeNil())
		})
	})

	Context("when tenant data service fails to read the requested applications", func() {
		It("should return the error returned by tenant data service", func() {
			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockTenantDataService.
				EXPECT().
				ReadAllApplications(validTenantID).
				Return(nil, expectedError)

			applications, err := tenantService.ReadAllApplications(validTenantID)

			Eventually(applications).Should(HaveLen(0))
			Expect(err).To(Equal(expectedError))
		})
	})
})

func TestReadAllApplications(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReadAllApplications method input parameters and dependency test")
	RunSpecs(t, "ReadAllApplications method behaviour")
}
