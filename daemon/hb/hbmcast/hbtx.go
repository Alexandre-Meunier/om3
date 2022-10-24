package hbmcast

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	reqjsonrpc "opensvc.com/opensvc/core/client/requester/jsonrpc"
	"opensvc.com/opensvc/core/hbtype"
	"opensvc.com/opensvc/core/rawconfig"
	"opensvc.com/opensvc/daemon/daemonlogctx"
	"opensvc.com/opensvc/daemon/hb/hbctrl"
	"opensvc.com/opensvc/util/hostname"
)

type (
	tx struct {
		ctx     context.Context
		id      string
		nodes   []string
		udpAddr *net.UDPAddr
		intf    string
		timeout time.Duration

		name   string
		log    zerolog.Logger
		cmdC   chan<- interface{}
		msgC   chan<- *hbtype.Msg
		cancel func()
	}
)

// Id implements the Id function of Transmitter interface for tx
func (t *tx) Id() string {
	return t.id
}

// Stop implements the Stop function of Transmitter interface for tx
func (t *tx) Stop() error {
	t.cancel()
	for _, node := range t.nodes {
		t.cmdC <- hbctrl.CmdDelWatcher{
			HbId:     t.id,
			Nodename: node,
		}
	}
	return nil
}

// Start implements the Start function of Transmitter interface for tx
func (t *tx) Start(cmdC chan<- interface{}, msgC <-chan []byte) error {
	started := make(chan bool)
	ctx, cancel := context.WithCancel(t.ctx)
	t.cancel = cancel
	t.cmdC = cmdC

	go func() {
		t.log.Info().Msg("starting")
		for _, node := range t.nodes {
			cmdC <- hbctrl.CmdAddWatcher{
				HbId:     t.id,
				Nodename: node,
				Ctx:      ctx,
				Timeout:  t.timeout,
			}
		}
		started <- true
		for {
			select {
			case <-ctx.Done():
				t.log.Info().Msg("stopped")
				return
			case b := <-msgC:
				go t.send(b)
			}
		}
	}()
	<-started
	t.log.Info().Msg("started")
	return nil
}

func (t *tx) encryptMessage(b []byte) ([]byte, error) {
	cluster := rawconfig.ClusterSection()
	msg := &reqjsonrpc.Message{
		NodeName:    hostname.Hostname(),
		ClusterName: cluster.Name,
		Key:         cluster.Secret,
		Data:        b,
	}
	return msg.Encrypt()
}

func (t *tx) send(b []byte) {
	//fmt.Println("xx >>>\n", hex.Dump(b))
	t.log.Debug().Msgf("send to udp %s", t.udpAddr)
	encMsg, err := t.encryptMessage(b)
	if err != nil {
		t.log.Debug().Err(err).Msg("encrypt")
		return
	}
	c, err := net.DialUDP("udp", nil, t.udpAddr)
	if err != nil {
		t.log.Debug().Err(err).Msgf("dial udp %s", t.udpAddr)
		return
	}
	defer c.Close()
	msgID := uuid.New().String()
	msgLength := len(encMsg)
	total := msgLength / MaxDatagramSize
	if (msgLength % MaxDatagramSize) != 0 {
		total += 1
	}
	for i := 1; i <= total; i += 1 {
		f := fragment{
			MsgID: msgID,
			Index: i,
			Total: total,
		}
		if i == total {
			f.Chunk = encMsg
		} else {
			f.Chunk = encMsg[:MaxDatagramSize]
			encMsg = encMsg[MaxDatagramSize:]
		}
		dgram, err := json.Marshal(f)
		if err != nil {
			t.log.Debug().Err(err).Msgf("marshal frame")
			return
		}
		if _, err := c.Write(dgram); err != nil {
			t.log.Debug().Err(err).Msgf("write in udp conn to %s", t.udpAddr)
			return
		}
	}
	for _, node := range t.nodes {
		t.cmdC <- hbctrl.CmdSetPeerSuccess{
			Nodename: node,
			HbId:     t.id,
			Success:  true,
		}
	}
}

func newTx(ctx context.Context, name string, nodes []string, udpAddr *net.UDPAddr, intf string, timeout time.Duration) *tx {
	id := name + ".tx"
	log := daemonlogctx.Logger(ctx).With().Str("id", id).Logger()
	return &tx{
		ctx:     ctx,
		id:      id,
		nodes:   nodes,
		udpAddr: udpAddr,
		intf:    intf,
		timeout: timeout,
		log:     log,
	}
}
