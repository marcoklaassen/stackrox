package networkentities

import (
	"context"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/sensor/service/common"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/sync"
)

type controller struct {
	clusterID    string
	netEntityMgr common.NetworkEntityManager

	requestSeqID int64

	stopSig  concurrency.ReadOnlyErrorSignal
	injector common.MessageInjector

	lock sync.Mutex
}

func newController(clusterID string, netEntityMgr common.NetworkEntityManager, injector common.MessageInjector, stopSig concurrency.ReadOnlyErrorSignal) *controller {
	return &controller{
		clusterID:    clusterID,
		netEntityMgr: netEntityMgr,
		stopSig:      stopSig,
		injector:     injector,
	}
}

func (c *controller) SyncNow(ctx context.Context) error {
	msg, err := c.getPushNetworkEntitiesRequestMsg(ctx)
	if err != nil {
		return err
	}

	// Stop early if the request message is outdated.
	if msg.GetPushNetworkEntitiesRequest().GetSeqID() != atomic.LoadInt64(&c.requestSeqID) {
		return nil
	}

	if err := c.injector.InjectMessage(ctx, msg); err != nil {
		return errors.Wrap(err, "could not send external network entities")
	}
	return nil
}

func (c *controller) getPushNetworkEntitiesRequestMsg(ctx context.Context) (*central.MsgToSensor, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	requestSeqID := atomic.AddInt64(&c.requestSeqID, 1)

	netEntities, err := c.netEntityMgr.GetAllEntitiesForCluster(ctx, c.clusterID)
	if err != nil {
		return nil, errors.Wrap(err, "could not obtain external network entities to sync")
	}

	srcs := make([]*storage.NetworkEntityInfo, 0, len(netEntities))
	for _, entity := range netEntities {
		srcs = append(srcs, entity.GetInfo())
	}

	return &central.MsgToSensor{
		Msg: &central.MsgToSensor_PushNetworkEntitiesRequest{
			PushNetworkEntitiesRequest: &central.PushNetworkEntitiesRequest{
				Entities: srcs,
				SeqID:    requestSeqID,
			},
		},
	}, nil
}
