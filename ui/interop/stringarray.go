package interop

import "sync"

type StringArray struct {
	strings []string
	mutex   sync.RWMutex
}

func (stringArray *StringArray) Add(data string) {
	stringArray.mutex.Lock()
	defer stringArray.mutex.Unlock()
	stringArray.strings = append(stringArray.strings, data)
}

func (stringArray *StringArray) Remove(data string) {
	stringArray.mutex.Lock()
	defer stringArray.mutex.Unlock()
	newElements := make([]string, len(stringArray.strings))

	for _, element := range stringArray.strings {
		if element != data {
			newElements = append(newElements, element)
		}
	}
	stringArray.strings = newElements
}

func (stringArray *StringArray) GetAll() []string {
	newElements := make([]string, len(stringArray.strings))
	stringArray.mutex.RLock()
	defer stringArray.mutex.RUnlock()

	elements := stringArray.strings
	copy(newElements, elements)
	return newElements
}

func (stringArray *StringArray) Has(data string) bool {
	stringArray.mutex.RLock()
	defer stringArray.mutex.RUnlock()
	for _, element := range stringArray.strings {
		if element == data {
			return true
		}
	}
	return false
}
