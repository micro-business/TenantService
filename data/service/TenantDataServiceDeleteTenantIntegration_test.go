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

var _ = Describe("DeleteTenant method behaviour", func() {
	var (
		mockCtrl                 *gomock.Controller
		tenantDataService        *service.TenantDataService
		mockUUIDGeneratorService *MockUUIDGeneratorService
		tenantID                 system.UUID
		validTenant              contract.Tenant
		clusterConfig            *gocql.ClusterConfig
	)

	BeforeEach(func() {
		clusterConfig = getClusterConfig()
		clusterConfig.Keyspace = keyspace

		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = NewMockUUIDGeneratorService(mockCtrl)

		tenantDataService = &service.TenantDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: clusterConfig}

		tenantID, _ = system.RandomUUID()

		randomValue, _ := system.RandomUUID()
		validTenant = contract.Tenant{SecretKey: randomValue.String()}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when deleting existing tenant", func() {
		It("should return error if tenant does not exist", func() {
			err := tenantDataService.DeleteTenant(tenantID)

			Expect(err).To(Equal(fmt.Errorf("Tenant not found. Tenant ID: %s", tenantID.String())))
		})

		It("should remove the record from tenant table", func() {
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(tenantID, nil)

			returnedTenantID, err := tenantDataService.CreateTenant(validTenant)

			Expect(err).To(BeNil())

			err = tenantDataService.DeleteTenant(returnedTenantID)

			Expect(err).To(BeNil())

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			Expect(err).To(BeNil())

			defer session.Close()

			var secretKey string

			iter := session.Query(
				"SELECT secret_key"+
					" FROM tenant"+
					" WHERE"+
					" tenant_id = ?",
				returnedTenantID.String()).Iter()

			defer iter.Close()

			Expect(iter.Scan(&secretKey)).To(BeFalse())
		})
	})
})

func TestDeleteTenantBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DeleteTenant method behaviour")
}
