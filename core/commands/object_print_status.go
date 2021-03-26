package commands

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"opensvc.com/opensvc/config"
	"opensvc.com/opensvc/core/client"
	"opensvc.com/opensvc/core/cluster"
	"opensvc.com/opensvc/core/object"
	"opensvc.com/opensvc/core/output"
)

type (
	// CmdObjectPrintStatus is the cobra flag set of the status command.
	CmdObjectPrintStatus struct {
		object.OptsStatus
	}
)

// Init configures a cobra command and adds it to the parent command.
func (t *CmdObjectPrintStatus) Init(kind string, parent *cobra.Command, selector *string) {
	cmd := t.cmd(kind, selector)
	parent.AddCommand(cmd)
	object.InstallFlags(cmd, t)
}

func (t *CmdObjectPrintStatus) cmd(kind string, selector *string) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Print selected service and instance status",
		Long: `Resources Flags:

(1) R   Running,           . Not Running
(2) M   Monitored,         . Not Monitored
(3) D   Disabled,          . Enabled
(4) O   Optional,          . Not Optional
(5) E   Encap,             . Not Encap
(6) P   Not Provisioned,   . Provisioned
(7) S   Standby,           . Not Standby
(8) <n> Remaining Restart, + if more than 10,   . No Restart

`,
		Run: func(cmd *cobra.Command, args []string) {
			t.run(selector, kind)
		},
	}
}

func (t *CmdObjectPrintStatus) extract(selector string, c *client.T) []object.Status {
	if data, err := t.extractFromDaemon(selector, c); err == nil {
		log.Debug().Err(err).Msg("extract cluster status")
		return data
	}
	if client.WantContext() {
		log.Error().Msg("can not fetch daemon data")
		return []object.Status{}
	}
	return t.extractLocal(selector)
}

func (t *CmdObjectPrintStatus) extractLocal(selector string) []object.Status {
	data := make([]object.Status, 0)
	sel := object.NewSelection(selector).SetLocal(true)
	for _, path := range sel.Expand() {
		obj := path.NewBaser()
		status, err := obj.Status(t.OptsStatus)
		if err != nil {
			log.Debug().Err(err).Str("path", path.String()).Msg("extract local")
			continue
		}
		o := object.Status{
			Path:   path,
			Compat: true,
			Object: object.AggregatedStatus{},
			Instances: map[string]object.InstanceStates{
				config.Node.Hostname: {
					Node: object.InstanceNode{
						Name:   config.Node.Hostname,
						Frozen: object.NewNode().Frozen(),
					},
					Status: status,
				},
			},
		}
		data = append(data, o)
	}
	return data
}

func (t *CmdObjectPrintStatus) extractFromDaemon(selector string, c *client.T) ([]object.Status, error) {
	var (
		err           error
		b             []byte
		clusterStatus cluster.Status
	)
	getter, err := client.NewGetDaemonStatusB(c, client.WithSelector(selector), client.WithRelatives(true))
	if err != nil {
		return []object.Status{}, err
	}
	b, err = getter.Get()
	if err != nil {
		return []object.Status{}, err
	}
	err = json.Unmarshal(b, &clusterStatus)
	if err != nil {
		return []object.Status{}, err
	}
	data := make([]object.Status, 0)
	for p := range clusterStatus.Monitor.Services {
		path, err := object.NewPathFromString(p)
		if err != nil {
			log.Debug().Err(err).Str("path", p).Msg("extractFromDaemon")
			continue
		}
		data = append(data, clusterStatus.GetObjectStatus(path))
	}
	return data, nil
}

func (t *CmdObjectPrintStatus) run(selector *string, kind string) {
	var data []object.Status
	mergedSelector := mergeSelector(*selector, t.Global.ObjectSelector, kind, "")
	c, err := client.New(client.URL(t.Global.Server))
	if err == nil {
		data = t.extract(mergedSelector, c)
	}
	output.Renderer{
		Format: t.Global.Format,
		Color:  t.Global.Color,
		Data:   data,
		HumanRenderer: func() string {
			s := ""
			for _, d := range data {
				s += d.Render()
			}
			return s
		},
	}.Print()
}
