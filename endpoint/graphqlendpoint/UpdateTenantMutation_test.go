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

var _ = Describe("UpdateTenant method input parameters and dependency test", func() {
	var (
		mockCtrl          *gomock.Controller
		mockTenantService *MockTenantService
		tenantID          system.UUID
		tenant            domain.Tenant
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantService = NewMockTenantService(mockCtrl)

		tenantID, _ = system.RandomUUID()

		randomValue, _ := system.RandomUUID()
		tenant = domain.Tenant{SecretKey: randomValue.String()}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Input Parameters", func() {
		It("should return error if tenantID not provided", func() {
			query := "mutation {updateTenant (tenant: {Name:\"" + tenant.SecretKey + "\"})}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if tenantID format is invalid", func() {
			query := "mutation {updateTenant (tenantID: \"Invalid UUID\", tenant: {Name:\"" + tenant.SecretKey + "\"})}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if no Tenant provided", func() {
			query := "mutation {updateTenant(tenantID: \"" + tenantID.String() + "\")}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
	})
})

var _ = Describe("UpdateTenant method behaviour", func() {
	var (
		mockCtrl          *gomock.Controller
		mockTenantService *MockTenantService
		tenantID          system.UUID
		tenant            domain.Tenant
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantService = NewMockTenantService(mockCtrl)

		tenantID, _ = system.RandomUUID()

		randomValue, _ := system.RandomUUID()
		tenant = domain.Tenant{SecretKey: randomValue.String()}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call tenant service UpdateTenant function", func() {
		mockTenantService.EXPECT().UpdateTenant(tenantID, tenant).Return(nil)

		query := "mutation {updateTenant (tenantID: \"" + tenantID.String() + "\", tenant: {SecretKey:\"" + tenant.SecretKey + "\"})}"

		graphqlendpoint.ExecuteQuery(query, mockTenantService)
	})

	It("should return error if tenant service UpdateTenant function returns error", func() {
		randomValue, _ := system.RandomUUID()
		mockTenantService.EXPECT().UpdateTenant(tenantID, tenant).Return(fmt.Errorf(randomValue.String()))

		query := "mutation {updateTenant (tenantID: \"" + tenantID.String() + "\", tenant: {SecretKey:\"" + tenant.SecretKey + "\"})}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(Equal(fmt.Errorf(randomValue.String())))
		Expect(result).To(BeNil())
	})

	It("should return tenant unique identifier if tenant service UpdateTenant function returns no error", func() {
		mockTenantService.EXPECT().UpdateTenant(tenantID, tenant).Return(nil)

		expectedTenant := &graphql.Result{
			Data: map[string]interface{}{
				"updateTenant": true,
			},
		}

		query := "mutation {updateTenant (tenantID: \"" + tenantID.String() + "\", tenant: {SecretKey:\"" + tenant.SecretKey + "\"})}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(expectedTenant))
	})
})

func TestUpdateTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateTenant method input parameters and dependency test")
	RunSpecs(t, "UpdateTenant method behaviour")
}
