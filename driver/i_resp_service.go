package driver

import (
	"context"
	"fmt"
)

type RespServiceName string

type IRespService interface {
	// Start service
	Start(ctx context.Context) (err error)
	// InitRespConn init resp connect session by select db index,
	// return IRespConn interface
	InitRespConn(ctx context.Context, dbIdx int) IRespConn
	// Close resp service
	Close() (err error)
	// Name
	Name() RespServiceName
	// SetStorager
	SetStorager(store IStorager)
}

var respCmdSrvs = map[RespServiceName]IRespService{}

func RegisterRespCmdSrv(s IRespService) error {
	name := s.Name()
	if _, ok := respCmdSrvs[name]; ok {
		return fmt.Errorf("RespCmdSrv %s is registered", s)
	}

	respCmdSrvs[name] = s
	return nil
}

func ListRespCmdSrvs() []string {
	s := []string{}
	for k := range respCmdSrvs {
		s = append(s, string(k))
	}

	return s
}

func GetRespCmdSrv(name RespServiceName) (IRespService, error) {
	s, ok := respCmdSrvs[name]
	if !ok {
		return nil, fmt.Errorf("kv RespCmdSrv %s is not registered", name)
	}

	return s, nil
}
