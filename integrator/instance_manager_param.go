package integrator

import (
	"io/ioutil"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type InstanceManagerParamState struct {
	m bridge.InstanceManager
}

func createInstanceManagerParamState(ctx *core.Context, params data.Map) (core.SharedState, error) {
	p, err := params.Get("file")
	if err != nil {
		return nil, err
	}
	path, err := data.AsString(p)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	managerConfig := string(b)
	s := &InstanceManagerParamState{}
	s.m = bridge.NewInstanceManager(managerConfig)

	return s, nil
}

func (s *InstanceManagerParamState) CreateNewState() func(*core.Context, data.Map) (core.SharedState, error) {
	return createInstanceManagerParamState
}

func (s *InstanceManagerParamState) TypeName() string {
	return "instance_manager_parameter"
}

func (s *InstanceManagerParamState) Terminate(ctx *core.Context) error {
	s.m.Delete()
	return nil
}
