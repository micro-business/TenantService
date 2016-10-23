package graphqlendpoint_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/business/domain"
	"github.com/microbusinesses/TenantService/endpoint/graphqlendpoint"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type application struct {
	ID   string
	Name string
}

var _ = Describe("ApplicationsQuery method input parameters and dependency test", func() {
	var (
		mockCtrl          *gomock.Controller
		mockTenantService *MockTenantService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantService = NewMockTenantService(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Input Parameters", func() {
		It("should return error if no TenantID provided", func() {
			query := "{applications{ID Name}}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return error if TenantID format is not UUID", func() {
			query := "{applications(tenantID:\"invalid UUID\"){ID Name}}"

			result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
	})
})

var _ = Describe("ApplicationsQuery method behaviour", func() {
	var (
		mockCtrl          *gomock.Controller
		mockTenantService *MockTenantService
		tenantID          system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantService = NewMockTenantService(mockCtrl)

		tenantID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call tenant service ReadAllApplications function", func() {
		applications := make(map[system.UUID]domain.Application)
		mockTenantService.EXPECT().ReadAllApplications(tenantID).Return(applications, nil)

		query := "{applications(tenantID:\"" + tenantID.String() + "\"){ID Name}}"

		graphqlendpoint.ExecuteQuery(query, mockTenantService)
	})

	It("should return error if tenant service ReadAllApplications function returns error", func() {
		randomValue, _ := system.RandomUUID()
		mockTenantService.EXPECT().ReadAllApplications(tenantID).Return(nil, fmt.Errorf(randomValue.String()))

		query := "{applications(tenantID:\"" + tenantID.String() + "\"){ID Name}}"

		result, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(Equal(fmt.Errorf(randomValue.String())))
		Expect(result).To(BeNil())
	})

	It("should return applications information if tenant service ReadAllApplications function returns applications information", func() {
		applicationCount := rand.Intn(10) + 1
		applications := make(map[system.UUID]domain.Application)
		expectedApps := make(map[string]application)

		for idx := 0; idx < applicationCount; idx++ {
			applicationID, _ := system.RandomUUID()
			applicationInfo := createApplicationInfo()
			applications[applicationID] = applicationInfo

			expectedApps[applicationID.String()] = application{applicationID.String(), applicationInfo.Name}
		}

		mockTenantService.EXPECT().ReadAllApplications(tenantID).Return(applications, nil)

		query := "{applications(tenantID:\"" + tenantID.String() + "\"){ID Name}}"

		returnedApplications, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(BeNil())

		returnedApps := make(map[string]application)

		for _, returnedApplicationInfo := range returnedApplications.(*graphql.Result).Data.(map[string]interface{})["applications"].([]interface{}) {
			appInfo := returnedApplicationInfo.(map[string]interface{})
			returnedApps[appInfo["ID"].(string)] = application{ID: appInfo["ID"].(string), Name: appInfo["Name"].(string)}
		}

		Expect(returnedApps).To(Equal(expectedApps))
	})
})

func TestApplicationsQuery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApplicationsQuery method input parameters and dependency test")
	RunSpecs(t, "ApplicationsQuery method behaviour")
}
