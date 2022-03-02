/*
	Package httpmux provides http mux

	It defines routing for Opensvc listener daemons

*/
package httpmux

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"opensvc.com/opensvc/daemon/listener/handlers/daemonhandler"
	"opensvc.com/opensvc/daemon/listener/handlers/dispatchhandler"
	"opensvc.com/opensvc/daemon/listener/mux/muxctx"
	"opensvc.com/opensvc/daemon/subdaemon"
)

type (
	T struct {
		log        zerolog.Logger
		mux        *chi.Mux
		rootDaemon subdaemon.RootManager
	}
)

// New returns *T with log, rootDaemon
// it prepares middlewares and routes for Opensvc daemon listeners
func New(log zerolog.Logger, rootDaemon subdaemon.RootManager) *T {
	t := &T{
		log:        log,
		rootDaemon: rootDaemon,
	}
	mux := chi.NewRouter()
	mux.Use(daemonMiddleWare(t.rootDaemon))
	mux.Use(logMiddleWare(t.log))
	mux.Post("/daemon_stop", daemonhandler.Stop)
	mux.Mount("/daemon", t.newDaemonRouter())

	t.mux = mux
	return t
}

// ServerHTTP implement http.Handler interface for T
func (t *T) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.mux.ServeHTTP(w, r)
}

func (t *T) newDaemonRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/running", dispatchhandler.New(daemonhandler.Running))
	r.Post("/stop", daemonhandler.Stop)
	r.Get("/eventsdemo", daemonhandler.Events)
	return r
}

func logMiddleWare(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uuid := uuid.New()
			ctx := muxctx.WithLogger(r.Context(), logger.With().Str("request-uuid", uuid.String()).Logger())
			logger.Info().Str("METHOD", r.Method).Str("PATH", r.URL.Path).Msg("request")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func daemonMiddleWare(manager subdaemon.RootManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := muxctx.WithDaemon(r.Context(), manager)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
