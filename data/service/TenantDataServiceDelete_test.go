package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delete method input parameters and dependency test", func() {
	var (
		mockCtrl          *gomock.Controller
		tenantDataService *service.TenantDataService
		tenantID          system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		tenantDataService = &service.TenantDataService{ClusterConfig: &gocql.ClusterConfig{}}

		tenantID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService.ClusterConfig = nil

			Ω(func() { tenantDataService.Delete(tenantID) }).Should(Panic())
		})
	})

	DescribeTable("Input Parameters",
		func(tenantID system.UUID) {
			Ω(func() { tenantDataService.Delete(tenantID) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provided", system.EmptyUUID))
})

func TestDelete(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Delete method input parameters and dependency test")
}
