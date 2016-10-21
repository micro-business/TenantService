package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadApplication method input parameters and dependency test", func() {
	var (
		mockCtrl           *gomock.Controller
		tenantDataService  *service.TenantDataService
		validTenantID      system.UUID
		validApplicationID system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		tenantDataService = &service.TenantDataService{ClusterConfig: &gocql.ClusterConfig{}}

		validTenantID, _ = system.RandomUUID()
		validApplicationID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService.ClusterConfig = nil

			Ω(func() { tenantDataService.ReadApplication(validTenantID, validApplicationID) }).Should(Panic())
		})
	})
})

func TestReadApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReadApplication method input parameters and dependency test")
}
