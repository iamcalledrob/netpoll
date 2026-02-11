//go:build linux || darwin || dragonfly || freebsd || netbsd || openbsd

package netpoll

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

func (h *Desc) Close() error {
	return unix.Close(h.fd)
}

// Handle creates new Desc with given conn and event.
// Returned descriptor could be used as argument to Start(), Resume() and
// Stop() methods of some Poller implementation.
// Caller must call Close() to avoid leaks.
func Handle(conn syscall.Conn, event Event) (*Desc, error) {
	sc, err := conn.SyscallConn()
	if err != nil {
		return nil, fmt.Errorf("SyscallConn: %w", err)
	}

	// Duplicate the fd to avoid races with the Go runtime,
	// which is liable to close/allow recycling of FDs when
	// conn failure is detected.
	//
	// By using dup (and not File()), the original socket fd
	// remains intact and doesn't enter blocking mode
	var dupFd int
	var dupErr error
	err = sc.Control(func(fd uintptr) {
		dupFd, dupErr = unix.Dup(int(fd))
	})
	if err != nil {
		return nil, fmt.Errorf("obtaining fd: %w", err)
	}

	if dupErr != nil {
		return nil, fmt.Errorf("duplicating fd: %w", dupErr)
	}

	return &Desc{
		fd:    dupFd,
		event: event,
	}, nil
}
