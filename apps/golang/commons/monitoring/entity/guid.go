package entity

import (
	"context"

	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/monitoring/newrelic/graphql"
)

const entityGuidQueryTemplate = `
{
	actor {
		entitySearch(query: "name = '{{ .EntityName }}'") {
			results {
				entities {
					guid
				}
			}
		}
	}
}
`

type entityGuidQueryVariables struct {
	EntityName string
}

type EntityGuidClient struct {
	logger        *logger.Logger
	graphQlClient *graphql.GraphQlClient
}

func NewEntityGuidClient(
	logger *logger.Logger,
	newrelicGraphQlEndpoint string,
	newrelicUserApiKey string,
) *EntityGuidClient {

	return &EntityGuidClient{
		logger:        logger,
		graphQlClient: graphql.New(logger, newrelicGraphQlEndpoint, newrelicUserApiKey),
	}
}

func (c *EntityGuidClient) Run(
	ctx context.Context,
	entityName string,
) (
	string,
	error,
) {

	// Create query variables for entity GUID
	qv := &entityGuidQueryVariables{
		EntityName: entityName,
	}

	// Perform GraphQL request to get entity GUID
	var res map[string]map[string]map[string]map[string]map[string][]map[string]string
	err := c.graphQlClient.Execute(ctx, "change_marker", entityGuidQueryTemplate, qv, &res)
	if err != nil {
		return "", err
	}
	// Parse entity GUID
	guid := res["data"]["actor"]["entitySearch"]["results"]["entities"][0]["guid"]
	return guid, nil
}
