package operations

import (
	"PLViewer/processor/environment"
	"fmt"
)

type CubicSpline struct {
	*operation
	name      string
	inName    string
	interpret int
}

func MakeCubicSpline() *CubicSpline {
	return &CubicSpline{
		operation: &operation{
			method: CubicSplineMethod,
		},
		name:      "",
		inName:    "",
		interpret: 1,
	}
}

func (operation *CubicSpline) GetMethod() Method {
	return operation.method
}

func (operation *CubicSpline) GetName() string {
	return operation.name
}

func (operation *CubicSpline) SetName(name string) {
	if name != "" {
		operation.name = name
	}
}

func (operation *CubicSpline) GetInName() string {
	return operation.inName
}

func (operation *CubicSpline) SetInName(inName string) {
	if inName != "" {
		operation.inName = inName
	}
}

func (operation *CubicSpline) GetInterpret() int {
	return operation.interpret
}

func (operation *CubicSpline) SetInterpret(interpret int) {
	if interpret > 0 {
		operation.interpret = interpret
	}
}

func (operation *CubicSpline) RunEnvironment(env *environment.Environment) *environment.Environment {
	newEnv := env.Copy()
	newEnv.SetSpectrumFromSpectrum(operation.name, operation.inName, "", "", [2]int{0, 0})

	return newEnv
}

func (operation *CubicSpline) Serialize() []string {
	return []string{fmt.Sprintf("-3:%s:%d", operation.inName, operation.interpret), operation.name}
}
