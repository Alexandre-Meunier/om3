// Package svcagg is responsible for of object.AggregatedStatus
//
// It provides the cluster data ["monitor", "services," <svcname>]
//
// worker ends when context is done or when no more service instance config exist
//
// worker watch on instance status updates to refresh object.AggregatedStatus
package svcagg

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"opensvc.com/opensvc/core/instance"
	"opensvc.com/opensvc/core/object"
	"opensvc.com/opensvc/core/path"
	"opensvc.com/opensvc/core/placement"
	"opensvc.com/opensvc/core/provisioned"
	"opensvc.com/opensvc/core/status"
	"opensvc.com/opensvc/daemon/daemondata"
	"opensvc.com/opensvc/daemon/msgbus"
	"opensvc.com/opensvc/util/pubsub"
)

type (
	svcAggStatus struct {
		status object.AggregatedStatus
		path   path.T
		id     string
		nodes  map[string]struct{}

		discoverCmdC chan<- any
		dataCmdC     chan<- any

		instStatus  map[string]instance.Status
		instMonitor map[string]instance.Monitor

		// srcEvent is the source event that triggered the svcAggStatus update
		srcEvent any

		ctx context.Context
		log zerolog.Logger

		subInstanceMonitorUpdated pubsub.Subscription
		subInstanceStatusUpdated  pubsub.Subscription
		subCfgUpdated             pubsub.Subscription
		subCfgDeleted             pubsub.Subscription
	}
)

// Start launch goroutine svcAggStatus worker for a service
func Start(ctx context.Context, p path.T, cfg instance.Config, svcAggDiscoverCmd chan<- any) error {
	id := p.String()
	o := &svcAggStatus{
		status:       object.AggregatedStatus{},
		path:         p,
		id:           id,
		discoverCmdC: svcAggDiscoverCmd,
		dataCmdC:     daemondata.BusFromContext(ctx),
		instStatus:   make(map[string]instance.Status),
		instMonitor:  make(map[string]instance.Monitor),
		ctx:          ctx,
		log:          log.Logger.With().Str("func", "svcagg").Stringer("object", p).Logger(),
	}
	o.startSubscriptions()

	go func() {
		defer o.stopSubscriptions()
		o.worker(cfg.Scope)
	}()
	return nil
}

func (o *svcAggStatus) stopSubscriptions() {
	o.subInstanceMonitorUpdated.Stop()
	o.subInstanceStatusUpdated.Stop()
	o.subCfgUpdated.Stop()
	o.subCfgDeleted.Stop()
}

func (o *svcAggStatus) startSubscriptions() {
	bus := pubsub.BusFromContext(o.ctx)
	label := pubsub.Label{"path", o.id}
	name := o.id + " svcagg"
	o.subInstanceMonitorUpdated = msgbus.Sub(bus, name, msgbus.InstanceMonitorUpdated{}, label)
	o.subInstanceStatusUpdated = msgbus.Sub(bus, name, msgbus.InstanceStatusUpdated{}, label)
	o.subCfgUpdated = msgbus.Sub(bus, name, msgbus.CfgUpdated{}, label)
	o.subCfgDeleted = msgbus.Sub(bus, name, msgbus.CfgDeleted{}, label)
}

func (o *svcAggStatus) worker(nodes []string) {
	o.log.Debug().Msg("started")
	defer o.log.Debug().Msg("done")

	for _, node := range nodes {
		o.instStatus[node] = daemondata.GetInstanceStatus(o.dataCmdC, o.path, node)
		o.instMonitor[node] = instance.Monitor{}
	}
	o.update()
	defer o.delete()
	for {
		if len(o.instStatus) == 0 {
			o.log.Info().Msg("no more nodes")
			return
		}
		select {
		case <-o.ctx.Done():
			return
		case i := <-o.subInstanceMonitorUpdated.C:
			c := i.(msgbus.InstanceMonitorUpdated)
			if _, ok := o.instMonitor[c.Node]; !ok {
				o.log.Debug().Msgf("skip instance monitor change from unknown node: %s", c.Node)
				continue
			}
			o.srcEvent = i
			o.instMonitor[c.Node] = c.Status
			o.updateStatus()
		case i := <-o.subInstanceStatusUpdated.C:
			c := i.(msgbus.InstanceStatusUpdated)
			if _, ok := o.instStatus[c.Node]; !ok {
				o.log.Debug().Msgf("skip instance status change from unknown node: %s", c.Node)
				continue
			}
			o.srcEvent = i
			o.instStatus[c.Node] = c.Status
			o.updateStatus()
		case i := <-o.subCfgUpdated.C:
			c := i.(msgbus.CfgUpdated)
			if _, ok := o.instStatus[c.Node]; ok {
				continue
			}
			o.srcEvent = i
			o.instStatus[c.Node] = daemondata.GetInstanceStatus(o.dataCmdC, o.path, c.Node)
			o.updateStatus()
		case i := <-o.subCfgDeleted.C:
			c := i.(msgbus.CfgDeleted)
			if _, ok := o.instStatus[c.Node]; ok {
				delete(o.instStatus, c.Node)
			}
			if _, ok := o.instMonitor[c.Node]; ok {
				delete(o.instMonitor, c.Node)
			}
			o.srcEvent = i
			o.updateStatus()
		}
	}
}

func (o *svcAggStatus) updateStatus() {
	// TODO update this simple aggregate status compute, perhaps already implemented
	updateAvail := func() {
		statusCount := make([]uint, 128, 128)
		var newAvail status.T
		for _, instStatus := range o.instStatus {
			statusCount[instStatus.Avail]++
		}
		if statusCount[status.Warn] > 0 {
			newAvail = status.Warn
		} else if statusCount[status.Up] > 0 {
			newAvail = status.Up
		} else if statusCount[status.Down] > 0 {
			newAvail = status.Down
		} else {
			newAvail = status.Undef
		}
		if o.status.Avail != newAvail {
			o.status.Avail = newAvail
		}
	}
	updateOverall := func() {
		if o.status.Avail == status.Warn {
			o.status.Overall = status.Warn
			return
		} else {
			o.status.Overall = status.Undef
		}
		for _, instStatus := range o.instStatus {
			o.status.Overall.Add(instStatus.Overall)
		}
	}
	updateProvisioned := func() {
		o.status.Provisioned = provisioned.Undef
		for _, instStatus := range o.instStatus {
			o.status.Provisioned = o.status.Provisioned.And(instStatus.Provisioned)
		}
	}
	updateFrozen := func() {
		m := map[bool]int{
			true:  0,
			false: 0,
		}
		for _, instStatus := range o.instStatus {
			m[instStatus.Frozen.IsZero()] += 1
		}
		n := len(o.instStatus)
		switch {
		case n == 0:
			o.status.Frozen = "n/a"
		case n == m[false]:
			o.status.Frozen = "frozen"
		case n == m[true]:
			o.status.Frozen = "thawed"
		default:
			o.status.Frozen = "mixed"
		}
	}
	updatePlacement := func() {
		o.status.Placement = placement.NotApplicable
		for node, instMonitor := range o.instMonitor {
			instStatus, ok := o.instStatus[node]
			if !ok {
				o.status.Placement = placement.NotApplicable
				break
			}
			if instMonitor.IsLeader && !instStatus.Avail.Is(status.Up, status.NotApplicable) {
				o.status.Placement = placement.NonOptimal
				break
			}
			if !instMonitor.IsLeader && !instStatus.Avail.Is(status.Down, status.NotApplicable) {
				o.status.Placement = placement.NonOptimal
				break
			}
			o.status.Placement = placement.Optimal
		}
	}

	updateAvail()
	updateOverall()
	updateProvisioned()
	updateFrozen()
	updatePlacement()
	o.update()
}

func (o *svcAggStatus) delete() {
	if err := daemondata.DelServiceAgg(o.dataCmdC, o.path); err != nil {
		o.log.Error().Err(err).Msg("DelServiceAgg")
	}
	o.discoverCmdC <- msgbus.ObjectAggDone{Path: o.path}
}

func (o *svcAggStatus) update() {
	value := o.status.DeepCopy()
	o.log.Debug().Msgf("update avail %s", value.Avail)
	if err := daemondata.SetServiceAgg(o.dataCmdC, o.path, *value, o.srcEvent); err != nil {
		o.log.Error().Err(err).Msg("SetServiceAgg")
	}
}
