// +build integration

package service_test

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/contract"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateTenant method behaviour", func() {
	var (
		tenantDataService *service.TenantDataService
		clusterConfig     *gocql.ClusterConfig
	)

	BeforeEach(func() {
		clusterConfig = getClusterConfig()
		clusterConfig.Keyspace = keyspace

		tenantDataService = &service.TenantDataService{ClusterConfig: clusterConfig}
	})

	Context("when updating an existing tenant", func() {
		It("should return error if tenant does not exist", func() {
			randomValue, _ := system.RandomUUID()
			tenant := contract.Tenant{SecretKey: randomValue.String()}

			invalidTenantID, _ := system.RandomUUID()
			err := tenantDataService.UpdateTenant(invalidTenantID, tenant)

			Expect(err).To(Equal(fmt.Errorf("Tenant not found. Tenant ID: %s", invalidTenantID.String())))
		})

		It("should update the record in tenant table", func() {
			tenantID, _, err := createTenant(keyspace)
			Expect(err).To(BeNil())

			randomValue, _ := system.RandomUUID()
			updatedTenant := contract.Tenant{SecretKey: randomValue.String()}

			Expect(tenantDataService.UpdateTenant(tenantID, updatedTenant)).To(BeNil())

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			defer session.Close()

			Expect(err).To(BeNil())

			iter := session.Query(
				"SELECT secret_key"+
					" FROM tenant"+
					" WHERE"+
					" tenant_id = ?",
				tenantID.String()).Iter()

			defer iter.Close()

			var secretKey string

			Expect(iter.Scan(&secretKey)).To(BeTrue())
			Expect(secretKey).To(Equal(updatedTenant.SecretKey))
		})
	})
})

func TestUpdateTenantBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateTenant method behaviour")
}
