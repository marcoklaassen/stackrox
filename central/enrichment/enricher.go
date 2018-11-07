package enrichment

import (
	deploymentDS "github.com/stackrox/rox/central/deployment/datastore"
	imageDS "github.com/stackrox/rox/central/image/datastore"
	imageIntegrationDS "github.com/stackrox/rox/central/imageintegration/datastore"
	multiplierDS "github.com/stackrox/rox/central/multiplier/store"
	"github.com/stackrox/rox/central/risk"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/images/enricher"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	logger = logging.LoggerForModule()
)

// Enricher enriches images with data from registries and scanners.
//go:generate mockgen-wrapper Enricher
type Enricher interface {
	Enrich(deployment *v1.Deployment) (bool, error)

	UpdateMultiplier(multiplier *v1.Multiplier)
	RemoveMultiplier(id string)

	ReprocessRiskAsync()
	ReprocessDeploymentRiskAsync(deployment *v1.Deployment)
}

// New creates and returns a new Enricher.
func New(deploymentStorage deploymentDS.DataStore,
	imageStorage imageDS.DataStore,
	imageIntegrationStorage imageIntegrationDS.DataStore,
	multiplierStorage multiplierDS.Store,
	imageEnricher enricher.ImageEnricher,
	scorer risk.Scorer) (Enricher, error) {
	e := &enricherImpl{
		deploymentStorage:       deploymentStorage,
		imageStorage:            imageStorage,
		imageIntegrationStorage: imageIntegrationStorage,
		multiplierStorage:       multiplierStorage,
		imageEnricher:           imageEnricher,
		scorer:                  scorer,
	}
	if err := e.initializeMultipliers(); err != nil {
		return nil, err
	}
	return e, nil
}
