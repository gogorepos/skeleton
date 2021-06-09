package rx

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogorepos/skeleton/ipssc"
)

type ListParam struct {
	IDList    []string
	Arguments []string
}

func (p ListParam) Get() (*garray.Array, error) {
	array := garray.New()
	command := make(map[string]interface{})
	if len(p.IDList) == 0 {
		command["idList"] = []string{}
	} else {
		command["idList"] = p.IDList
	}
	if len(p.Arguments) > 0 {
		for _, a := range p.Arguments {
			command[a] = "y"
		}
	}
	r, err := ipssc.Send("hscmd-get-rx-config", command)
	if err != nil {
		return nil, err
	}
	err = r.GetJson("parameter").Struct(array)
	return array, err
}

func GetAll() (*garray.Array, error) {
	p := ListParam{}
	return p.Get()
}
