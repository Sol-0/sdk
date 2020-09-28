// Code generated by "go-syncmap -output per_connection_file_map.gen.go -type perConnectionFileMapMap<string,*perConnectionFileMap>"; DO NOT EDIT.

package recvfd

import "sync"

func _() {
	// An "cannot convert perConnectionFileMapMap literal (type perConnectionFileMapMap) to type sync.Map" compiler error signifies that the base type have changed.
	// Re-run the go-syncmap command to generate them again.
	_ = (sync.Map)(perConnectionFileMapMap{})
}
func (m *perConnectionFileMapMap) Store(key string, value *perConnectionFileMap) {
	(*sync.Map)(m).Store(key, value)
}

func (m *perConnectionFileMapMap) LoadOrStore(key string, value *perConnectionFileMap) (*perConnectionFileMap, bool) {
	actual, loaded := (*sync.Map)(m).LoadOrStore(key, value)
	if actual == nil {
		return nil, loaded
	}
	return actual.(*perConnectionFileMap), loaded
}

func (m *perConnectionFileMapMap) Load(key string) (*perConnectionFileMap, bool) {
	value, ok := (*sync.Map)(m).Load(key)
	if value == nil {
		return nil, ok
	}
	return value.(*perConnectionFileMap), ok
}

func (m *perConnectionFileMapMap) Delete(key string) {
	(*sync.Map)(m).Delete(key)
}

func (m *perConnectionFileMapMap) Range(f func(key string, value *perConnectionFileMap) bool) {
	(*sync.Map)(m).Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*perConnectionFileMap))
	})
}