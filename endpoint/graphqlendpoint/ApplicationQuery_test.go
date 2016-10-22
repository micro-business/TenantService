package graphqlendpoint_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/business/domain"
	"github.com/microbusinesses/TenantService/endpoint/graphqlendpoint"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ApplicationQuery method input parameters and dependency test", func() {
	var (
		mockCtrl          *gomock.Controller
		mockTenantService *MockTenantService
		tenantID          system.UUID
		applicationID     system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantService = NewMockTenantService(mockCtrl)

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Input Parameters", func() {
		It("should return error if no TenantID provided", func() {
			query := "{application(applicationID:\"" + applicationID.String() + "\"){ID Name}}"

			_, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
		})

		It("should return error if no ApplicationID provided", func() {
			query := "{application(tenantID:\"" + tenantID.String() + "\"){ID Name}}"

			_, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
		})

		It("should return error if no TenantID format is not UUID", func() {
			query := "{application(tenantID:\"invalid UUID\", applicationID:\"" + applicationID.String() + "\"){ID Name}}"

			_, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
		})

		It("should return error if no ApplicationID format is not UUID", func() {
			query := "{application(tenantID:\"" + tenantID.String() + "\", applicationID:\"invalid UUID\"){ID Name}}"

			_, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
			Expect(err).NotTo(BeNil())
		})
	})
})

var _ = Describe("ApplicationQuery method behaviour", func() {
	var (
		mockCtrl          *gomock.Controller
		mockTenantService *MockTenantService
		tenantID          system.UUID
		applicationID     system.UUID
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTenantService = NewMockTenantService(mockCtrl)

		tenantID, _ = system.RandomUUID()
		applicationID, _ = system.RandomUUID()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should call tenant service ReadApplication function", func() {
		randomValue, _ := system.RandomUUID()
		application := domain.Application{Name: randomValue.String()}
		mockTenantService.EXPECT().ReadApplication(tenantID, applicationID).Return(application, nil)

		query := "{application(tenantID:\"" + tenantID.String() + "\", applicationID:\"" + applicationID.String() + "\"){ID Name}}"

		graphqlendpoint.ExecuteQuery(query, mockTenantService)
	})

	It("should return error if tenant service ReadApplication function returns error", func() {
		randomValue, _ := system.RandomUUID()
		mockTenantService.EXPECT().ReadApplication(tenantID, applicationID).Return(domain.Application{}, fmt.Errorf(randomValue.String()))

		query := "{application(tenantID:\"" + tenantID.String() + "\", applicationID:\"" + applicationID.String() + "\"){ID Name}}"

		returnedApplication, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(Equal(fmt.Errorf(randomValue.String())))
		Expect(returnedApplication).To(BeNil())
	})

	It("should return application information if tenant service ReadApplication function returns an application information", func() {
		randomValue, _ := system.RandomUUID()
		application := domain.Application{Name: randomValue.String()}
		mockTenantService.EXPECT().ReadApplication(tenantID, applicationID).Return(application, nil)

		expectedApplication := &graphql.Result{
			Data: map[string]interface{}{
				"application": map[string]interface{}{
					"ID":   applicationID.String(),
					"Name": randomValue.String(),
				},
			},
		}

		query := "{application(tenantID:\"" + tenantID.String() + "\", applicationID:\"" + applicationID.String() + "\"){ID Name}}"

		returnedApplication, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(BeNil())
		Expect(returnedApplication).To(Equal(expectedApplication))
	})

	It("should return application information (only ID) if tenant service ReadApplication function returns an application information", func() {
		randomValue, _ := system.RandomUUID()
		application := domain.Application{Name: randomValue.String()}
		mockTenantService.EXPECT().ReadApplication(tenantID, applicationID).Return(application, nil)

		expectedApplication := &graphql.Result{
			Data: map[string]interface{}{
				"application": map[string]interface{}{
					"ID": applicationID.String(),
				},
			},
		}

		query := "{application(tenantID:\"" + tenantID.String() + "\", applicationID:\"" + applicationID.String() + "\"){ID}}"

		returnedApplication, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(BeNil())
		Expect(returnedApplication).To(Equal(expectedApplication))
	})

	It("should return application information (only Name) if tenant service ReadApplication function returns an application information", func() {
		randomValue, _ := system.RandomUUID()
		application := domain.Application{Name: randomValue.String()}
		mockTenantService.EXPECT().ReadApplication(tenantID, applicationID).Return(application, nil)

		expectedApplication := &graphql.Result{
			Data: map[string]interface{}{
				"application": map[string]interface{}{
					"Name": randomValue.String(),
				},
			},
		}

		query := "{application(tenantID:\"" + tenantID.String() + "\", applicationID:\"" + applicationID.String() + "\"){Name}}"

		returnedApplication, err := graphqlendpoint.ExecuteQuery(query, mockTenantService)
		Expect(err).To(BeNil())
		Expect(returnedApplication).To(Equal(expectedApplication))
	})
})

func TestApplicationQuery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApplicationQuery method input parameters and dependency test")
	RunSpecs(t, "ApplicationQuery method behaviour")
}
