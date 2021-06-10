package iswitch

import (
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gogorepos/skeleton/proto/snmp"
)

func Run(ips ...string) (map[string][]IfUnit, error) {
	var (
		connectedIP  = make(map[string]bool)
		switchIPUnit = make(map[string][]IfUnit)
		c            = make(chan string, len(ips)+1)
		wg           sync.WaitGroup
	)
	// 如果不指定 IP 地址，则从 169.254.0.1 到 169.254.0.254 中能连接的服务
	if len(ips) == 0 {
		ips = Ping("169.254.0.1", "169.254.0.254")
	}
	go func() {
		for {
			select {
			case ip := <-c:
				if !connectedIP[ip] {
					wg.Add(1)
					connectedIP[ip] = true
					go func() {
						defer wg.Done()
						ifUnitSlice, err := Fetch(ip, c)
						if err != nil {
							return
						}
						switchIPUnit[ip] = ifUnitSlice
					}()
				}
			}
		}
	}()
	for _, ip := range ips {
		c <- ip
	}
	time.Sleep(time.Second)
	wg.Wait()
	return switchIPUnit, nil
}

func Fetch(ip string, c chan<- string) ([]IfUnit, error) {
	s, err := snmp.NewSNMP(ip)
	if err != nil {
		return nil, err
	}
	defer s.Close()
	// 获取接口数量
	ifNumber, err := GetIfNumber(s)
	if err != nil {
		return nil, err
	}
	// 创建接口设备列表
	ifUnitSlice := make([]IfUnit, ifNumber)
	// 接口描述哈希表，接口描述 => 接口列表的下标
	ifDescIndex := make(map[string]int)
	r, err := s.Walk(IfDescOid)
	if err != nil {
		return nil, err
	}
	// 遍历结果，保存接口描述，并在 ifDescIndex 记录每个接口描述对应的下标
	for i, result := range r {
		desc := result.String()
		ifUnitSlice[i].Desc = desc
		ifDescIndex[desc] = i
	}
	// 获取交换机连接其他交换机的端口数
	r, err = s.Walk(OccupiedPortOid)
	if err != nil {
		return nil, err
	}
	count := len(r)
	for i := 1; i <= count; i++ {
		iString := strconv.Itoa(i)
		// 获取本地端口
		num, err := s.GetNext(IndexLocalPortOid + iString)
		if err != nil {
			continue
		}
		// 获取本地端口 ID
		num, err = s.Get(IndexLocalIDOid + iString)
		if err != nil {
			log.Printf("id %v", err)
			continue
		}
		id := num.String()
		// 获取本地端口描述
		num, err = s.Get(IndexLocalDesOid + iString)
		if err != nil {
			continue
		}
		description := num.String()
		num, err = s.GetNext(IndexRemoteIPOid + iString)
		if err != nil {
			continue
		}
		if index, ok := ifDescIndex[description]; ok {
			remoteIP := num.String()
			ifUnitSlice[index].ID = id
			ifUnitSlice[index].IP = remoteIP
			if remoteIP != "" {
				c <- remoteIP
			}
		}
	}
	// 获取每个接口的状态
	r, err = s.Walk(IfStatusOid)
	if err == nil {
		for i, result := range r {
			ifUnitSlice[i].Status = result.Int()
		}
	}
	// 获取每个接口的带宽
	r, err = s.Walk(IfSpeedOid)
	if err == nil {
		for i, result := range r {
			ifUnitSlice[i].Speed = result.Int()
		}
	}
	// 获取每个接口的物理地址
	r, err = s.Walk(IfPMacOid)
	if err == nil {
		for i, result := range r {
			ifUnitSlice[i].Mac = result.MacString()
		}
	}
	// 获取端口和连接 Mac 地址映射表
	portMac, err := GetMacAddress(s)
	if err == nil {
		for i, _ := range ifUnitSlice {
			if m, ok := portMac[i+1]; ok && len(m) == 1 {
				ifUnitSlice[i].Mac = m[0]
			}
		}
	}
	return ifUnitSlice, nil
}

// Ping 获取 <start> 到 <end> IP 网段能够连接 snmp 服务器的 IP 地址
func Ping(start, end string) []string {
	startIP := net.ParseIP(start).To4()
	endIP := net.ParseIP(end).To4()
	// 如果其中有一个不能转换为 IPv4
	if startIP == nil || endIP == nil {
		return nil
	}
	var (
		ips  []string                               // 保存结果
		ip   = make(chan string)                    // IP 地址传递
		flag = make(chan bool)                      // 结束信号传递
		t    = time.NewTicker(time.Millisecond * 5) // 请求限速
	)
	go func() {
		for ; startIP[2] <= endIP[2]; startIP[2]++ {
			sIP := net.ParseIP(startIP.String()).To4()
			for ; sIP[3] <= endIP[3]; sIP[3]++ {
				<-t.C
				go func() {
					// 好像随便一个地址都可以连接，不会报错
					// 所以只能连接后请请求数据，如果正确回应则表示连接成功
					i := sIP.String()
					if s, err := snmp.NewSNMP(i); err == nil {
						if _, err = s.Get(OccupiedPortOid); err == nil {
							ip <- i
						}
						s.Close()
					}
				}()
				// net.IP 实际是 uint8，所以当 255 + 1 时，会得到 0
				// 只能手动判断如果是 255，那么跳到下一个网段
				if sIP[3] == 255 {
					break
				}
				// 当开始后结束 IP 地址相同的时候，表示扫描完毕
				// 给结束信道发送信号，返回结果
				if sIP.Equal(endIP) {
					flag <- true
				}
			}
		}
	}()
	for {
		select {
		case i := <-ip:
			ips = append(ips, i)
		case <-flag:
			return ips
		}
	}
}
