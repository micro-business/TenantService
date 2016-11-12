// +build integration

package service_test

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/micro-business/Micro-Business-Core/system"
	"github.com/micro-business/TenantService/data/contract"
	"github.com/micro-business/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadApplication method behaviour", func() {
	var (
		tenantDataService *service.TenantDataService
		clusterConfig     *gocql.ClusterConfig
	)

	BeforeEach(func() {
		clusterConfig = getClusterConfig()
		clusterConfig.Keyspace = keyspace

		tenantDataService = &service.TenantDataService{ClusterConfig: clusterConfig}
	})

	It("should return error if tenant does not exist", func() {
		_, _, applicationID, _, err := createApplication(keyspace)
		Expect(err).To(BeNil())

		invalidTenantID, _ := system.RandomUUID()
		application, err := tenantDataService.ReadApplication(invalidTenantID, applicationID)

		Expect(err).To(Equal(fmt.Errorf("Tenant not found. Tenant ID: %s", invalidTenantID.String())))
		Expect(application).To(Equal(contract.Application{}))
	})

	It("should return error if application does not exist", func() {
		tenantID, _, err := createTenant(keyspace)
		Expect(err).To(BeNil())

		invalidApplicationID, _ := system.RandomUUID()
		application, err := tenantDataService.ReadApplication(tenantID, invalidApplicationID)
		Expect(err).To(Equal(fmt.Errorf("Tenant Application not found. Tenant ID: %s, Application ID: %s", tenantID.String(), invalidApplicationID.String())))
		Expect(application).To(Equal(contract.Application{}))
	})

	It("should return the existing tenant", func() {
		tenantID, _, applicationID, expectedApplication, err := createApplication(keyspace)
		Expect(err).To(BeNil())

		returnedApplication, err := tenantDataService.ReadApplication(tenantID, applicationID)

		Expect(err).To(BeNil())
		Expect(returnedApplication).To(Equal(expectedApplication))
	})
})

func TestReadApplicationBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReadApplication method behaviour")
}
