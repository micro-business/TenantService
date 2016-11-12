package endpoint

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/micro-business/Micro-Business-Core/common/diagnostics"
	"github.com/micro-business/TenantService/business/contract"
	"github.com/micro-business/TenantService/config"
	"github.com/micro-business/TenantService/endpoint/graphqlendpoint"
	"golang.org/x/net/context"
)

// Endpoint implements method to start the service. The structure contains all the dependencies required by the Endpoint service.
type Endpoint struct {
	ConfigurationReader config.ConfigurationReader
	TenantService       contract.TenantService
}

// StartServer creates all the endpoints and starts the server.
func (endpoint Endpoint) StartServer() {
	diagnostics.IsNotNil(endpoint.TenantService, "endpoint.TenantService", "TenantService must be provided.")
	diagnostics.IsNotNil(endpoint.ConfigurationReader, "endpoint.ConfigurationReader", "ConfigurationReader must be provided.")

	http.Handle("/Api", httptransport.NewServer(
		context.Background(),
		createAPIEndpoint(endpoint.TenantService),
		decodeAPIRequest,
		encodeAPIResponse))

	if listeningPort, err := endpoint.ConfigurationReader.GetListeningPort(); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(listeningPort), nil))
	}
}

func createAPIEndpoint(tenantService contract.TenantService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return graphqlendpoint.ExecuteQuery(request.(string), tenantService)
	}
}

/// decodeAPIRequest decodes the request message sent by the client. The request message can be sent using GET method as part of
// URL or can be the payload of a POST HTTP message.
func decodeAPIRequest(context context.Context, httpRequest *http.Request) (interface{}, error) {
	if httpRequest.Method == "GET" {
		return httpRequest.URL.Query()["query"][0], nil
	}

	query, err := ioutil.ReadAll(httpRequest.Body)

	if err != nil {
		return nil, err
	}

	return string(query), nil
}

// encodeAPIResponse encodes the response message before sending back to the client
func encodeAPIResponse(context context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST")
	writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

	return json.NewEncoder(writer).Encode(response)
}
