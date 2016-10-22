// +build integration

package service_test

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateApplication method behaviour", func() {
	var (
		tenantDataService *service.TenantDataService
		clusterConfig     *gocql.ClusterConfig
	)

	BeforeEach(func() {
		clusterConfig = getClusterConfig()
		clusterConfig.Keyspace = keyspace

		tenantDataService = &service.TenantDataService{ClusterConfig: clusterConfig}
	})

	Context("when updating an existing application", func() {
		It("should return error if tenant does not exist", func() {
			_, _, applicationID, _, err := createApplication(keyspace)
			Expect(err).To(BeNil())

			invalidTenantID, _ := system.RandomUUID()
			Expect(tenantDataService.UpdateApplication(invalidTenantID, applicationID, createApplicationInfo())).To(Equal(fmt.Errorf("Tenant not found. Tenant ID: %s", invalidTenantID.String())))
		})

		It("should return error if application does not exist", func() {
			tenantID, _, err := createTenant(keyspace)
			Expect(err).To(BeNil())

			invalidApplicationID, _ := system.RandomUUID()

			Expect(tenantDataService.UpdateApplication(tenantID, invalidApplicationID, createApplicationInfo())).To(Equal(fmt.Errorf("Tenant Application not found. Tenant ID: %s, Application ID: %s", tenantID.String(), invalidApplicationID.String())))
		})

		It("should update the record in application table", func() {
			tenantID, _, applicationID, _, err := createApplication(keyspace)
			Expect(err).To(BeNil())

			updatedTenant := createApplicationInfo()
			Expect(tenantDataService.UpdateApplication(tenantID, applicationID, updatedTenant)).To(BeNil())

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			defer session.Close()

			Expect(err).To(BeNil())

			iter := session.Query(
				"SELECT name"+
					" FROM application"+
					" WHERE"+
					" tenant_id = ?"+
					" AND application_id = ?",
				tenantID.String(),
				applicationID.String()).Iter()

			defer iter.Close()

			var name string

			Expect(iter.Scan(&name)).To(BeTrue())
			Expect(name).To(Equal(updatedTenant.Name))
		})
	})
})

func TestUpdateApplicationBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateApplication method behaviour")
}
