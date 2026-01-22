package netpoll

import (
	"fmt"
	"syscall"
)

// Desc is a network connection within netpoll descriptor.
// It's methods are not goroutine safe.
type Desc struct {
	conn  syscall.Conn // reference to the conn, not its fd, to avoid garbage collection during netpoll ops
	event Event
}

// withFd invokes f with a guaranteed valid fd, propagating errors to the caller
func (h *Desc) withFd(f func(fd uintptr) error) error {
	rawConn, err := h.conn.SyscallConn()
	if err != nil {
		return fmt.Errorf("SyscallConn: %w", err)
	}

	var invokeErr error
	err = rawConn.Control(func(fd uintptr) {
		invokeErr = f(fd)
	})
	if err != nil {
		return fmt.Errorf("rawConn: getting file descriptor: %w", err)
	}
	if invokeErr != nil {
		return fmt.Errorf("rawConn: invoking f: %w", invokeErr)
	}

	return nil
}

// HandleRead creates read descriptor for further use in Poller methods.
// It is the same as Handle(conn, EventRead|EventEdgeTriggered).
func HandleRead(conn syscall.Conn) *Desc {
	return Handle(conn, EventRead|EventEdgeTriggered)
}

// HandleReadOnce creates read descriptor for further use in Poller methods.
// It is the same as Handle(conn, EventRead|EventOneShot).
func HandleReadOnce(conn syscall.Conn) *Desc {
	return Handle(conn, EventRead|EventOneShot)
}

// HandleWrite creates write descriptor for further use in Poller methods.
// It is the same as Handle(conn, EventWrite|EventEdgeTriggered).
func HandleWrite(conn syscall.Conn) *Desc {
	return Handle(conn, EventWrite|EventEdgeTriggered)
}

// HandleWriteOnce creates write descriptor for further use in Poller methods.
// It is the same as Handle(conn, EventWrite|EventOneShot).
func HandleWriteOnce(conn syscall.Conn) *Desc {
	return Handle(conn, EventWrite|EventOneShot)
}

// HandleReadWrite creates read and write descriptor for further use in Poller
// methods.
// It is the same as Handle(conn, EventRead|EventWrite|EventEdgeTriggered).
func HandleReadWrite(conn syscall.Conn) *Desc {
	return Handle(conn, EventRead|EventWrite|EventEdgeTriggered)
}

// Handle creates new Desc with given conn and event.
// Returned descriptor could be used as argument to Start(), Resume() and
// Stop() methods of some Poller implementation.
func Handle(conn syscall.Conn, event Event) *Desc {
	return &Desc{
		conn:  conn,
		event: event,
	}
}
