// +build integration

package service_test

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/microbusinesslimited/Micro-Business-Core/system"
	"github.com/microbusinesslimited/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadTenant method behaviour", func() {
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
		invalidTenantID, _ := system.RandomUUID()
		_, err := tenantDataService.ReadTenant(invalidTenantID)

		Expect(err).To(Equal(fmt.Errorf("Tenant not found. Tenant ID: %s", invalidTenantID.String())))
	})

	It("should return the existing tenant", func() {
		tenantID, expectedTenant, err := createTenant(keyspace)
		Expect(err).To(BeNil())

		returnedTenant, err := tenantDataService.ReadTenant(tenantID)

		Expect(err).To(BeNil())
		Expect(returnedTenant).To(Equal(expectedTenant))
	})
})

func TestReadTenantBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReadTenant method behaviour")
}
