//go:build linux

package poolloop

import (
	"opensvc.com/opensvc/core/driver"
	"opensvc.com/opensvc/core/pool"
	"opensvc.com/opensvc/util/df"
	"opensvc.com/opensvc/util/sizeconv"
)

type (
	T struct {
		pool.T
	}
)

var (
	drvID = driver.NewID(driver.GroupPool, "loop")
)

func init() {
	driver.Register(drvID, NewPooler)
}

func NewPooler() pool.Pooler {
	t := New()
	var i interface{} = t
	return i.(pool.Pooler)
}

func New() *T {
	t := T{}
	return &t
}

func (t T) Head() string {
	return t.path()
}

func (t T) Capabilities() []string {
	return []string{"rox", "rwx", "roo", "rwo", "blk"}
}

func (t T) Usage() (pool.StatusUsage, error) {
	head := t.Head()
	entries, err := df.MountUsage(head)
	if err != nil {
		return pool.StatusUsage{}, err
	}
	if len(entries) == 0 {
		return pool.StatusUsage{}, err
	}
	e := entries[0]
	var used int64
	if size > 0 {
		size = e.Size / 1024
		free = e.Free / 1024
		used = e.Used / 1024
	} else {
		size = 0
		free = 0
	}
	usage := pool.StatusUsage{
		Size: float64(size),
		Free: float64(free),
		Used: float64(used),
	}
	return usage, nil
}

func (t *T) Translate(name string, size float64, shared bool) []string {
	data := t.BlkTranslate(name, size, shared)
	data = append(data, t.AddFS(name, shared, 1, 0, "disk#0")...)
	return data
}

func (t *T) BlkTranslate(name string, size float64, shared bool) []string {
	data := []string{
		"disk#0.type=loop",
		"disk#0.name=" + name,
		"disk#0.size=" + sizeconv.ExactBSizeCompact(size),
	}
	return data
}

func (t T) path() string {
	return t.GetString("path")
}
