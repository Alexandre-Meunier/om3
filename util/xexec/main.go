package xexec

import (
	"github.com/anmitsu/go-shlex"
	"opensvc.com/opensvc/util/limits"
	"os/exec"
	"strings"
)

// CommandFromString wrapper to exec.Command from a string command 's'
// When string command 's' contains multiple commands,
//   exec.Command("/bin/sh", "-c", s)
// else
//   exec.Command from shlex.Split(s)
func CommandFromString(s string) (*exec.Cmd, error) {
	var needShell bool
	switch {
	case strings.Contains(s, "|"):
		needShell = true
	case strings.Contains(s, "&&"):
		needShell = true
	case strings.Contains(s, ";"):
		needShell = true
	}
	if needShell {
		return exec.Command("/bin/sh", "-c", s), nil
	}
	sSplit, err := shlex.Split(s, true)
	if err != nil {
		return nil, err
	}
	if len(sSplit) == 0 {
		return nil, nil
	}
	return exec.Command(sSplit[0], sSplit[1:]...), nil
}

// CommandFromLimits wrapper to exec.Command from a string command 's'
// and 'l' limits.T
func CommandFromLimits(l limits.T, s string) (*exec.Cmd, error) {
	limitCommands := shLimitCommands(l)
	if len(limitCommands) > 0 {
		s = strings.Join(limitCommands, " && ") + " && " + s
	}
	return CommandFromString(s)
}
