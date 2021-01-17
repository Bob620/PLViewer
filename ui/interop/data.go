package interop

import "sync"

type InteropData struct {
	strings           map[string]string
	stringsMutex      sync.RWMutex
	stringArrays      map[string]*StringArray
	stringArraysMutex sync.RWMutex
}

func MakeInteropData() *InteropData {
	interop := &InteropData{
		strings:           map[string]string{},
		stringsMutex:      sync.RWMutex{},
		stringArrays:      map[string]*StringArray{},
		stringArraysMutex: sync.RWMutex{},
	}

	return interop
}

func (interopData *InteropData) SetString(key, data string) {
	interopData.stringsMutex.Lock()
	defer interopData.stringsMutex.Unlock()
	interopData.strings[key] = data
}

func (interopData *InteropData) GetString(key string) string {
	interopData.stringsMutex.RLock()
	defer interopData.stringsMutex.RUnlock()
	data := interopData.strings[key]
	return data
}

func (interopData *InteropData) SetStringArray(key string, data []string) {
	interopData.stringArraysMutex.Lock()
	defer interopData.stringArraysMutex.Unlock()
	interopData.stringArrays[key] = &StringArray{
		strings: data,
		mutex:   sync.RWMutex{},
	}
}

func (interopData *InteropData) GetStringArray(key string) *StringArray {
	interopData.stringArraysMutex.RLock()
	data := interopData.stringArrays[key]
	interopData.stringArraysMutex.RUnlock()
	if data == nil {
		interopData.SetStringArray(key, []string{})
		return interopData.GetStringArray(key)
	}
	return data
}
