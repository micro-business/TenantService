// +build integration

package service_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateApplication method behaviour", func() {
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
		It("should return the new UUID as application unique identifier and no error", func() {
			tenantID, _, err := createTenant(keyspace)
			Expect(err).To(BeNil())

			expectedApplicationID, _ := system.RandomUUID()
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(expectedApplicationID, nil)

			newApplicationID, err := tenantDataService.CreateApplication(tenantID, createApplicationInfo())

			Expect(newApplicationID).To(Equal(expectedApplicationID))
			Expect(err).To(BeNil())
		})
	})

	Context("when UUID generator service fails to create the new UUID", func() {
		It("should return application unique identifier as empty UUID and the returned error by tenant data service", func() {
			tenantID, _, err := createTenant(keyspace)
			Expect(err).To(BeNil())

			expectedErrorID, _ := system.RandomUUID()
			expectedError := errors.New(expectedErrorID.String())
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(system.EmptyUUID, expectedError)

			newApplicationID, err := tenantDataService.CreateApplication(tenantID, createApplicationInfo())

			Expect(newApplicationID).To(Equal(system.EmptyUUID))
			Expect(err).To(Equal(expectedError))
		})
	})

	Context("when creating new application", func() {
		It("should return error if tenant does not exist", func() {
			invalidTenantID, _ := system.RandomUUID()
			newApplicationID, err := tenantDataService.CreateApplication(invalidTenantID, createApplicationInfo())

			Expect(newApplicationID).To(Equal(system.EmptyUUID))
			Expect(err).To(Equal(fmt.Errorf("Tenant not found. Tenant ID: %s", invalidTenantID.String())))
		})

		It("should insert the record into application table", func() {
			tenantID, _, err := createTenant(keyspace)
			Expect(err).To(BeNil())

			applicationID, _ := system.RandomUUID()
			mockUUIDGeneratorService.
				EXPECT().
				GenerateRandomUUID().
				Return(applicationID, nil)

			application := createApplicationInfo()
			newApplicationID, err := tenantDataService.CreateApplication(tenantID, application)

			Expect(newApplicationID).To(Equal(applicationID))
			Expect(err).To(BeNil())

			config := getClusterConfig()
			config.Keyspace = keyspace

			session, err := config.CreateSession()

			defer session.Close()

			Expect(err).To(BeNil())

			iter := session.Query(
				"SELECT name"+
					" FROM application"+
					" WHERE"+
					" tenant_id = ?"+
					" AND application_id = ?",
				tenantID.String(),
				applicationID.String()).Iter()

			defer iter.Close()

			var name string

			Expect(iter.Scan(&name)).To(BeTrue())
			Expect(name).To(Equal(application.Name))
		})
	})
})

func TestCreateApplicationBehaviour(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateApplication method behaviour")
}
