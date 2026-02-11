package netpoll

import "syscall"

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
