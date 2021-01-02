package operations

import (
	"PLViewer/processor/environment"
	"fmt"
	"strings"
)

type Csv struct {
	*operation
	dataNames []string
	uri       string
}

func MakeCsv() *Csv {
	return &Csv{
		operation: &operation{
			method: ExportToCsvMethod,
		},
		dataNames: []string{},
		uri:       "",
	}
}

func (operation *Csv) GetMethod() Method {
	return operation.method
}

func (operation *Csv) GetDataNames() string {
	return strings.Join(operation.dataNames, ", ")
}

func (operation *Csv) SetDataNames(names string) {
	rawNames := strings.Split(names, ",")
	outputNames := []string{}
	for _, name := range rawNames {
		outputNames = append(outputNames, strings.TrimSpace(name))
	}

	operation.dataNames = outputNames
}

func (operation *Csv) GetUri() string {
	return operation.uri
}

func (operation *Csv) SetUri(uri string) {
	if uri != "" {
		if !strings.HasSuffix(uri, ".csv") {
			uri = uri + ".csv"
		}
		operation.uri = uri
	}
}

func (operation *Csv) RunEnvironment(env *environment.Environment) *environment.Environment {
	return env.Copy()
}

func (operation *Csv) Serialize() []string {
	return []string{fmt.Sprintf("-c:%s", strings.Join(operation.dataNames, ":")), operation.uri}
}
