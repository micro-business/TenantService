package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadApplication method input parameters and dependency test", func() {
	var (
		tenantDataService *service.TenantDataService
	)

	BeforeEach(func() {
		tenantDataService = &service.TenantDataService{ClusterConfig: &gocql.ClusterConfig{}}

	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService.ClusterConfig = nil

			validTenantID, _ := system.RandomUUID()
			validApplicationID, _ := system.RandomUUID()

			Ω(func() { tenantDataService.ReadApplication(validTenantID, validApplicationID) }).Should(Panic())
		})
	})
})

func TestReadApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReadApplication method input parameters and dependency test")
}
