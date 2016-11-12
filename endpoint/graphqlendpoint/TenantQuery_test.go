package graphqlendpoint_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/graphql-go/graphql"
	"github.com/microbusinesslimited/Micro-Business-Core/system"
	"github.com/microbusinesslimited/TenantService/business/domain"
	"github.com/microbusinesslimited/TenantService/endpoint/graphqlendpoint"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TenantQuery method input parameters and dependency test", func() {
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
			query := "{tenant{ID SecretKey}}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if TenantID format is not UUID", func() {
			query := "{tenant(tenantID:\"invalid UUID\"){ID SecretKey}}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
	})
})

var _ = Describe("TenantQuery method behaviour", func() {
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

	It("should call tenant service ReadTenant function", func() {
		randomValue, _ := system.RandomUUID()
		tenant := domain.Tenant{SecretKey: randomValue.String()}
		mockTenantService.EXPECT().ReadTenant(tenantID).Return(tenant, nil)

		query := "{tenant(tenantID:\"" + tenantID.String() + "\"){ID SecretKey}}"

		graphqlendpoint.ExecuteQuery(query, mockTenantService)
	})

	It("should return error if tenant service ReadTenant function returns error", func() {
		randomValue, _ := system.RandomUUID()
		mockTenantService.EXPECT().ReadTenant(tenantID).Return(domain.Tenant{}, fmt.Errorf(randomValue.String()))

		query := "{tenant(tenantID:\"" + tenantID.String() + "\"){ID SecretKey}}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(Equal(fmt.Errorf(randomValue.String())))
		Expect(result).To(BeNil())
	})

	It("should return tenant information if tenant service ReadTenant function returns an tenant information", func() {
		randomValue, _ := system.RandomUUID()
		tenant := domain.Tenant{SecretKey: randomValue.String()}
		mockTenantService.EXPECT().ReadTenant(tenantID).Return(tenant, nil)

		expectedTenant := &graphql.Result{
			Data: map[string]interface{}{
				"tenant": map[string]interface{}{
					"ID":        tenantID.String(),
					"SecretKey": randomValue.String(),
				},
			},
		}

		query := "{tenant(tenantID:\"" + tenantID.String() + "\"){ID SecretKey}}"

		returnedTenant, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(BeNil())
		Expect(returnedTenant).To(Equal(expectedTenant))
	})

	It("should return tenant information (only ID) if tenant service ReadTenant function returns an tenant information", func() {
		randomValue, _ := system.RandomUUID()
		tenant := domain.Tenant{SecretKey: randomValue.String()}
		mockTenantService.EXPECT().ReadTenant(tenantID).Return(tenant, nil)

		expectedTenant := &graphql.Result{
			Data: map[string]interface{}{
				"tenant": map[string]interface{}{
					"ID": tenantID.String(),
				},
			},
		}

		query := "{tenant(tenantID:\"" + tenantID.String() + "\"){ID}}"

		returnedTenant, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(BeNil())
		Expect(returnedTenant).To(Equal(expectedTenant))
	})

	It("should return tenant information (only SecretKey) if tenant service ReadTenant function returns an tenant information", func() {
		randomValue, _ := system.RandomUUID()
		tenant := domain.Tenant{SecretKey: randomValue.String()}
		mockTenantService.EXPECT().ReadTenant(tenantID).Return(tenant, nil)

		expectedTenant := &graphql.Result{
			Data: map[string]interface{}{
				"tenant": map[string]interface{}{
					"SecretKey": randomValue.String(),
				},
			},
		}

		query := "{tenant(tenantID:\"" + tenantID.String() + "\"){SecretKey}}"

		returnedTenant, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(BeNil())
		Expect(returnedTenant).To(Equal(expectedTenant))
	})
})

func TestTenantQuery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TenantQuery method input parameters and dependency test")
	RunSpecs(t, "TenantQuery method behaviour")
}
