package processor

import (
	"PLViewer/processor/environment"
	"PLViewer/processor/operations"
)

func RunOperations(order []string, ops map[string]operations.Interface) *environment.Environment {
	env := environment.MakeEnvironment()
	for _, id := range order {
		op := ops[id]
		op.SetEnvironment(env)
		env = op.RunEnvironment(env)
	}

	return env
}

func SerializeOperations(order []string, ops map[string]operations.Interface) []string {
	RunOperations(order, ops)

	output := []string{}
	for _, id := range order {
		op := ops[id]
		output = append(output, op.Serialize()...)
	}

	return output
}
