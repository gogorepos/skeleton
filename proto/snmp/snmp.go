package snmp

import (
	"errors"
	"time"

	"github.com/gosnmp/gosnmp"
)

type SNMP struct {
	*gosnmp.GoSNMP
}

var NoneErr = errors.New("not result")

func NewSNMP(ip string) (*SNMP, error) {
	snmp := &gosnmp.GoSNMP{
		Target:             ip,
		Port:               161,
		Transport:          "udp",
		Community:          "public",
		Version:            gosnmp.Version2c,
		Timeout:            time.Second,
		Retries:            0,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
	}
	err := snmp.Connect()
	return &SNMP{snmp}, err
}

func (s *SNMP) Close() error {
	return s.Conn.Close()
}

func (s *SNMP) Get(oid string) (*Var, error) {
	oids := []string{oid}
	p, err := s.GoSNMP.Get(oids)
	if err != nil {
		return &Var{}, err
	}
	return getFirstResultFromPacket(p)
}

func (s *SNMP) GetAll(oids []string) ([]*Var, error) {
	p, err := s.GoSNMP.Get(oids)
	if err != nil {
		return nil, err
	}
	return getAllResultFromPacket(p)
}

func (s *SNMP) GetNext(oid string) (*Var, error) {
	oids := []string{oid}
	p, err := s.GoSNMP.GetNext(oids)
	if err != nil {
		return &Var{}, err
	}
	return getFirstResultFromPacket(p)
}

func (s *SNMP) GetNextAll(oids []string) ([]*Var, error) {
	p, err := s.GoSNMP.GetNext(oids)
	if err != nil {
		return nil, err
	}
	return getAllResultFromPacket(p)
}

func (s *SNMP) Walk(oid string) ([]*Var, error) {
	p, err := s.GoSNMP.WalkAll(oid)
	if err != nil {
		return nil, err
	}
	var result []*Var
	for _, pdu := range p {
		result = append(result, NewVar(pdu))
	}
	return result, nil
}

func (s *SNMP) WalkFunc(oid string, fun gosnmp.WalkFunc) error {
	return s.GoSNMP.Walk(oid, fun)
}

func (s *SNMP) BulkWalk(oid string) ([]*Var, error) {
	p, err := s.GoSNMP.BulkWalkAll(oid)
	if err != nil {
		return nil, err
	}
	var result []*Var
	for _, pdu := range p {
		result = append(result, NewVar(pdu))
	}
	return result, nil
}

func (s *SNMP) BulkWalkFunc(oid string, fun gosnmp.WalkFunc) error {
	return s.GoSNMP.BulkWalk(oid, fun)
}

// getFirstResultFromPacket 根据 <packet> 获取第一条结果
func getFirstResultFromPacket(packet *gosnmp.SnmpPacket) (*Var, error) {
	var result *Var
	if len(packet.Variables) == 0 {
		return result, NoneErr
	}
	result = NewVar(packet.Variables[0])
	return result, nil
}

// getAllResultFromPacket 根据 <packet> 获取所有结果
func getAllResultFromPacket(packet *gosnmp.SnmpPacket) ([]*Var, error) {
	var result []*Var
	for _, pdu := range packet.Variables {
		result = append(result, NewVar(pdu))
	}
	return result, nil
}
