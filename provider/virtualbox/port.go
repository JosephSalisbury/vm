package virtualbox

import (
	"net"
)

const (
	localhostNet = "localhost:0"
	tcp          = "tcp"
)

// getFreePort returns a port number that can be used.
func getFreePort() int {
	addr, err := net.ResolveTCPAddr(tcp, localhostNet)
	if err != nil {
		return 0
	}

	l, err := net.ListenTCP(tcp, addr)
	if err != nil {
		return 0
	}
	defer l.Close()

	port := l.Addr().(*net.TCPAddr).Port

	return port
}
