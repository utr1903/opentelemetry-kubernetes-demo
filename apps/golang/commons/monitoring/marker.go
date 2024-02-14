package monitoring

import (
	"context"

	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	newrelicmarker "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/monitoring/newrelic/marker"
)

const NEWRELIC = "newrelic"

type Marker interface {
	Run(
		ctx context.Context,
		entityName string,
		changeLog string,
		commit string,
		description string,
		groupId string,
		version string,
	) error
}

func NewMarker(
	logger *logger.Logger,
	observabilityBackendName string,
	observabilityBackendEndpoint string,
	observabilityBackendApiKey string,
) Marker {
	if observabilityBackendName == NEWRELIC {
		return newrelicmarker.New(logger, observabilityBackendEndpoint, observabilityBackendApiKey)
	}

	return nil
}
