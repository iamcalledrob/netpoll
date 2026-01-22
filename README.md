# Netpoll

Adaptation of https://github.com/mailru/easygo/ as a Go module and patched
to use `syscall.Conn` instead of `os.File`, to avoid blocking/deadline issues
associated with using a File to access the Fd.

Minor breaking API changes, methods explicitly take a `syscall.Conn`.