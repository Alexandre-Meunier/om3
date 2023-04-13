package daemondata

import (
	"sort"

	"github.com/opensvc/om3/core/cluster"
	"github.com/opensvc/om3/daemon/hbcache"
	"github.com/opensvc/om3/daemon/msgbus"
	"github.com/opensvc/om3/util/pubsub"
	"github.com/opensvc/om3/util/stringslice"
)

func (d *data) setDaemonHb() {
	hbModes := make([]cluster.HbMode, 0)
	nodes := make([]string, 0)
	for node := range d.hbMsgMode {
		if !stringslice.Has(node, d.clusterData.Cluster.Config.Nodes) {
			// Drop not anymore in cluster config nodes
			hbcache.DropPeer(node)
			continue
		}
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)
	for _, node := range nodes {
		hbModes = append(hbModes, cluster.HbMode{
			Node: node,
			Mode: d.hbMsgMode[node],
			Type: d.hbMsgType[node],
		})
	}

	subHb := cluster.DaemonHb{
		Streams: hbcache.Heartbeats(),
		Modes:   hbModes,
	}
	d.clusterData.Daemon.Hb = subHb
	d.bus.Pub(&msgbus.DaemonHb{Node: d.localNode, Value: subHb}, d.labelLocalNode)
}

func (d *data) setHbMsgMode(node string, mode string) {
	d.hbMsgMode[node] = mode
}

// setHbMsgType update the sub.hb.mode.x.Type for node,
// if value is changed publish msgbus.HbMessageTypeUpdated
func (d *data) setHbMsgType(node string, msgType string) {
	previous := d.hbMsgType[node]
	if msgType != previous {
		d.hbMsgType[node] = msgType
		joinedNodes := make([]string, 0)
		for n, v := range d.hbMsgType {
			if v == "patch" {
				joinedNodes = append(joinedNodes, n)
			}
		}
		d.bus.Pub(&msgbus.HbMessageTypeUpdated{
			Node:          node,
			From:          previous,
			To:            msgType,
			Nodes:         append([]string{}, d.clusterData.Cluster.Config.Nodes...),
			JoinedNodes:   joinedNodes,
			InstalledGens: d.deepCopyLocalGens(),
		}, pubsub.Label{"node", node})
	}
}
