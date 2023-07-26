// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/sac/resources"
	pgSearch "github.com/stackrox/rox/pkg/search/postgres"
	"gorm.io/gorm"
)

const (
	baseTable = "deployments"
	storeName = "Deployment"

	// using copyFrom, we may not even want to batch.  It would probably be simpler
	// to deal with failures if we just sent it all.  Something to think about as we
	// proceed and move into more e2e and larger performance testing
	batchSize = 10000
)

var (
	log            = logging.LoggerForModule()
	schema         = pkgSchema.DeploymentsSchema
	targetResource = resources.Deployment
)

type storeType = storage.Deployment

// Store is the interface to interact with the storage for storage.Deployment
type Store interface {
	Upsert(ctx context.Context, obj *storeType) error
	UpsertMany(ctx context.Context, objs []*storeType) error
	Delete(ctx context.Context, id string) error
	DeleteByQuery(ctx context.Context, q *v1.Query) error
	DeleteMany(ctx context.Context, identifiers []string) error

	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)

	Get(ctx context.Context, id string) (*storeType, bool, error)
	GetByQuery(ctx context.Context, query *v1.Query) ([]*storeType, error)
	GetMany(ctx context.Context, identifiers []string) ([]*storeType, []int, error)
	GetIDs(ctx context.Context) ([]string, error)

	Walk(ctx context.Context, fn func(obj *storeType) error) error
}

// New returns a new Store instance using the provided sql instance.
func New(db postgres.DB) Store {
	return pgSearch.NewGenericStore[storeType, *storeType](
		db,
		schema,
		pkGetter,
		insertIntoDeployments,
		copyFromDeployments,
		metricsSetAcquireDBConnDuration,
		metricsSetPostgresOperationDurationTime,
		isUpsertAllowed,
		targetResource,
	)
}

// region Helper functions

func pkGetter(obj *storeType) string {
	return obj.GetId()
}

func metricsSetPostgresOperationDurationTime(start time.Time, op ops.Op) {
	metrics.SetPostgresOperationDurationTime(start, op, storeName)
}

func metricsSetAcquireDBConnDuration(start time.Time, op ops.Op) {
	metrics.SetAcquireDBConnDuration(start, op, storeName)
}
func isUpsertAllowed(ctx context.Context, objs ...*storeType) error {
	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_WRITE_ACCESS).Resource(targetResource)
	if scopeChecker.IsAllowed() {
		return nil
	}
	var deniedIDs []string
	for _, obj := range objs {
		subScopeChecker := scopeChecker.ClusterID(obj.GetClusterId()).Namespace(obj.GetNamespace())
		if !subScopeChecker.IsAllowed() {
			deniedIDs = append(deniedIDs, obj.GetId())
		}
	}
	if len(deniedIDs) != 0 {
		return errors.Wrapf(sac.ErrResourceAccessDenied, "modifying deployments with IDs [%s] was denied", strings.Join(deniedIDs, ", "))
	}
	return nil
}

func insertIntoDeployments(ctx context.Context, batch *pgx.Batch, obj *storage.Deployment) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(obj.GetId()),
		obj.GetName(),
		obj.GetType(),
		obj.GetNamespace(),
		pgutils.NilOrUUID(obj.GetNamespaceId()),
		obj.GetOrchestratorComponent(),
		pgutils.EmptyOrMap(obj.GetLabels()),
		pgutils.EmptyOrMap(obj.GetPodLabels()),
		pgutils.NilOrTime(obj.GetCreated()),
		pgutils.NilOrUUID(obj.GetClusterId()),
		obj.GetClusterName(),
		pgutils.EmptyOrMap(obj.GetAnnotations()),
		obj.GetPriority(),
		obj.GetImagePullSecrets(),
		obj.GetServiceAccount(),
		obj.GetServiceAccountPermissionLevel(),
		obj.GetRiskScore(),
		serialized,
	}

	finalStr := "INSERT INTO deployments (Id, Name, Type, Namespace, NamespaceId, OrchestratorComponent, Labels, PodLabels, Created, ClusterId, ClusterName, Annotations, Priority, ImagePullSecrets, ServiceAccount, ServiceAccountPermissionLevel, RiskScore, serialized) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) ON CONFLICT(Id) DO UPDATE SET Id = EXCLUDED.Id, Name = EXCLUDED.Name, Type = EXCLUDED.Type, Namespace = EXCLUDED.Namespace, NamespaceId = EXCLUDED.NamespaceId, OrchestratorComponent = EXCLUDED.OrchestratorComponent, Labels = EXCLUDED.Labels, PodLabels = EXCLUDED.PodLabels, Created = EXCLUDED.Created, ClusterId = EXCLUDED.ClusterId, ClusterName = EXCLUDED.ClusterName, Annotations = EXCLUDED.Annotations, Priority = EXCLUDED.Priority, ImagePullSecrets = EXCLUDED.ImagePullSecrets, ServiceAccount = EXCLUDED.ServiceAccount, ServiceAccountPermissionLevel = EXCLUDED.ServiceAccountPermissionLevel, RiskScore = EXCLUDED.RiskScore, serialized = EXCLUDED.serialized"
	batch.Queue(finalStr, values...)

	var query string

	for childIndex, child := range obj.GetContainers() {
		if err := insertIntoDeploymentsContainers(ctx, batch, child, obj.GetId(), childIndex); err != nil {
			return err
		}
	}

	query = "delete from deployments_containers where deployments_Id = $1 AND idx >= $2"
	batch.Queue(query, pgutils.NilOrUUID(obj.GetId()), len(obj.GetContainers()))
	for childIndex, child := range obj.GetPorts() {
		if err := insertIntoDeploymentsPorts(ctx, batch, child, obj.GetId(), childIndex); err != nil {
			return err
		}
	}

	query = "delete from deployments_ports where deployments_Id = $1 AND idx >= $2"
	batch.Queue(query, pgutils.NilOrUUID(obj.GetId()), len(obj.GetPorts()))
	return nil
}

func insertIntoDeploymentsContainers(ctx context.Context, batch *pgx.Batch, obj *storage.Container, deploymentID string, idx int) error {

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(deploymentID),
		idx,
		obj.GetImage().GetId(),
		obj.GetImage().GetName().GetRegistry(),
		obj.GetImage().GetName().GetRemote(),
		obj.GetImage().GetName().GetTag(),
		obj.GetImage().GetName().GetFullName(),
		obj.GetSecurityContext().GetPrivileged(),
		obj.GetSecurityContext().GetDropCapabilities(),
		obj.GetSecurityContext().GetAddCapabilities(),
		obj.GetSecurityContext().GetReadOnlyRootFilesystem(),
		obj.GetResources().GetCpuCoresRequest(),
		obj.GetResources().GetCpuCoresLimit(),
		obj.GetResources().GetMemoryMbRequest(),
		obj.GetResources().GetMemoryMbLimit(),
	}

	finalStr := "INSERT INTO deployments_containers (deployments_Id, idx, Image_Id, Image_Name_Registry, Image_Name_Remote, Image_Name_Tag, Image_Name_FullName, SecurityContext_Privileged, SecurityContext_DropCapabilities, SecurityContext_AddCapabilities, SecurityContext_ReadOnlyRootFilesystem, Resources_CpuCoresRequest, Resources_CpuCoresLimit, Resources_MemoryMbRequest, Resources_MemoryMbLimit) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) ON CONFLICT(deployments_Id, idx) DO UPDATE SET deployments_Id = EXCLUDED.deployments_Id, idx = EXCLUDED.idx, Image_Id = EXCLUDED.Image_Id, Image_Name_Registry = EXCLUDED.Image_Name_Registry, Image_Name_Remote = EXCLUDED.Image_Name_Remote, Image_Name_Tag = EXCLUDED.Image_Name_Tag, Image_Name_FullName = EXCLUDED.Image_Name_FullName, SecurityContext_Privileged = EXCLUDED.SecurityContext_Privileged, SecurityContext_DropCapabilities = EXCLUDED.SecurityContext_DropCapabilities, SecurityContext_AddCapabilities = EXCLUDED.SecurityContext_AddCapabilities, SecurityContext_ReadOnlyRootFilesystem = EXCLUDED.SecurityContext_ReadOnlyRootFilesystem, Resources_CpuCoresRequest = EXCLUDED.Resources_CpuCoresRequest, Resources_CpuCoresLimit = EXCLUDED.Resources_CpuCoresLimit, Resources_MemoryMbRequest = EXCLUDED.Resources_MemoryMbRequest, Resources_MemoryMbLimit = EXCLUDED.Resources_MemoryMbLimit"
	batch.Queue(finalStr, values...)

	var query string

	for childIndex, child := range obj.GetConfig().GetEnv() {
		if err := insertIntoDeploymentsContainersEnvs(ctx, batch, child, deploymentID, idx, childIndex); err != nil {
			return err
		}
	}

	query = "delete from deployments_containers_envs where deployments_Id = $1 AND deployments_containers_idx = $2 AND idx >= $3"
	batch.Queue(query, pgutils.NilOrUUID(deploymentID), idx, len(obj.GetConfig().GetEnv()))
	for childIndex, child := range obj.GetVolumes() {
		if err := insertIntoDeploymentsContainersVolumes(ctx, batch, child, deploymentID, idx, childIndex); err != nil {
			return err
		}
	}

	query = "delete from deployments_containers_volumes where deployments_Id = $1 AND deployments_containers_idx = $2 AND idx >= $3"
	batch.Queue(query, pgutils.NilOrUUID(deploymentID), idx, len(obj.GetVolumes()))
	for childIndex, child := range obj.GetSecrets() {
		if err := insertIntoDeploymentsContainersSecrets(ctx, batch, child, deploymentID, idx, childIndex); err != nil {
			return err
		}
	}

	query = "delete from deployments_containers_secrets where deployments_Id = $1 AND deployments_containers_idx = $2 AND idx >= $3"
	batch.Queue(query, pgutils.NilOrUUID(deploymentID), idx, len(obj.GetSecrets()))
	return nil
}

func insertIntoDeploymentsContainersEnvs(_ context.Context, batch *pgx.Batch, obj *storage.ContainerConfig_EnvironmentConfig, deploymentID string, deploymentContainerIdx int, idx int) error {

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(deploymentID),
		deploymentContainerIdx,
		idx,
		obj.GetKey(),
		obj.GetValue(),
		obj.GetEnvVarSource(),
	}

	finalStr := "INSERT INTO deployments_containers_envs (deployments_Id, deployments_containers_idx, idx, Key, Value, EnvVarSource) VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT(deployments_Id, deployments_containers_idx, idx) DO UPDATE SET deployments_Id = EXCLUDED.deployments_Id, deployments_containers_idx = EXCLUDED.deployments_containers_idx, idx = EXCLUDED.idx, Key = EXCLUDED.Key, Value = EXCLUDED.Value, EnvVarSource = EXCLUDED.EnvVarSource"
	batch.Queue(finalStr, values...)

	return nil
}

func insertIntoDeploymentsContainersVolumes(_ context.Context, batch *pgx.Batch, obj *storage.Volume, deploymentID string, deploymentContainerIdx int, idx int) error {

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(deploymentID),
		deploymentContainerIdx,
		idx,
		obj.GetName(),
		obj.GetSource(),
		obj.GetDestination(),
		obj.GetReadOnly(),
		obj.GetType(),
	}

	finalStr := "INSERT INTO deployments_containers_volumes (deployments_Id, deployments_containers_idx, idx, Name, Source, Destination, ReadOnly, Type) VALUES($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT(deployments_Id, deployments_containers_idx, idx) DO UPDATE SET deployments_Id = EXCLUDED.deployments_Id, deployments_containers_idx = EXCLUDED.deployments_containers_idx, idx = EXCLUDED.idx, Name = EXCLUDED.Name, Source = EXCLUDED.Source, Destination = EXCLUDED.Destination, ReadOnly = EXCLUDED.ReadOnly, Type = EXCLUDED.Type"
	batch.Queue(finalStr, values...)

	return nil
}

func insertIntoDeploymentsContainersSecrets(_ context.Context, batch *pgx.Batch, obj *storage.EmbeddedSecret, deploymentID string, deploymentContainerIdx int, idx int) error {

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(deploymentID),
		deploymentContainerIdx,
		idx,
		obj.GetName(),
		obj.GetPath(),
	}

	finalStr := "INSERT INTO deployments_containers_secrets (deployments_Id, deployments_containers_idx, idx, Name, Path) VALUES($1, $2, $3, $4, $5) ON CONFLICT(deployments_Id, deployments_containers_idx, idx) DO UPDATE SET deployments_Id = EXCLUDED.deployments_Id, deployments_containers_idx = EXCLUDED.deployments_containers_idx, idx = EXCLUDED.idx, Name = EXCLUDED.Name, Path = EXCLUDED.Path"
	batch.Queue(finalStr, values...)

	return nil
}

func insertIntoDeploymentsPorts(ctx context.Context, batch *pgx.Batch, obj *storage.PortConfig, deploymentID string, idx int) error {

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(deploymentID),
		idx,
		obj.GetContainerPort(),
		obj.GetProtocol(),
		obj.GetExposure(),
	}

	finalStr := "INSERT INTO deployments_ports (deployments_Id, idx, ContainerPort, Protocol, Exposure) VALUES($1, $2, $3, $4, $5) ON CONFLICT(deployments_Id, idx) DO UPDATE SET deployments_Id = EXCLUDED.deployments_Id, idx = EXCLUDED.idx, ContainerPort = EXCLUDED.ContainerPort, Protocol = EXCLUDED.Protocol, Exposure = EXCLUDED.Exposure"
	batch.Queue(finalStr, values...)

	var query string

	for childIndex, child := range obj.GetExposureInfos() {
		if err := insertIntoDeploymentsPortsExposureInfos(ctx, batch, child, deploymentID, idx, childIndex); err != nil {
			return err
		}
	}

	query = "delete from deployments_ports_exposure_infos where deployments_Id = $1 AND deployments_ports_idx = $2 AND idx >= $3"
	batch.Queue(query, pgutils.NilOrUUID(deploymentID), idx, len(obj.GetExposureInfos()))
	return nil
}

func insertIntoDeploymentsPortsExposureInfos(_ context.Context, batch *pgx.Batch, obj *storage.PortConfig_ExposureInfo, deploymentID string, deploymentPortIdx int, idx int) error {

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(deploymentID),
		deploymentPortIdx,
		idx,
		obj.GetLevel(),
		obj.GetServiceName(),
		obj.GetServicePort(),
		obj.GetNodePort(),
		obj.GetExternalIps(),
		obj.GetExternalHostnames(),
	}

	finalStr := "INSERT INTO deployments_ports_exposure_infos (deployments_Id, deployments_ports_idx, idx, Level, ServiceName, ServicePort, NodePort, ExternalIps, ExternalHostnames) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT(deployments_Id, deployments_ports_idx, idx) DO UPDATE SET deployments_Id = EXCLUDED.deployments_Id, deployments_ports_idx = EXCLUDED.deployments_ports_idx, idx = EXCLUDED.idx, Level = EXCLUDED.Level, ServiceName = EXCLUDED.ServiceName, ServicePort = EXCLUDED.ServicePort, NodePort = EXCLUDED.NodePort, ExternalIps = EXCLUDED.ExternalIps, ExternalHostnames = EXCLUDED.ExternalHostnames"
	batch.Queue(finalStr, values...)

	return nil
}

func copyFromDeployments(ctx context.Context, s pgSearch.Deleter, tx *postgres.Tx, objs ...*storage.Deployment) error {
	inputRows := make([][]interface{}, 0, batchSize)

	// This is a copy so first we must delete the rows and re-add them
	// Which is essentially the desired behaviour of an upsert.
	deletes := make([]string, 0, batchSize)

	copyCols := []string{
		"id",
		"name",
		"type",
		"namespace",
		"namespaceid",
		"orchestratorcomponent",
		"labels",
		"podlabels",
		"created",
		"clusterid",
		"clustername",
		"annotations",
		"priority",
		"imagepullsecrets",
		"serviceaccount",
		"serviceaccountpermissionlevel",
		"riskscore",
		"serialized",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		serialized, marshalErr := obj.Marshal()
		if marshalErr != nil {
			return marshalErr
		}

		inputRows = append(inputRows, []interface{}{
			pgutils.NilOrUUID(obj.GetId()),
			obj.GetName(),
			obj.GetType(),
			obj.GetNamespace(),
			pgutils.NilOrUUID(obj.GetNamespaceId()),
			obj.GetOrchestratorComponent(),
			pgutils.EmptyOrMap(obj.GetLabels()),
			pgutils.EmptyOrMap(obj.GetPodLabels()),
			pgutils.NilOrTime(obj.GetCreated()),
			pgutils.NilOrUUID(obj.GetClusterId()),
			obj.GetClusterName(),
			pgutils.EmptyOrMap(obj.GetAnnotations()),
			obj.GetPriority(),
			obj.GetImagePullSecrets(),
			obj.GetServiceAccount(),
			obj.GetServiceAccountPermissionLevel(),
			obj.GetRiskScore(),
			serialized,
		})

		// Add the ID to be deleted.
		deletes = append(deletes, obj.GetId())

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if err := s.DeleteMany(ctx, deletes); err != nil {
				return err
			}
			// clear the inserts and vals for the next batch
			deletes = deletes[:0]

			if _, err := tx.CopyFrom(ctx, pgx.Identifier{"deployments"}, copyCols, pgx.CopyFromRows(inputRows)); err != nil {
				return err
			}
			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	for idx, obj := range objs {
		_ = idx // idx may or may not be used depending on how nested we are, so avoid compile-time errors.

		if err := copyFromDeploymentsContainers(ctx, s, tx, obj.GetId(), obj.GetContainers()...); err != nil {
			return err
		}
		if err := copyFromDeploymentsPorts(ctx, s, tx, obj.GetId(), obj.GetPorts()...); err != nil {
			return err
		}
	}

	return nil
}

func copyFromDeploymentsContainers(ctx context.Context, s pgSearch.Deleter, tx *postgres.Tx, deploymentID string, objs ...*storage.Container) error {
	inputRows := make([][]interface{}, 0, batchSize)

	copyCols := []string{
		"deployments_id",
		"idx",
		"image_id",
		"image_name_registry",
		"image_name_remote",
		"image_name_tag",
		"image_name_fullname",
		"securitycontext_privileged",
		"securitycontext_dropcapabilities",
		"securitycontext_addcapabilities",
		"securitycontext_readonlyrootfilesystem",
		"resources_cpucoresrequest",
		"resources_cpucoreslimit",
		"resources_memorymbrequest",
		"resources_memorymblimit",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{
			pgutils.NilOrUUID(deploymentID),
			idx,
			obj.GetImage().GetId(),
			obj.GetImage().GetName().GetRegistry(),
			obj.GetImage().GetName().GetRemote(),
			obj.GetImage().GetName().GetTag(),
			obj.GetImage().GetName().GetFullName(),
			obj.GetSecurityContext().GetPrivileged(),
			obj.GetSecurityContext().GetDropCapabilities(),
			obj.GetSecurityContext().GetAddCapabilities(),
			obj.GetSecurityContext().GetReadOnlyRootFilesystem(),
			obj.GetResources().GetCpuCoresRequest(),
			obj.GetResources().GetCpuCoresLimit(),
			obj.GetResources().GetMemoryMbRequest(),
			obj.GetResources().GetMemoryMbLimit(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if _, err := tx.CopyFrom(ctx, pgx.Identifier{"deployments_containers"}, copyCols, pgx.CopyFromRows(inputRows)); err != nil {
				return err
			}
			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	for idx, obj := range objs {
		_ = idx // idx may or may not be used depending on how nested we are, so avoid compile-time errors.

		if err := copyFromDeploymentsContainersEnvs(ctx, s, tx, deploymentID, idx, obj.GetConfig().GetEnv()...); err != nil {
			return err
		}
		if err := copyFromDeploymentsContainersVolumes(ctx, s, tx, deploymentID, idx, obj.GetVolumes()...); err != nil {
			return err
		}
		if err := copyFromDeploymentsContainersSecrets(ctx, s, tx, deploymentID, idx, obj.GetSecrets()...); err != nil {
			return err
		}
	}

	return nil
}

func copyFromDeploymentsContainersEnvs(ctx context.Context, s pgSearch.Deleter, tx *postgres.Tx, deploymentID string, deploymentContainerIdx int, objs ...*storage.ContainerConfig_EnvironmentConfig) error {
	inputRows := make([][]interface{}, 0, batchSize)

	copyCols := []string{
		"deployments_id",
		"deployments_containers_idx",
		"idx",
		"key",
		"value",
		"envvarsource",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{
			pgutils.NilOrUUID(deploymentID),
			deploymentContainerIdx,
			idx,
			obj.GetKey(),
			obj.GetValue(),
			obj.GetEnvVarSource(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if _, err := tx.CopyFrom(ctx, pgx.Identifier{"deployments_containers_envs"}, copyCols, pgx.CopyFromRows(inputRows)); err != nil {
				return err
			}
			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return nil
}

func copyFromDeploymentsContainersVolumes(ctx context.Context, s pgSearch.Deleter, tx *postgres.Tx, deploymentID string, deploymentContainerIdx int, objs ...*storage.Volume) error {
	inputRows := make([][]interface{}, 0, batchSize)

	copyCols := []string{
		"deployments_id",
		"deployments_containers_idx",
		"idx",
		"name",
		"source",
		"destination",
		"readonly",
		"type",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{
			pgutils.NilOrUUID(deploymentID),
			deploymentContainerIdx,
			idx,
			obj.GetName(),
			obj.GetSource(),
			obj.GetDestination(),
			obj.GetReadOnly(),
			obj.GetType(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if _, err := tx.CopyFrom(ctx, pgx.Identifier{"deployments_containers_volumes"}, copyCols, pgx.CopyFromRows(inputRows)); err != nil {
				return err
			}
			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return nil
}

func copyFromDeploymentsContainersSecrets(ctx context.Context, s pgSearch.Deleter, tx *postgres.Tx, deploymentID string, deploymentContainerIdx int, objs ...*storage.EmbeddedSecret) error {
	inputRows := make([][]interface{}, 0, batchSize)

	copyCols := []string{
		"deployments_id",
		"deployments_containers_idx",
		"idx",
		"name",
		"path",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{
			pgutils.NilOrUUID(deploymentID),
			deploymentContainerIdx,
			idx,
			obj.GetName(),
			obj.GetPath(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if _, err := tx.CopyFrom(ctx, pgx.Identifier{"deployments_containers_secrets"}, copyCols, pgx.CopyFromRows(inputRows)); err != nil {
				return err
			}
			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return nil
}

func copyFromDeploymentsPorts(ctx context.Context, s pgSearch.Deleter, tx *postgres.Tx, deploymentID string, objs ...*storage.PortConfig) error {
	inputRows := make([][]interface{}, 0, batchSize)

	copyCols := []string{
		"deployments_id",
		"idx",
		"containerport",
		"protocol",
		"exposure",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{
			pgutils.NilOrUUID(deploymentID),
			idx,
			obj.GetContainerPort(),
			obj.GetProtocol(),
			obj.GetExposure(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if _, err := tx.CopyFrom(ctx, pgx.Identifier{"deployments_ports"}, copyCols, pgx.CopyFromRows(inputRows)); err != nil {
				return err
			}
			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	for idx, obj := range objs {
		_ = idx // idx may or may not be used depending on how nested we are, so avoid compile-time errors.

		if err := copyFromDeploymentsPortsExposureInfos(ctx, s, tx, deploymentID, idx, obj.GetExposureInfos()...); err != nil {
			return err
		}
	}

	return nil
}

func copyFromDeploymentsPortsExposureInfos(ctx context.Context, s pgSearch.Deleter, tx *postgres.Tx, deploymentID string, deploymentPortIdx int, objs ...*storage.PortConfig_ExposureInfo) error {
	inputRows := make([][]interface{}, 0, batchSize)

	copyCols := []string{
		"deployments_id",
		"deployments_ports_idx",
		"idx",
		"level",
		"servicename",
		"serviceport",
		"nodeport",
		"externalips",
		"externalhostnames",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{
			pgutils.NilOrUUID(deploymentID),
			deploymentPortIdx,
			idx,
			obj.GetLevel(),
			obj.GetServiceName(),
			obj.GetServicePort(),
			obj.GetNodePort(),
			obj.GetExternalIps(),
			obj.GetExternalHostnames(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if _, err := tx.CopyFrom(ctx, pgx.Identifier{"deployments_ports_exposure_infos"}, copyCols, pgx.CopyFromRows(inputRows)); err != nil {
				return err
			}
			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return nil
}

// endregion Helper functions

// region Used for testing

// CreateTableAndNewStore returns a new Store instance for testing.
func CreateTableAndNewStore(ctx context.Context, db postgres.DB, gormDB *gorm.DB) Store {
	pkgSchema.ApplySchemaForTable(ctx, gormDB, baseTable)
	return New(db)
}

// Destroy drops the tables associated with the target object type.
func Destroy(ctx context.Context, db postgres.DB) {
	dropTableDeployments(ctx, db)
}

func dropTableDeployments(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments CASCADE")
	dropTableDeploymentsContainers(ctx, db)
	dropTableDeploymentsPorts(ctx, db)

}

func dropTableDeploymentsContainers(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments_containers CASCADE")
	dropTableDeploymentsContainersEnvs(ctx, db)
	dropTableDeploymentsContainersVolumes(ctx, db)
	dropTableDeploymentsContainersSecrets(ctx, db)

}

func dropTableDeploymentsContainersEnvs(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments_containers_envs CASCADE")

}

func dropTableDeploymentsContainersVolumes(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments_containers_volumes CASCADE")

}

func dropTableDeploymentsContainersSecrets(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments_containers_secrets CASCADE")

}

func dropTableDeploymentsPorts(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments_ports CASCADE")
	dropTableDeploymentsPortsExposureInfos(ctx, db)

}

func dropTableDeploymentsPortsExposureInfos(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments_ports_exposure_infos CASCADE")

}

// endregion Used for testing
