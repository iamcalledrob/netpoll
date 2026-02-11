//go:build !linux && !darwin && !dragonfly && !freebsd && !netbsd && !openbsd

package netpoll

import (
	"fmt"
	"runtime"
	"syscall"
)

func (h *Desc) Close() error {
	return fmt.Errorf("unsupported on %s", runtime.GOOS)
}

func Handle(syscall.Conn, Event) (*Desc, error) {
	return nil, fmt.Errorf("unsupported on %s", runtime.GOOS)
}
