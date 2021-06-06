package mudp

import (
	"net"

	"github.com/gogf/gf/net/gudp"
)

// NewConn 根据 <address> 创建 UDP 连接。
// 注意：该 *gudp.Conn 只能发送数据，无法接受数据
func NewConn(address string) (*gudp.Conn, error) {
	host, port, err := getHostAndPortFromAddress(address)
	if err != nil {
		return nil, err
	}
	ip := net.ParseIP(host)
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: port}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		return nil, err
	}
	return gudp.NewConnByNetConn(conn), nil
}
