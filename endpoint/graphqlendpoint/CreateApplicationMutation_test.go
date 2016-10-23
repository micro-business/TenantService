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

var _ = Describe("CreateApplication method input parameters and dependency test", func() {
	var (
		mockCtrl          *gomock.Controller
		mockTenantService *MockTenantService
		tenantID          system.UUID
		application       domain.Application
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantService = NewMockTenantService(mockCtrl)

		tenantID, _ = system.RandomUUID()

		randomValue, _ := system.RandomUUID()
		application = domain.Application{Name: randomValue.String()}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Input Parameters", func() {
		It("should return error if tenantID not provided", func() {
			query := "mutation {createApplication (application: {Name:\"" + application.Name + "\"})}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if tenantID format is invalid", func() {
			query := "mutation {createApplication (tenantID: \"Invalid UUID\", application: {Name:\"" + application.Name + "\"})}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if no Application provided", func() {
			query := "mutation {createApplication(tenantID: \"" + tenantID.String() + "\")}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
	})
})

var _ = Describe("CreateApplication method behaviour", func() {
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

	It("should call tenant service CreateApplication function", func() {
		mockTenantService.EXPECT().CreateApplication(tenantID, application).Return(applicationID, nil)

		query := "mutation {createApplication (tenantID: \"" + tenantID.String() + "\", application: {Name:\"" + application.Name + "\"})}"

		graphqlendpoint.ExecuteQuery(query, mockTenantService)
	})

	It("should return error if tenant service CreateApplication function returns error", func() {
		randomValue, _ := system.RandomUUID()
		mockTenantService.EXPECT().CreateApplication(tenantID, application).Return(system.EmptyUUID, fmt.Errorf(randomValue.String()))

		query := "mutation {createApplication (tenantID: \"" + tenantID.String() + "\", application: {Name:\"" + application.Name + "\"})}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(Equal(fmt.Errorf(randomValue.String())))
		Expect(result).To(BeNil())
	})

	It("should return application unique identifier if tenant service CreateApplication function returns no error", func() {
		mockTenantService.EXPECT().CreateApplication(tenantID, application).Return(applicationID, nil)

		expectedApplication := &graphql.Result{
			Data: map[string]interface{}{
				"createApplication": applicationID.String(),
			},
		}

		query := "mutation {createApplication (tenantID: \"" + tenantID.String() + "\", application: {Name:\"" + application.Name + "\"})}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(expectedApplication))
	})
})

func TestCreateApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateApplication method input parameters and dependency test")
	RunSpecs(t, "CreateApplication method behaviour")
}
