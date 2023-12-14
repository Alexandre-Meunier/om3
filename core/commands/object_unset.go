package commands

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/opensvc/om3/core/actioncontext"
	"github.com/opensvc/om3/core/client"
	"github.com/opensvc/om3/core/naming"
	"github.com/opensvc/om3/core/object"
	"github.com/opensvc/om3/core/objectaction"
	"github.com/opensvc/om3/core/objectselector"
	"github.com/opensvc/om3/daemon/api"
	"github.com/opensvc/om3/util/key"
	"github.com/opensvc/om3/util/xsession"
)

type (
	CmdObjectUnset struct {
		OptsGlobal
		OptsLock
		Keywords []string
		Sections []string
	}
)

func (t *CmdObjectUnset) Run(selector, kind string) error {
	mergedSelector := mergeSelector(selector, t.ObjectSelector, kind, "")
	if t.Local || (t.NodeSelector != "") {
		return t.doObjectAction(mergedSelector)
	}
	c, err := client.New()
	if err != nil {
		return err
	}
	sel := objectselector.NewSelection(mergedSelector, objectselector.SelectionWithClient(c))
	paths, err := sel.Expand()
	if err != nil {
		return err
	}
	for _, p := range paths {
		params := api.PostObjectConfigUnsetParams{}
		params.Kw = &t.Keywords
		response, err := c.PostObjectConfigUnsetWithResponse(context.Background(), p.Namespace, p.Kind, p.Name, &params)
		if err != nil {
			return err
		}
		switch response.StatusCode() {
		case 200:
			return nil
		case 400:
			return fmt.Errorf("%s: %s", p, *response.JSON400)
		case 401:
			return fmt.Errorf("%s: %s", p, *response.JSON401)
		case 403:
			return fmt.Errorf("%s: %s", p, *response.JSON403)
		case 500:
			return fmt.Errorf("%s: %s", p, *response.JSON500)
		default:
			return fmt.Errorf("%s: unexpected response: %s", p, response.Status())
		}
	}
	return nil
}

func (t *CmdObjectUnset) doObjectAction(mergedSelector string) error {

	return objectaction.New(
		objectaction.LocalFirst(),
		objectaction.WithLocal(t.Local),
		objectaction.WithColor(t.Color),
		objectaction.WithOutput(t.Output),
		objectaction.WithObjectSelector(mergedSelector),
		objectaction.WithRemoteNodes(t.NodeSelector),
		objectaction.WithRemoteRun(func(ctx context.Context, p naming.Path, nodename string) (interface{}, error) {
			c, err := client.New(client.WithURL(nodename))
			if err != nil {
				return nil, err
			}
			params := api.PostInstanceActionUnsetParams{}
			if t.OptsLock.Disable {
				v := true
				params.NoLock = &v
			}
			if t.OptsLock.Timeout != 0 {
				v := fmt.Sprint(t.OptsLock.Timeout)
				params.WaitLock = &v
			}
			{
				sid := xsession.ID
				params.RequesterSid = &sid
				params.Kw = &t.Keywords
			}
			response, err := c.PostInstanceActionUnsetWithResponse(ctx, p.Namespace, p.Kind, p.Name, &params)
			if err != nil {
				return nil, err
			}
			switch {
			case response.JSON200 != nil:
				return *response.JSON200, nil
			case response.JSON401 != nil:
				return nil, fmt.Errorf("%s: node %s: %s", p, nodename, *response.JSON401)
			case response.JSON403 != nil:
				return nil, fmt.Errorf("%s: node %s: %s", p, nodename, *response.JSON403)
			case response.JSON500 != nil:
				return nil, fmt.Errorf("%s: node %s: %s", p, nodename, *response.JSON500)
			default:
				return nil, fmt.Errorf("%s: node %s: unexpected response: %s", p, nodename, response.Status())
			}
		}),

		objectaction.WithLocalRun(func(ctx context.Context, p naming.Path) (interface{}, error) {
			// TODO: one commit on Unset, one commit on DeleteSection. Change to single commit ?
			o, err := object.NewConfigurer(p)
			if err != nil {
				return nil, err
			}
			ctx = actioncontext.WithLockDisabled(ctx, t.Disable)
			ctx = actioncontext.WithLockTimeout(ctx, t.Timeout)
			kws := key.ParseStrings(t.Keywords)
			if len(kws) > 0 {
				log.Debug().Msgf("unsetting %s keywords: %s", p, kws)
				if err = o.Unset(ctx, kws...); err != nil {
					return nil, err
				}
			}
			sections := make([]string, 0)
			for _, r := range t.Sections {
				if r != "DEFAULT" {
					sections = append(sections, r)
				}
			}
			if len(sections) > 0 {
				log.Debug().Msgf("deleting %s sections: %s", p, sections)
				if err = o.DeleteSection(ctx, sections...); err != nil {
					return nil, err
				}
			}
			return nil, nil
		}),
	).Do()
}
