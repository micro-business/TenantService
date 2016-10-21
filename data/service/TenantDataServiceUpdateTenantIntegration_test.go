// +build integration

package service_test

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/contract"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateTenant method behaviour", func() {
	var (
		mockCtrl                 *gomock.Controller
		tenantDataService        *service.TenantDataService
		mockUUIDGeneratorService *MockUUIDGeneratorService
		validTenantID            system.UUID
		validTenant              contract.Tenant
		clusterConfig            *gocql.ClusterConfig
	)

	BeforeEach(func() {
		clusterConfig = getClusterConfig()
		clusterConfig.Keyspace = keyspace

		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = NewMockUUIDGeneratorService(mockCtrl)

		tenantDataService = &service.TenantDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: clusterConfig}

		validTenantID, _ = system.RandomUUID()

		randomValue, _ := system.RandomUUID()
		validTenant = contract.Tenant{SecretKey: randomValue.String()}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when updating an existing tenant", func() {
		It("should return error if tenant does not exist", func() {
			err := tenantDataService.UpdateTenant(validTenantID, validTenant)

			Expect(err).To(Equal(fmt.Errorf("Tenant not found. Tenant ID: %s", validTenantID.String())))
		})

		It("should update the record in tenant table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(validTenantID, nil)

			returnedTenantID, err := tenantDataService.CreateTenant(validTenant)

			Expect(err).To(BeNil())

			randomValue, _ := system.RandomUUID()
			updatedTenant := contract.Tenant{SecretKey: randomValue.String()}

			err = tenantDataService.UpdateTenant(returnedTenantID, updatedTenant)

			Expect(err).To(BeNil())

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			Expect(err).To(BeNil())

			defer session.Close()

			iter := session.Query(
				"SELECT secret_key"+
					" FROM tenant"+
					" WHERE"+
					" tenant_id = ?",
				validTenantID.String()).Iter()

			defer iter.Close()

			var secretKey string

			Expect(iter.Scan(&secretKey)).To(BeTrue())
			Expect(updatedTenant.SecretKey).To(Equal(secretKey))
		})
	})
})

func TestUpdateTenantBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateTenant method behaviour")
}
