//go:build windows

package installer

import (
	"errors"
	"syscall"
)

func crossDevice(err error) bool {
	var sysErr syscall.Errno
	if errors.As(err, &sysErr) {
		return err == syscall.EXDEV || err == syscall.ERROR_NOT_SAME_DEVICE
	}
	return false
}
