package snmp

import (
	"strconv"
	"strings"

	"github.com/gosnmp/gosnmp"
)

type Var struct {
	gosnmp.SnmpPDU
}

func New(s gosnmp.SnmpPDU) *Var {
	return &Var{s}
}

func (v *Var) String() string {
	if n, ok := v.Value.([]byte); ok {
		return string(n)
	}
	if n, ok := v.Value.(int); ok {
		return strconv.Itoa(n)
	}
	return ""
}

func (v *Var) Int() int {
	return int(v.Int64())
}

func (v *Var) Uint() uint {
	return uint(v.Uint64())
}

// MacUpperString 获取大写十六进制的 Mac 地址
func (v *Var) MacUpperString() string {
	var b []string
	if s, ok := v.Value.([]uint8); ok {
		// 不是规范的 mac 地址，返回空字符串
		if len(s) != 6 {
			return ""
		}
		for _, u := range s {
			// 将十进制数字转换为大写十六进制的字符串
			n := strings.ToUpper(strconv.FormatUint(uint64(u), 16))
			// 如果只有一位，补零。
			if len(n) < 2 {
				n = "0" + n
			}
			b = append(b, n)
		}
	}
	// 各个字符用 ":" 分隔
	return strings.Join(b, ":")
}

// MacLowerString 获取小写十六进制的 Mac 地址
func (v *Var) MacLowerString() string {
	var b []string
	if s, ok := v.Value.([]uint8); ok {
		// 不是规范的 mac 地址，返回空字符串
		if len(s) != 6 {
			return ""
		}
		for _, u := range s {
			// 将十进制数字转换为小写十六进制的字符串
			n := strings.ToLower(strconv.FormatUint(uint64(u), 16))
			// 如果只有一位，补零。
			if len(n) < 2 {
				n = "0" + n
			}
			b = append(b, n)
		}
	}
	// 各个字符用 ":" 分隔
	return strings.Join(b, ":")
}

func (v *Var) MacString() string {
	return v.MacUpperString()
}

func (v *Var) Int64() int64 {
	b := gosnmp.ToBigInt(v.Value)
	if b.IsUint64() {
		return int64(b.Uint64())
	}
	return b.Int64()
}

func (v *Var) Uint64() uint64 {
	b := gosnmp.ToBigInt(v.Value)
	if b.IsInt64() {
		return uint64(b.Int64())
	}
	return b.Uint64()
}
