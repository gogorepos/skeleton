package mudp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/gudp"
)

// NewMulticastConn 根据 <address> 创建并返回一个组播连接 *net.UDPConn
func NewMulticastConn(address string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	/**
	 * net.ListenMulticastUDP 第二个参数官方文档中说如果使用 nil，会使用系统分配的网卡，但不建议这么做。
	 * 但是，如果自己搜索网卡，有时候搜索不到，造成服务器无法启动，反而使用 nil 可以快速启动。
	 * 如果后续出现问题，可以尝试下面自己搜索网卡的方式。
	 */
	// host, _, err := getHostAndPortFromAddress(address)
	// if err != nil {
	// 	return nil, err
	// }
	// conn, err := net.ListenMulticastUDP("udp", findIfiForever(host), addr)
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// NewMulticastConnByIfiName 根据 <address> 创建并返回一个组播连接 *net.UDPConn
func NewMulticastConnByIfiName(address, name string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	ifi, err := findIfiByName(name)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenMulticastUDP("udp", ifi, addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Send 向组播地址 <address> 发送信息
func Send(address string, data []byte, retry ...gudp.Retry) error {
	conn, err := NewConn(address)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Send(data, retry...)
}

// SendString 向组播地址 <address> 发送字符串信息
func SendString(address, data string) error {
	return Send(address, []byte(data))
}

// SendJson 将 <data> JSON 序列化后发送到组播地址 <address>
func SendJson(address string, data interface{}) error {
	jsonByte, err := gjson.Encode(data)
	if err != nil {
		return err
	}
	return Send(address, jsonByte)
}

// getHostAndPortFromAddress 从地址 <address> 中获取 IP 和端口
func getHostAndPortFromAddress(address string) (string, int, error) {
	if !strings.Contains(address, ":") {
		return "", 0, gerror.New(fmt.Sprintf("错误的 IP 地址：%s", address))
	}
	strArr := strings.Split(address, ":")
	if len(strArr) != 2 {
		return "", 0, gerror.New(fmt.Sprintf("错误的 IP 地址：%s", address))
	}
	host := strArr[0]
	port, err := strconv.Atoi(strArr[1])
	if err != nil {
		return "", 0, err
	}
	return host, port, nil
}

// findIfiForever 每 5 秒查找一次对应组播地址 <host> 的网卡，直到找到位置。
func findIfiForever(host string) *net.Interface {
	ticker := time.NewTicker(5 * time.Second)
	timer := time.NewTimer(time.Millisecond)
	for {
		select {
		case <-ticker.C:
		case <-timer.C:
			if ifi, err := findIfi(host); err == nil {
				return ifi
			}
		}
	}
}

// findIfiByName 通过 <name> 获取网卡
func findIfiByName(name string) (*net.Interface, error) {
	var (
		ifis []net.Interface
		err  error
	)
	if ifis, err = net.Interfaces(); err != nil {
		return nil, err
	}
	for _, ifi := range ifis {
		if strings.Compare(ifi.Name, name) == 0 {
			return &ifi, nil
		}
	}
	return nil, gerror.New(fmt.Sprintf("not found ifi for %s", name))
}

// findIfi 查找对应组播地址 <host> 网卡。
func findIfi(host string) (*net.Interface, error) {
	var (
		ifis  []net.Interface
		addrs []net.Addr
		err   error
	)
	// 获取设备的所有网络接口
	if ifis, err = net.Interfaces(); err != nil {
		return nil, err
	}
	for _, ifi := range ifis {
		// 检查接口是否有组播地址
		if addrs, err = ifi.MulticastAddrs(); err != nil {
			// 不存在组播地址则直接跳过
			continue
		}
		for _, addr := range addrs {
			// 查看 <host> 和组播地址前 3 位是否一致
			// 比如 234.0.0.0 和 234.234.1.1 是匹配的
			// log.Printf("%s ==> %s", ifi.Name, addr.String())
			if strings.HasPrefix(addr.String(), strings.Split(host, ".")[0]) {
				return &ifi, nil
			}
		}
	}
	return nil, gerror.New(fmt.Sprintf("not found ifi to %s", host))
}
