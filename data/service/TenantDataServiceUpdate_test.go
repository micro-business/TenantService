package service_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/contract"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update method input parameters and dependency test", func() {
	var (
		mockCtrl                                 *gomock.Controller
		tenantDataService                        *service.TenantDataService
		validTenantID                            system.UUID
		validTenant                              contract.Tenant
		tenantWithEmptySecretKey                 contract.Tenant
		tenantWithSecretKeyContainWhitespaceOnly contract.Tenant
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		tenantDataService = &service.TenantDataService{ClusterConfig: &gocql.ClusterConfig{}}
		validTenantID, _ = system.RandomUUID()

		randomValue, _ := system.RandomUUID()
		validTenant = contract.Tenant{SecretKey: randomValue.String()}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when cluster configuration not provided", func() {
		It("should panic", func() {
			tenantDataService.ClusterConfig = nil

			Ω(func() { tenantDataService.Update(validTenantID, validTenant) }).Should(Panic())
		})
	})

	DescribeTable("Input Parameters",
		func(tenantID system.UUID, tenant contract.Tenant) {
			Ω(func() { tenantDataService.Update(tenantID, tenant) }).Should(Panic())
		},
		Entry("should panic when empty tenant unique identifier provide", system.EmptyUUID, validTenant),
		Entry("should panic when tenant with empty secret key provided", validTenantID, tenantWithEmptySecretKey),
		Entry("should panic when tenant with secret key contains whitespace only provided", validTenantID, tenantWithSecretKeyContainWhitespaceOnly))
})

func TestUpdate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Update method input parameters and dependency test")
}
