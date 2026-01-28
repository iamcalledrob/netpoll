package netpoll

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

// Desc is a network connection within netpoll descriptor.
// It's methods are not goroutine safe.
type Desc struct {
	fd    int
	event Event
}

// withFd invokes f with a guaranteed valid fd, propagating errors to the caller
// callback is used for forwards-compatibility with SyscallConn.Control() if needed in the future
func (h *Desc) withFd(f func(fd int) error) error {
	return f(h.fd)
}

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

// HandleRead creates read descriptor for further use in Poller methods.
// It is the same as Handle(conn, EventRead|EventEdgeTriggered).
func HandleRead(conn syscall.Conn) (*Desc, error) {
	return Handle(conn, EventRead|EventEdgeTriggered)
}

// HandleReadOnce creates read descriptor for further use in Poller methods.
// It is the same as Handle(conn, EventRead|EventOneShot).
func HandleReadOnce(conn syscall.Conn) (*Desc, error) {
	return Handle(conn, EventRead|EventOneShot)
}

// HandleWrite creates write descriptor for further use in Poller methods.
// It is the same as Handle(conn, EventWrite|EventEdgeTriggered).
func HandleWrite(conn syscall.Conn) (*Desc, error) {
	return Handle(conn, EventWrite|EventEdgeTriggered)
}

// HandleWriteOnce creates write descriptor for further use in Poller methods.
// It is the same as Handle(conn, EventWrite|EventOneShot).
func HandleWriteOnce(conn syscall.Conn) (*Desc, error) {
	return Handle(conn, EventWrite|EventOneShot)
}

// HandleReadWrite creates read and write descriptor for further use in Poller
// methods.
// It is the same as Handle(conn, EventRead|EventWrite|EventEdgeTriggered).
func HandleReadWrite(conn syscall.Conn) (*Desc, error) {
	return Handle(conn, EventRead|EventWrite|EventEdgeTriggered)
}
