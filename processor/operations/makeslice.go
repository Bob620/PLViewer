package operations

import (
	"PLViewer/processor/environment"
	"fmt"
	"strconv"
)

type MakeSlice struct {
	*operation
	name   string
	inName string
	high   int
	low    int
}

func MakeMakeSlice() *MakeSlice {
	return &MakeSlice{
		operation: &operation{
			method: MakeSliceMethod,
		},
		name:   "",
		inName: "",
		high:   0,
		low:    0,
	}
}

func (operation *MakeSlice) GetMethod() Method {
	return operation.method
}

func (operation *MakeSlice) GetName() string {
	return operation.name
}

func (operation *MakeSlice) SetName(name string) {
	if name != "" {
		operation.name = name
	}
}

func (operation *MakeSlice) GetInName() string {
	return operation.inName
}

func (operation *MakeSlice) SetInName(inName string) {
	if inName != "" {
		operation.inName = inName
	}
}

func (operation *MakeSlice) GetHigh() int {
	return operation.high
}

func (operation *MakeSlice) SetHigh(high string) {
	if high == "" {
		operation.high = 0
	} else {
		highInt, err := strconv.Atoi(high)
		if err == nil && highInt > operation.low {
			operation.high = highInt
		}
	}
}

func (operation *MakeSlice) GetLow() int {
	return operation.low
}

func (operation *MakeSlice) SetLow(low string) {
	if low == "" {
		operation.low = 0
	} else {
		lowInt, err := strconv.Atoi(low)
		if err == nil && (lowInt < operation.high || operation.high == 0) {
			operation.low = lowInt
		}
	}
}

func (operation *MakeSlice) RunEnvironment(env *environment.Environment) *environment.Environment {
	newEnv := env.Copy()
	newEnv.SetSpectrumFromSpectrum(operation.name, operation.inName, "", "", [2]int{operation.low, operation.high})

	return newEnv
}

func (operation *MakeSlice) Serialize() []string {
	return []string{fmt.Sprintf("-s:%s:%d-%d", operation.inName, operation.low, operation.high), operation.name}
}
