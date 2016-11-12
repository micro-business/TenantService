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

var _ = Describe("DeleteTenant method input parameters and dependency test", func() {
	var (
		mockCtrl          *gomock.Controller
		mockTenantService *MockTenantService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantService = NewMockTenantService(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Input Parameters", func() {
		It("should return error if no TenantID provided", func() {
			query := "mutation {deleteTenant}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if TenantID format is not UUID", func() {
			query := "mutation {deleteTenant (tenantID: \"Invalid UUID\")}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
	})
})

var _ = Describe("DeleteTenant method behaviour", func() {
	var (
		mockCtrl          *gomock.Controller
		mockTenantService *MockTenantService
		tenantID          system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantService = NewMockTenantService(mockCtrl)

		tenantID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call tenant service DeleteTenant function", func() {
		mockTenantService.EXPECT().DeleteTenant(tenantID).Return(nil)

		query := "mutation {deleteTenant (tenantID: \"" + tenantID.String() + "\")}"

		graphqlendpoint.ExecuteQuery(query, mockTenantService)
	})

	It("should return error if tenant service DeleteTenant function returns error", func() {
		randomValue, _ := system.RandomUUID()
		mockTenantService.EXPECT().DeleteTenant(tenantID).Return(fmt.Errorf(randomValue.String()))

		query := "mutation {deleteTenant (tenantID: \"" + tenantID.String() + "\")}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(Equal(fmt.Errorf(randomValue.String())))
		Expect(result).To(BeNil())
	})

	It("should return true if tenant service DeleteTenant function returns no error", func() {
		mockTenantService.EXPECT().DeleteTenant(tenantID).Return(nil)

		expectedTenant := &graphql.Result{
			Data: map[string]interface{}{
				"deleteTenant": true,
			},
		}

		query := "mutation {deleteTenant (tenantID: \"" + tenantID.String() + "\")}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(expectedTenant))
	})
})

func TestDeleteTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DeleteTenant method input parameters and dependency test")
	RunSpecs(t, "DeleteTenant method behaviour")
}
