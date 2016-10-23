package graphqlendpoint_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/business/domain"
	"github.com/microbusinesses/TenantService/endpoint/graphqlendpoint"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateApplication method input parameters and dependency test", func() {
	var (
		mockCtrl          *gomock.Controller
		mockTenantService *MockTenantService
		tenantID          system.UUID
		applicationID     system.UUID
		application       domain.Application
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantService = NewMockTenantService(mockCtrl)

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()

		randomValue, _ := system.RandomUUID()
		application = domain.Application{Name: randomValue.String()}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Input Parameters", func() {
		It("should return error if tenantID not provided", func() {
			query := "mutation {updateApplication (applicationID: \"" + applicationID.String() + "\", application: {Name:\"" + application.Name + "\"})}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if tenantID format is invalid", func() {
			query := "mutation {updateApplication (tenantID: \"Invalid UUID\", applicationID: \"" + applicationID.String() + "\", application: {Name:\"" + application.Name + "\"})}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if applicationID not provided", func() {
			query := "mutation {updateApplication (tenantID: \"" + tenantID.String() + "\", application: {Name:\"" + application.Name + "\"})}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if applicationID format is invalid", func() {
			query := "mutation {updateApplication (applicationID: \"Invalid UUID\", tenantID: \"" + tenantID.String() + "\", application: {Name:\"" + application.Name + "\"})}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if no Tenant provided", func() {
			query := "mutation {updateApplication(tenantID: \"" + tenantID.String() + "\", applicationID: \"" + applicationID.String() + "\")}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
	})
})

var _ = Describe("UpdateApplication method behaviour", func() {
	var (
		mockCtrl          *gomock.Controller
		mockTenantService *MockTenantService
		tenantID          system.UUID
		applicationID     system.UUID
		application       domain.Application
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantService = NewMockTenantService(mockCtrl)

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()

		randomValue, _ := system.RandomUUID()
		application = domain.Application{Name: randomValue.String()}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call tenant service UpdateApplication function", func() {
		mockTenantService.EXPECT().UpdateApplication(tenantID, applicationID, application).Return(nil)

		query := "mutation {updateApplication (tenantID: \"" + tenantID.String() + "\", applicationID: \"" + applicationID.String() + "\", application: {Name:\"" + application.Name + "\"})}"

		graphqlendpoint.ExecuteQuery(query, mockTenantService)
	})

	It("should return error if tenant service UpdateApplication function returns error", func() {
		randomValue, _ := system.RandomUUID()
		mockTenantService.EXPECT().UpdateApplication(tenantID, applicationID, application).Return(fmt.Errorf(randomValue.String()))

		query := "mutation {updateApplication (tenantID: \"" + tenantID.String() + "\", applicationID: \"" + applicationID.String() + "\", application: {Name:\"" + application.Name + "\"})}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(Equal(fmt.Errorf(randomValue.String())))
		Expect(result).To(BeNil())
	})

	It("should return tenant unique identifier if tenant service UpdateApplication function returns no error", func() {
		mockTenantService.EXPECT().UpdateApplication(tenantID, applicationID, application).Return(nil)

		expectedApplication := &graphql.Result{
			Data: map[string]interface{}{
				"updateApplication": true,
			},
		}

		query := "mutation {updateApplication (tenantID: \"" + tenantID.String() + "\", applicationID: \"" + applicationID.String() + "\", application: {Name:\"" + application.Name + "\"})}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(expectedApplication))
	})
})

func TestUpdateApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateApplication method input parameters and dependency test")
	RunSpecs(t, "UpdateApplication method behaviour")
}
