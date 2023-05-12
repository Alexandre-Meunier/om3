package daemonapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/opensvc/om3/daemon/api"
	"github.com/opensvc/om3/daemon/daemonlogctx"
	"github.com/opensvc/om3/daemon/relay"
)

func (a *DaemonApi) PostRelayMessage(w http.ResponseWriter, r *http.Request) {
	var (
		payload api.PostRelayMessage
		value   api.RelayMessage
	)
	log := daemonlogctx.Logger(r.Context()).With().Str("func", "PostRelayMessage").Logger()
	log.Debug().Msg("starting")

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteProblemf(w, http.StatusBadRequest, "Invalid body", "%s", err)
		return
	}

	value.ClusterName = payload.ClusterName
	value.ClusterId = payload.ClusterId
	value.Nodename = payload.Nodename
	value.Msg = payload.Msg
	value.Updated = time.Now()
	value.Addr = r.RemoteAddr

	relay.Map.Store(payload.ClusterId, payload.Nodename, value)
	log.Debug().Msgf("stored %s %s", payload.ClusterId, payload.Nodename)
	WriteProblemf(w, http.StatusOK, "stored", "at %s from %s", value.Updated, value.Addr)
}
