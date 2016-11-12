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

var _ = Describe("DeleteApplication method behaviour", func() {
	var (
		tenantDataService *service.TenantDataService
		clusterConfig     *gocql.ClusterConfig
	)

	BeforeEach(func() {
		clusterConfig = getClusterConfig()
		clusterConfig.Keyspace = keyspace

		tenantDataService = &service.TenantDataService{ClusterConfig: clusterConfig}
	})

	Context("when deleting existing application", func() {
		It("should return error if tenant does not exist", func() {
			_, _, applicationID, _, err := createApplication(keyspace)
			Expect(err).To(BeNil())

			invalidTenantID, _ := system.RandomUUID()
			Expect(tenantDataService.DeleteApplication(invalidTenantID, applicationID)).To(Equal(fmt.Errorf("Tenant not found. Tenant ID: %s", invalidTenantID.String())))
		})

		It("should return error if application does not exist", func() {
			tenantID, _, err := createTenant(keyspace)
			Expect(err).To(BeNil())

			invalidApplicationID, _ := system.RandomUUID()
			Expect(tenantDataService.DeleteApplication(tenantID, invalidApplicationID)).To(Equal(fmt.Errorf("Tenant Application not found. Tenant ID: %s, Application ID: %s", tenantID.String(), invalidApplicationID.String())))
		})

		It("should remove the record from application table", func() {
			tenantID, _, applicationID, _, err := createApplication(keyspace)
			Expect(err).To(BeNil())

			Expect(tenantDataService.DeleteApplication(tenantID, applicationID)).To(BeNil())

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			defer session.Close()

			Expect(err).To(BeNil())

			var name string

			iter := session.Query(
				"SELECT name"+
					" FROM application"+
					" WHERE"+
					" tenant_id = ?"+
					" AND application_id = ?",
				tenantID.String(),
				applicationID.String()).Iter()

			defer iter.Close()

			Expect(iter.Scan(&name)).To(BeFalse())
		})
	})
})

func TestDeleteApplicationBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DeleteApplication method behaviour")
}
