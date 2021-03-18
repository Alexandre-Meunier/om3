package object

import (
	"opensvc.com/opensvc/core/provisioned"
	"opensvc.com/opensvc/core/status"
)

type (
	// InstanceMonitor describes the in-daemon states of an instance
	InstanceMonitor struct {
		GlobalExpect        string  `json:"global_expect"`
		LocalExpect         string  `json:"local_expect"`
		Status              string  `json:"status"`
		StatusUpdated       float64 `json:"status_updated"`
		GlobalExpectUpdated float64 `json:"global_expect_updated"`
		Placement           string  `json:"placement"`
	}

	// InstanceConfig describes a configuration file content checksum,
	// timestamp of last change and the nodes it should be installed on.
	InstanceConfig struct {
		Nodename string   `json:"-"`
		Path     Path     `json:"-"`
		Checksum string   `json:"csum"`
		Scope    []string `json:"scope"`
		Updated  float64
	}

	// InstanceStatus describes the instance status.
	InstanceStatus struct {
		Nodename    string                    `json:"-"`
		Path        Path                      `json:"-"`
		App         string                    `json:"app,omitempty"`
		Avail       status.T               `json:"avail,omitempty"`
		DRP         bool                      `json:"drp,omitempty"`
		Overall     status.T               `json:"overall,omitempty"`
		Csum        string                    `json:"csum,omitempty"`
		Env         string                    `json:"env,omitempty"`
		Frozen      float64                   `json:"frozen,omitempty"`
		Kind        Kind                      `json:"kind"`
		Monitor     InstanceMonitor           `json:"monitor"`
		Optional    status.T               `json:"optional,omitempty"`
		Orchestrate string                    `json:"orchestrate,omitempty"` // TODO enum
		Topology    string                    `json:"topology,omitempty"`    // TODO enum
		Placement   string                    `json:"placement,omitempty"`   // TODO enum
		Priority    int                       `json:"priority,omitempty"`
		Provisioned provisioned.T          `json:"provisioned,omitempty"`
		Preserved   bool                      `json:"preserved,omitempty"`
		Updated     float64                   `json:"updated"`
		FlexTarget  int                       `json:"flex_target,omitempty"`
		FlexMin     int                       `json:"flex_min,omitempty"`
		FlexMax     int                       `json:"flex_max,omitempty"`
		Subsets     map[string]SubsetStatus   `json:"subsets,omitempty"`
		Resources   map[string]ResourceStatus `json:"resources,omitempty"`
		Running     []string                  `json:"running,omitempty"`
	}

	// SubsetStatus describes a resource subset properties.
	SubsetStatus struct {
		Parallel bool `json:"parallel,omitempty"`
	}

	// ResourceStatus describes the status of a resource of an instance of an object.
	ResourceStatus struct {
		Label       string                  `json:"label"`
		Log         []string                `json:"log"`
		Status      status.T             `json:"status"`
		Type        string                  `json:"type"`
		Provisioned ResourceStatusProvision `json:"provisioned"`
	}

	// ResourceStatusProvision define if and when the resource became provisioned.
	ResourceStatusProvision struct {
		Mtime float64          `json:"mtime"`
		State provisioned.T `json:"state"`
	}
)
