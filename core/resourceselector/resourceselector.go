package resourceselector

import (
	"context"
	"strings"

	"opensvc.com/opensvc/core/actioncontext"
	"opensvc.com/opensvc/core/actionresdeps"
	"opensvc.com/opensvc/core/ordering"
	"opensvc.com/opensvc/core/resource"
	"opensvc.com/opensvc/util/funcopt"
)

type (
	// T contains options accepted by all actions manipulating resources
	T struct {
		rid    string
		subset string
		tag    string
		order  ordering.T
		lister ResourceLister
		action string
	}

	// ResourceLister is the interface required to list resource.T and see the ordering
	ResourceLister interface {
		Resources() resource.Drivers
		ReconfigureResource(r resource.Driver) error
		IsDesc() bool
	}

	rider interface {
		ResourceSelectorRID() string
	}
	tager interface {
		ResourceSelectorTag() string
	}
	subseter interface {
		ResourceSelectorSubset() string
	}

	depser interface {
		GetActionResDeps() *actionresdeps.Store
	}
)

func WithRID(s string) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.rid = s
		return nil
	})
}

func WithTag(s string) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.tag = s
		return nil
	})
}

func WithSubset(s string) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.subset = s
		return nil
	})
}

func WithOrder(s ordering.T) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.order = s
		return nil
	})
}

func WithAction(s string) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.action = s
		return nil
	})
}

func New(l ResourceLister, opts ...funcopt.O) *T {
	t := &T{
		lister: l,
	}
	_ = funcopt.Apply(t, opts...)
	return t
}

func (t *T) SetResourceLister(l ResourceLister) {
	t.lister = l
}

func (t T) IsDesc() bool {
	return t.order.IsDesc()
}

func (t T) ReconfigureResource(r resource.Driver) error {
	return t.lister.ReconfigureResource(r)
}

func (t T) Resources() resource.Drivers {
	l := t.lister.Resources()
	if t.rid == "" && t.tag == "" && t.subset == "" {
		return l
	}
	var dp *actionresdeps.Store
	if i, ok := t.lister.(depser); ok {
		dp = i.GetActionResDeps()
	}
	fl := make(resource.Drivers, 0)
	f := func(c rune) bool { return c == ',' }
	rids := strings.FieldsFunc(t.rid, f)
	tags := strings.FieldsFunc(t.tag, f)
	subsets := strings.FieldsFunc(t.subset, f)
	for _, r := range l {
		for _, e := range rids {
			if r.MatchRID(e) {
				goto add
			}
		}
		for _, e := range subsets {
			if r.MatchSubset(e) {
				goto add
			}
		}
		for _, e := range tags {
			if r.MatchTag(e) {
				goto add
			}
		}
		continue
	add:
		fl = fl.Add(r)
		if dp != nil {
			deps := dp.SelectDependencies(t.action, r.RID())
			for _, rid := range deps {
				if dep := l.GetRID(rid); dep != nil {
					r.Log().Debug().Msgf("add %s to satisfy %s dependency", dep.RID(), t.action)
					fl = fl.Add(dep)
				}
			}
		}
	}
	if t.order == ordering.Desc {
		fl.Reverse()
	} else {
		fl.Sort()
	}
	return fl
}

func (t T) IsZero() bool {
	switch {
	case t.rid != "":
		return false
	case t.subset != "":
		return false
	case t.tag != "":
		return false
	default:
		return true
	}
}

func FromContext(ctx context.Context, l ResourceLister) *T {
	props := actioncontext.Props(ctx)
	return New(
		l,
		WithRID(RIDFromContext(ctx)),
		WithTag(TagFromContext(ctx)),
		WithSubset(SubsetFromContext(ctx)),
		WithOrder(props.Order),
		WithAction(props.Name),
	)
}

func RIDFromContext(ctx context.Context) string {
	if o, ok := actioncontext.Value(ctx).Options.(rider); ok {
		return o.ResourceSelectorRID()
	}
	return ""
}

func TagFromContext(ctx context.Context) string {
	if o, ok := actioncontext.Value(ctx).Options.(tager); ok {
		return o.ResourceSelectorTag()
	}
	return ""
}

func SubsetFromContext(ctx context.Context) string {
	if o, ok := actioncontext.Value(ctx).Options.(subseter); ok {
		return o.ResourceSelectorSubset()
	}
	return ""
}
