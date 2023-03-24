//go:build linux

package raw

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/opensvc/om3/util/file"
)

func TestRaw(t *testing.T) {
	rawdev := "/dev/raw"
	if v, err := file.ExistsAndDir(rawdev); err != nil {
		t.Fatalf("%s", err)
	} else if !v {
		t.Skipf("no %s, skip test", rawdev)
	}
	if os.Getuid() != 0 {
		t.Skip("skipped for non root user")
	}
	log := &zerolog.Logger{}
	ra := New(
		WithLogger(log),
	)
	t.Logf("data")
	if os.Getuid() != 0 {
		t.Skip("skipped for non root user")
	}
	data, err := ra.QueryAll()
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(data), 0)
	//
	// BEWARE: uncomment only after setting a proper bdevpath for your env
	//
	//minor, err := ra.Bind("/dev/sda")
	//assert.Nil(t, err)
	//assert.GreaterOrEqual(t, minor, 1)
	//err = ra.Unbind(minor)
	//assert.Nil(t, err)
}
