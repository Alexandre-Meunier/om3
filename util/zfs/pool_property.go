package zfs

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/opensvc/om3/util/command"
)

// GetProperty returns a dataset property value
func (t *Pool) GetProperty(prop string) (string, error) {
	cmd := command.New(
		command.WithName("zpool"),
		command.WithVarArgs("get", "-Hp", "-o", "value", prop, t.Name),
		command.WithBufferedStdout(),
		command.WithLogger(t.Log),
		command.WithCommandLogLevel(zerolog.DebugLevel),
		command.WithStdoutLogLevel(zerolog.DebugLevel),
		command.WithStderrLogLevel(zerolog.DebugLevel),
	)
	b, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func (t *Pool) SetProperty(prop, value string) error {
	s := fmt.Sprintf("%s=%s", prop, value)
	cmd := command.New(
		command.WithName("zpool"),
		command.WithVarArgs("set", s, t.Name),
		command.WithLogger(t.Log),
		command.WithCommandLogLevel(zerolog.InfoLevel),
		command.WithStdoutLogLevel(zerolog.InfoLevel),
		command.WithStderrLogLevel(zerolog.ErrorLevel),
	)
	return cmd.Run()
}
