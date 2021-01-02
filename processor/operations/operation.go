package operations

import (
	"PLViewer/processor/environment"
)

type Method int

const (
	AddZipMethod Method = iota
	LoadSpectrumMethod
	MakeSliceMethod
	ExportToCsvMethod
)

type Interface interface {
	GetMethod() Method
	RunEnvironment(*environment.Environment) *environment.Environment
	SetEnvironment(*environment.Environment)
	GetEnvironment() *environment.Environment
	Serialize() []string
}

type operation struct {
	method Method
	env    *environment.Environment
}

func (op *operation) SetEnvironment(environment *environment.Environment) {
	op.env = environment
}

func (op *operation) GetEnvironment() *environment.Environment {
	return op.env
}
