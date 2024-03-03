//go:build windows

package installer

import (
	"errors"
	"syscall"

	"golang.org/x/sys/windows"
)

func crossDevice(err error) bool {
	var sysErr syscall.Errno
	if errors.As(err, &sysErr) {
		return err == syscall.EXDEV || err == windows.ERROR_NOT_SAME_DEVICE
	}
	return false
}
