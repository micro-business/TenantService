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

var _ = Describe("ReadAllApplications method behaviour", func() {
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
		_, _, _, _, err := createApplication(keyspace)
		Expect(err).To(BeNil())

		invalidTenantID, _ := system.RandomUUID()
		applications, err := tenantDataService.ReadAllApplications(invalidTenantID)

		Expect(err).To(Equal(fmt.Errorf("Tenant not found. Tenant ID: %s", invalidTenantID.String())))
		Eventually(applications).Should(HaveLen(0))
	})

	It("should return empty list if tenant does not have any registered application", func() {
		tenantID, _, err := createTenant(keyspace)
		Expect(err).To(BeNil())

		applications, err := tenantDataService.ReadAllApplications(tenantID)

		Expect(err).To(BeNil())
		Eventually(applications).Should(HaveLen(0))
	})

	It("should return the existing tenant", func() {
		tenantID, _, expectedApplications, err := createApplications(keyspace)
		Expect(err).To(BeNil())

		returnedApplications, err := tenantDataService.ReadAllApplications(tenantID)

		Expect(err).To(BeNil())
		Expect(returnedApplications).To(Equal(expectedApplications))
	})
})

func TestReadAllApplicationsBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReadAllApplications method behaviour")
}
