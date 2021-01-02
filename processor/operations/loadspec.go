package operations

import (
	"PLViewer/processor/environment"
	"fmt"
)

type LoadSpec struct {
	*operation
	name     string
	zip      string
	typeName string
	uuid     string
}

func MakeLoadSpec() *LoadSpec {
	return &LoadSpec{
		operation: &operation{
			method: LoadSpectrumMethod,
		},
		name:     "",
		zip:      "",
		typeName: "",
		uuid:     "",
	}
}

func (operation *LoadSpec) GetMethod() Method {
	return operation.method
}

func (operation *LoadSpec) GetName() string {
	return operation.name
}

func (operation *LoadSpec) SetName(name string) {
	if name != "" {
		operation.name = name
	}
}

func (operation *LoadSpec) GetZip() string {
	return operation.zip
}

func (operation *LoadSpec) SetZip(zip string) {
	if zip != "" {
		operation.zip = zip
	}
}

func (operation *LoadSpec) GetUuid() string {
	return operation.uuid
}

func (operation *LoadSpec) SetUuid(uuid string) {
	if uuid != "" {
		operation.uuid = uuid
	}
}

func (operation *LoadSpec) GetType() string {
	return operation.typeName
}

func (operation *LoadSpec) SetType(typeName string) {
	if typeName != "" {
		operation.typeName = typeName
	}
}

func (operation *LoadSpec) RunEnvironment(env *environment.Environment) *environment.Environment {
	newEnv := env.Copy()
	newEnv.SetSpectrumFromZip(operation.name, operation.zip, operation.typeName, operation.uuid)

	return newEnv
}

func (operation *LoadSpec) Serialize() []string {
	return []string{fmt.Sprintf("-l:%s:%s:%s", operation.name, operation.uuid, operation.typeName), operation.env.GetZip(operation.zip).Uri}
}
