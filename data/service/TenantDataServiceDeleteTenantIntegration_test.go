// +build integration

package service_test

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/micro-business/Micro-Business-Core/system"
	"github.com/micro-business/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteTenant method behaviour", func() {
	var (
		tenantDataService *service.TenantDataService
		clusterConfig     *gocql.ClusterConfig
	)

	BeforeEach(func() {
		clusterConfig = getClusterConfig()
		clusterConfig.Keyspace = keyspace

		tenantDataService = &service.TenantDataService{ClusterConfig: clusterConfig}
	})

	Context("when deleting existing tenant", func() {
		It("should return error if tenant does not exist", func() {
			invalidTenantID, _ := system.RandomUUID()
			Expect(tenantDataService.DeleteTenant(invalidTenantID)).To(Equal(fmt.Errorf("Tenant not found. Tenant ID: %s", invalidTenantID.String())))
		})

		It("should remove the record from tenant table", func() {
			tenantID, _, err := createTenant(keyspace)
			Expect(err).To(BeNil())

			err = tenantDataService.DeleteTenant(tenantID)
			Expect(err).To(BeNil())

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			defer session.Close()

			Expect(err).To(BeNil())

			var secretKey string

			iter := session.Query(
				"SELECT secret_key"+
					" FROM tenant"+
					" WHERE"+
					" tenant_id = ?",
				tenantID.String()).Iter()

			defer iter.Close()

			Expect(iter.Scan(&secretKey)).To(BeFalse())
		})
	})
})

func TestDeleteTenantBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DeleteTenant method behaviour")
}
