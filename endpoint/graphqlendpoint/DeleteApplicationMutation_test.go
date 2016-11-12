package graphqlendpoint_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/graphql-go/graphql"
	"github.com/microbusinesslimited/Micro-Business-Core/system"
	"github.com/microbusinesslimited/TenantService/endpoint/graphqlendpoint"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteApplication method input parameters and dependency test", func() {
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
			query := "mutation {deleteApplication (applicationID: \"" + applicationID.String() + "\")}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if TenantID format is not UUID", func() {
			query := "mutation {deleteApplication (tenantID: \"Invalid UUID\", applicationID: \"" + applicationID.String() + "\")}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if no ApplicationID provided", func() {
			query := "mutation {deleteApplication (tenantID: \"" + tenantID.String() + "\")}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if ApplicationID format is not UUID", func() {
			query := "mutation {deleteApplication (applicationID: \"Invalid UUID\", tenantID: \"" + tenantID.String() + "\")}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
	})
})

var _ = Describe("DeleteApplication method behaviour", func() {
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

	It("should call tenant service DeleteApplication function", func() {
		mockTenantService.EXPECT().DeleteApplication(tenantID, applicationID).Return(nil)

		query := "mutation {deleteApplication (tenantID: \"" + tenantID.String() + "\", applicationID: \"" + applicationID.String() + "\")}"

		graphqlendpoint.ExecuteQuery(query, mockTenantService)
	})

	It("should return error if tenant service DeleteApplication function returns error", func() {
		randomValue, _ := system.RandomUUID()
		mockTenantService.EXPECT().DeleteApplication(tenantID, applicationID).Return(fmt.Errorf(randomValue.String()))

		query := "mutation {deleteApplication (tenantID: \"" + tenantID.String() + "\", applicationID: \"" + applicationID.String() + "\")}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(Equal(fmt.Errorf(randomValue.String())))
		Expect(result).To(BeNil())
	})

	It("should return true if tenant service DeleteApplication function returns no error", func() {
		mockTenantService.EXPECT().DeleteApplication(tenantID, applicationID).Return(nil)

		expectedApplication := &graphql.Result{
			Data: map[string]interface{}{
				"deleteApplication": true,
			},
		}

		query := "mutation {deleteApplication (tenantID: \"" + tenantID.String() + "\", applicationID: \"" + applicationID.String() + "\")}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(expectedApplication))
	})
})

func TestDeleteApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DeleteApplication method input parameters and dependency test")
	RunSpecs(t, "DeleteApplication method behaviour")
}
