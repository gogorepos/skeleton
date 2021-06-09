package iswitch

import (
	"github.com/gogorepos/skeleton/proto/snmp"
	"github.com/gosnmp/gosnmp"
)

// GetIfNumber 获取交换机端口个数
func GetIfNumber(s *snmp.SNMP) (int, error) {
	num, err := s.Get(IfNumberOid)
	if err != nil {
		return 0, err
	}
	return num.Int(), err
}

func GetMacAddress(s *snmp.SNMP) (map[int][]string, error) {
	// 端口对应 Mac 地址
	portMac := make(map[int][]string)
	// Mac 地址切片
	var macSlice []string
	// 获取所有 Mac 地址
	if err := s.WalkFunc(IfMacOid, func(u gosnmp.SnmpPDU) error {
		macSlice = append(macSlice, snmp.NewVar(u).MacString())
		return nil
	}); err != nil {
		return nil, err
	}
	// 获取 Mac 地址对应端口
	r, err := s.BulkWalk(IfMacPortOid)
	if err != nil {
		return nil, err
	}
	for i, result := range r {
		index := result.Int()
		if _, ok := portMac[index]; ok {
			portMac[index] = append(portMac[index], macSlice[i])
		} else {
			portMac[index] = []string{macSlice[i]}
		}
	}
	return portMac, nil
}
