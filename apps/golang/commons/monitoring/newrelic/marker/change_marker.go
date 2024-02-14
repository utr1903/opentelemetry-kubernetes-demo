package newrelic

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/monitoring/entity"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/monitoring/newrelic/graphql"
)

const changeMarkerQueryTemplate = `
mutation {
  changeTrackingCreateDeployment(
    deployment: {
      changelog: "{{ .ChangeLog }}",
      commit: "{{ .Commit }}",
      description: "{{ .Description }}",
      entityGuid: "{{ .EntityGuid }}",
      groupId: "{{ .GroupId }}",
      version: "{{ .Version }}"
    }
  ) {
    deploymentId
  }
}
`

type changeMarkerQueryVariables struct {
	ChangeLog   string
	Commit      string
	Description string
	EntityGuid  string
	GroupId     string
	Version     string
}

type ChangeMarkerClient struct {
	logger           *logger.Logger
	entityGuidClient *entity.EntityGuidClient
	graphQlClient    *graphql.GraphQlClient
}

func New(
	logger *logger.Logger,
	newrelicGraphQlEndpoint string,
	newrelicUserApiKey string,
) *ChangeMarkerClient {

	return &ChangeMarkerClient{
		logger:           logger,
		entityGuidClient: entity.NewEntityGuidClient(logger, newrelicGraphQlEndpoint, newrelicUserApiKey),
		graphQlClient:    graphql.New(logger, newrelicGraphQlEndpoint, newrelicUserApiKey),
	}
}

func (c *ChangeMarkerClient) Run(
	ctx context.Context,
	entityName string,
	changeLog string,
	commit string,
	description string,
	groupId string,
	version string,
) error {

	// Perform GraphQL request to get entity GUID
	entityGuid, err := c.entityGuidClient.Run(ctx, entityName)
	if err != nil {
		return err
	}

	// Create query variables for change marker
	qv := &changeMarkerQueryVariables{
		ChangeLog:   changeLog,
		Commit:      commit,
		Description: description,
		EntityGuid:  entityGuid,
		GroupId:     groupId,
		Version:     version,
	}

	// Perform GraphQL request to deploy change marker
	var res map[string]map[string]map[string]string
	err = c.graphQlClient.Execute(ctx, "change_marker", changeMarkerQueryTemplate, qv, &res)
	if err != nil {
		return err
	}
	deploymentId := res["data"]["changeTrackingCreateDeployment"]["deploymentId"]
	c.logger.Log(logrus.InfoLevel, ctx, "", "Change markes is deployed successfully. Deployment ID: ["+deploymentId+"]")
	return nil
}
