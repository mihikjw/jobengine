package queue

import (
	"encoding/json"
)

//JSONHandler is an interface for JSONHandle
type JSONHandler interface {
	Unmarshal(in []byte) (map[string]interface{}, error)
	Marshal(in interface{}) ([]byte, error)
}

//JSONHandle is an implementation of ConfigParser for the JSON file format
type JSONHandle struct{}

//Unmarshal transforms a byte slice of JSON to a map
func (h *JSONHandle) Unmarshal(in []byte) (map[string]interface{}, error) {
	rawData := make(map[string]interface{})
	err := json.Unmarshal(in, &rawData)
	return rawData, err
}

//Marshal transforms an input interface{} to a byte slice of JSON
func (h *JSONHandle) Marshal(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

//MockJSONHandle is a mock object implementing JSONHandler
type MockJSONHandle struct {
	UnmarshalResult map[string]interface{}
	UnmarshalError  error
	MarshalResult   []byte
	MarshalError    error
}

//Unmarshal returns UnmarshalResult & UnmarshalError
func (h *MockJSONHandle) Unmarshal(in []byte) (map[string]interface{}, error) {
	return h.UnmarshalResult, h.UnmarshalError
}

//Marshal returns MarshalResult & MarshalError
func (h *MockJSONHandle) Marshal(in interface{}) ([]byte, error) {
	return h.MarshalResult, h.MarshalError
}
