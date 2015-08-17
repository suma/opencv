package integrator

import (
	"io/ioutil"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// InstanceManagerParamState is a shared state used by Instance Manager UDF/UDSF.
type InstanceManagerParamState struct {
	m bridge.InstanceManager
}

func createInstanceManagerParamState(ctx *core.Context, params data.Map) (
	core.SharedState, error) {
	p, err := params.Get(filePath)
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

// CreateNewState creates a state of Instance Manager parameters. The parameter
// is collected on JSON, see `scouter::InstanceManager::Config`.
//
// Usage of WITH parameter:
//  "file": instance manager parameters file path
func (s *InstanceManagerParamState) CreateNewState() func(
	*core.Context, data.Map) (core.SharedState, error) {
	return createInstanceManagerParamState
}

func (s *InstanceManagerParamState) TypeName() string {
	return "scouter_instance_manager_param"
}

func (s *InstanceManagerParamState) Terminate(ctx *core.Context) error {
	s.m.Delete()
	return nil
}
