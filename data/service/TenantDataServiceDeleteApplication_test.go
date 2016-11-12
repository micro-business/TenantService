package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/micro-business/Micro-Business-Core/system"
	"github.com/micro-business/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteApplication method input parameters and dependency test", func() {
	var (
		tenantDataService *service.TenantDataService
	)

	BeforeEach(func() {
		tenantDataService = &service.TenantDataService{ClusterConfig: &gocql.ClusterConfig{}}

	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService.ClusterConfig = nil

			tenantID, _ := system.RandomUUID()
			applicationID, _ := system.RandomUUID()

			Ω(func() { tenantDataService.DeleteApplication(tenantID, applicationID) }).Should(Panic())
		})
	})
})

func TestDeleteApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DeleteApplication method input parameters and dependency test")
}
