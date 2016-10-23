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

var _ = Describe("CreateTenant method input parameters and dependency test", func() {
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
		It("should return error if no Tenant provided", func() {
			query := "mutation {createTenant}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
	})
})

var _ = Describe("CreateTenant method behaviour", func() {
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

	It("should call tenant service CreateTenant function", func() {
		mockTenantService.EXPECT().CreateTenant(tenant).Return(tenantID, nil)

		query := "mutation {createTenant (tenant: {SecretKey:\"" + tenant.SecretKey + "\"})}"

		graphqlendpoint.ExecuteQuery(query, mockTenantService)
	})

	It("should return error if tenant service CreateTenant function returns error", func() {
		randomValue, _ := system.RandomUUID()
		mockTenantService.EXPECT().CreateTenant(tenant).Return(system.EmptyUUID, fmt.Errorf(randomValue.String()))

		query := "mutation {createTenant (tenant: {SecretKey:\"" + tenant.SecretKey + "\"})}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(Equal(fmt.Errorf(randomValue.String())))
		Expect(result).To(BeNil())
	})

	It("should return tenant unique identifier if tenant service CreateTenant function returns no error", func() {
		mockTenantService.EXPECT().CreateTenant(tenant).Return(tenantID, nil)

		expectedTenant := &graphql.Result{
			Data: map[string]interface{}{
				"createTenant": tenantID.String(),
			},
		}

		query := "mutation {createTenant (tenant: {SecretKey:\"" + tenant.SecretKey + "\"})}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(expectedTenant))
	})
})

func TestCreateTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateTenant method input parameters and dependency test")
	RunSpecs(t, "CreateTenant method behaviour")
}
