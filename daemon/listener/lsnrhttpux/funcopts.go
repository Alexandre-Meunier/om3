package lsnrhttpux

import (
	"context"
	"net/http"

	"opensvc.com/opensvc/daemon/routinehelper"
	"opensvc.com/opensvc/util/funcopt"
)

func WithContext(parent context.Context) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.Ctx, t.CancelFunc = context.WithCancel(parent)
		return nil
	})
}

func WithRoutineTracer(o routinehelper.Tracer) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.SetTracer(o)
		return nil
	})
}

func WithHandler(o http.Handler) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.handler = o
		return nil
	})
}

func WithAddr(o string) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.addr = o
		return nil
	})
}

func WithCertFile(o string) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.certFile = o
		return nil
	})
}

func WithKeyFile(o string) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.keyFile = o
		return nil
	})
}
