// +build integration

package service_test

import (
	"errors"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesslimited/Micro-Business-Core/system"
	"github.com/microbusinesslimited/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateTenant method behaviour", func() {
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

	Context("when UUID generator service succeeds to create the new UUID", func() {
		It("should return the new UUID as tenant uniuqe identifier and no error", func() {
			expectedTenantID, _ := system.RandomUUID()
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(expectedTenantID, nil)

			newTenantID, err := tenantDataService.CreateTenant(createTenantInfo())

			Expect(newTenantID).To(Equal(expectedTenantID))
			Expect(err).To(BeNil())
		})
	})

	Context("when UUID generator service fails to create the new UUID", func() {
		It("should return tenant unique identifier as empty UUID and the returned error by tenant data service", func() {
			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(system.EmptyUUID, expectedError)

			newTenantID, err := tenantDataService.CreateTenant(createTenantInfo())

			Expect(newTenantID).To(Equal(system.EmptyUUID))
			Expect(err).To(Equal(expectedError))
		})
	})

	Context("when creating new tenant", func() {
		It("should insert the record into tenant table", func() {
			tenantID, _ := system.RandomUUID()
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(tenantID, nil)

			tenant := createTenantInfo()
			newTenantID, err := tenantDataService.CreateTenant(tenant)

			Expect(newTenantID).To(Equal(tenantID))
			Expect(err).To(BeNil())

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			defer session.Close()

			Expect(err).To(BeNil())

			iter := session.Query(
				"SELECT secret_key"+
					" FROM tenant"+
					" WHERE"+
					" tenant_id = ?",
				tenantID.String()).Iter()

			defer iter.Close()

			var secretKey string

			Expect(iter.Scan(&secretKey)).To(BeTrue())
			Expect(secretKey).To(Equal(tenant.SecretKey))
		})
	})
})

func TestCreateTenantBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateTenant method behaviour")
}
