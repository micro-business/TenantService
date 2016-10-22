package graphqlendpoint_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/endpoint/graphqlendpoint"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteTenant method input parameters and dependency test", func() {
	var (
		mockCtrl          *gomock.Controller
		mockTenantService *MockTenantService
		tenantID          system.UUID
		applicationID     system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantService = NewMockTenantService(mockCtrl)

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Input Parameters", func() {
		It("should return error if no TenantID provided", func() {
			query := "{application(ApplicationID:\"" + applicationID.String() + "\"){ID Name}}"

			_, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
		})

		It("should return error if no ApplicationID provided", func() {
			query := "{application(TenantID:\"" + tenantID.String() + "\"){ID Name}}"

			_, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
		})

		It("should return error if no TenantID format is not UUID", func() {
			query := "{application(TenantID:\"invalid UUID\", ApplicationID:\"" + applicationID.String() + "\"){ID Name}}"

			_, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
		})

		It("should return error if no ApplicationID format is not UUID", func() {
			query := "{application(TenantID:\"" + tenantID.String() + "\", ApplicationID:\"invalid UUID\"){ID Name}}"

			_, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
		})
	})
})

func TestApplicationQuery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApplicationQuery method input parameters and dependency test")
	RunSpecs(t, "ApplicationQuery method behaviour")
}
