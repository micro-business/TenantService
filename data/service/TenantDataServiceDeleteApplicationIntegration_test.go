// +build integration

package service_test

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteApplication method behaviour", func() {
	var (
		mockCtrl                 *gomock.Controller
		tenantDataService        *service.TenantDataService
		mockUUIDGeneratorService *MockUUIDGeneratorService
		clusterConfig            *gocql.ClusterConfig
	)

	BeforeEach(func() {
		clusterConfig = getClusterConfig()
		clusterConfig.Keyspace = keyspace

		mockCtrl = gomock.NewController(GinkgoT())
		mockUUIDGeneratorService = NewMockUUIDGeneratorService(mockCtrl)

		tenantDataService = &service.TenantDataService{UUIDGeneratorService: mockUUIDGeneratorService, ClusterConfig: clusterConfig}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when deleting existing application", func() {
		It("should return error if tenant does not exist", func() {
			_, applicationID, err := createApplication(keyspace)
			Expect(err).To(BeNil())

			invalidTenantID, _ := system.RandomUUID()
			err = tenantDataService.DeleteApplication(invalidTenantID, applicationID)

			Expect(err).To(Equal(fmt.Errorf("Tenant not found. Tenant ID: %s", invalidTenantID.String())))
		})

		It("should return error if application does not exist", func() {
			tenantID, err := createTenant(keyspace)
			Expect(err).To(BeNil())

			invalidApplicationID, _ := system.RandomUUID()
			err = tenantDataService.DeleteApplication(tenantID, invalidApplicationID)

			Expect(err).To(Equal(fmt.Errorf("Tenant Application not found. Tenant ID: %s, Application ID: %s", tenantID.String(), invalidApplicationID.String())))
		})

		It("should remove the record from application table", func() {
			tenantID, applicationID, err := createApplication(keyspace)
			Expect(err).To(BeNil())

			err = tenantDataService.DeleteApplication(tenantID, applicationID)
			Expect(err).To(BeNil())

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
