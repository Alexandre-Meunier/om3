package cluster

import (
	"encoding/json"

	"opensvc.com/opensvc/core/objectselector"
	"opensvc.com/opensvc/core/path"
)

type (
	// Status describes the full Cluster state.
	Status struct {
		Cluster    TCluster                         `json:"cluster"`
		Collector  CollectorThreadStatus            `json:"collector"`
		DNS        DNSThreadStatus                  `json:"dns"`
		Scheduler  SchedulerThreadStatus            `json:"scheduler"`
		Listener   ListenerThreadStatus             `json:"listener"`
		Monitor    MonitorThreadStatus              `json:"monitor"`
		Heartbeats map[string]HeartbeatThreadStatus `json:"-"`
	}

	TCluster struct {
		Config ClusterConfig  `json:"config"`
		Status TClusterStatus `json:"status"`
	}

	TClusterStatus struct {
		Compat bool `json:"compat"`
		Frozen bool `json:"frozen"`
	}

	// ClusterConfig decribes the cluster id, name and nodes
	// The cluster name is used as the right most part of cluster dns
	// names.
	ClusterConfig struct {
		ID    string   `json:"id"`
		Name  string   `json:"name"`
		Nodes []string `json:"nodes"`
	}
)

func (s *Status) DeepCopy() *Status {
	b, err := json.Marshal(s)
	if err != nil {
		return nil
	}
	newStatus := Status{}
	if err := json.Unmarshal(b, &newStatus); err != nil {
		return nil
	}
	return &newStatus
}

// WithSelector purges the dataset from objects not matching the selector expression
func (s *Status) WithSelector(selector string) *Status {
	if selector == "" {
		return s
	}
	paths, err := objectselector.NewSelection(
		selector,
		objectselector.SelectionWithLocal(true),
	).Expand()
	if err != nil {
		return s
	}
	selected := paths.StrMap()
	for nodename, nodeData := range s.Monitor.Nodes {
		for ps, _ := range nodeData.Services.Config {
			if !selected.Has(ps) {
				delete(s.Monitor.Nodes[nodename].Services.Config, ps)
			}
		}
		for ps, _ := range nodeData.Services.Smon {
			if !selected.Has(ps) {
				delete(s.Monitor.Nodes[nodename].Services.Smon, ps)
			}
		}
		for ps, _ := range nodeData.Services.Status {
			if !selected.Has(ps) {
				delete(s.Monitor.Nodes[nodename].Services.Status, ps)
			}
		}
	}
	for ps, _ := range s.Monitor.Services {
		if !selected.Has(ps) {
			delete(s.Monitor.Services, ps)
		}
	}
	return s
}

// WithSelector purges the dataset from objects not matching the namespace
func (s *Status) WithNamespace(namespace string) *Status {
	if namespace == "" {
		return s
	}
	for nodename, nodeData := range s.Monitor.Nodes {
		for ps, _ := range nodeData.Services.Config {
			p, _ := path.Parse(ps)
			if p.Namespace != namespace {
				delete(s.Monitor.Nodes[nodename].Services.Config, ps)
			}
		}
		for ps, _ := range nodeData.Services.Smon {
			p, _ := path.Parse(ps)
			if p.Namespace != namespace {
				delete(s.Monitor.Nodes[nodename].Services.Smon, ps)
			}
		}
		for ps, _ := range nodeData.Services.Status {
			p, _ := path.Parse(ps)
			if p.Namespace != namespace {
				delete(s.Monitor.Nodes[nodename].Services.Status, ps)
			}
		}
	}
	for ps, _ := range s.Monitor.Services {
		p, _ := path.Parse(ps)
		if p.Namespace != namespace {
			delete(s.Monitor.Services, ps)
		}
	}
	return s
}
