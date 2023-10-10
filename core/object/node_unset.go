package object

import (
	"context"

	"github.com/opensvc/om3/util/key"
)

// Unset removes keywords from node config
func (t *Node) Unset(ctx context.Context, kws ...key.T) error {
	return unsetKeys(t.config, kws...)
}
