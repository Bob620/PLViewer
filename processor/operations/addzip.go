package operations

import (
	"PLViewer/processor/environment"
)

type AddZip struct {
	*operation
	name string
	uri  string
}

func MakeAddZip() *AddZip {
	return &AddZip{
		operation: &operation{
			method: AddZipMethod,
		},
		name: "",
		uri:  "",
	}
}

func (operation *AddZip) GetMethod() Method {
	return operation.method
}

func (operation *AddZip) GetName() string {
	return operation.name
}

func (operation *AddZip) SetName(name string) {
	if name != "" {
		operation.name = name
	}
}

func (operation *AddZip) GetUri() string {
	return operation.uri
}

func (operation *AddZip) SetUri(uri string) {
	if uri != "" {
		operation.uri = uri
	}
}

func (operation *AddZip) RunEnvironment(env *environment.Environment) *environment.Environment {
	newEnv := env.Copy()
	newEnv.SetZip(operation.name, operation.uri)

	return newEnv
}

func (operation *AddZip) Serialize() []string {
	return []string{""}
}
