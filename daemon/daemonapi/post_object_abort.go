package daemonapi

import (
	"encoding/json"
	"net/http"

	"github.com/opensvc/om3/core/instance"
	"github.com/opensvc/om3/core/path"
	"github.com/opensvc/om3/daemon/msgbus"
	"github.com/opensvc/om3/util/hostname"
	"github.com/opensvc/om3/util/pubsub"
)

func (a *DaemonApi) PostObjectAbort(w http.ResponseWriter, r *http.Request) {
	var (
		payload = PostObjectAbort{}
		p       path.T
		err     error
	)
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	p, err = path.Parse(payload.Path)
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid path: "+payload.Path)
		return
	}
	globalExpect := instance.MonitorGlobalExpectAborted
	instMonitor := instance.MonitorUpdate{
		GlobalExpect: &globalExpect,
	}
	bus := pubsub.BusFromContext(r.Context())
	bus.Pub(&msgbus.SetInstanceMonitor{Path: p, Node: hostname.Hostname(), Value: instMonitor},
		pubsub.Label{"path", p.String()}, labelApi)
	w.WriteHeader(http.StatusOK)
}
