package v

import (
	"fmt"
	"net"
)

func getAddrString(conn net.Conn) string {
	localAddr := conn.LocalAddr().String()
	remoteAddr := conn.RemoteAddr().String()
	return fmt.Sprintf("%s (%s)", localAddr, remoteAddr)
}
