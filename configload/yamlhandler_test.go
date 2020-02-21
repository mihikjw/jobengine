package configload

import (
	"testing"
)

//MockYAMLHandler mocks YAMLHandler
type MockYAMLHandler struct {
	UnmarshalResult    error
	UnmarshalOutput    map[interface{}]interface{}
	MarshalByteResult  []byte
	MarshalErrorResult error
}

//Unmarshal transforms a byte slice of YAML to an interface{}
func (h *MockYAMLHandler) Unmarshal(in []byte) (map[interface{}]interface{}, error) {
	return h.UnmarshalOutput, h.UnmarshalResult
}

//Marshal transforms an input interface{} to a byte slice of YAML
func (h *MockYAMLHandler) Marshal(in interface{}) ([]byte, error) {
	return h.MarshalByteResult, h.MarshalErrorResult
}

//TestMarshal1 tests loading yaml bytes into an object
func TestMarshal1(t *testing.T) {
	testData := map[string]string{
		"Hello": "World",
	}
	client := new(YAMLHandler)

	result, err := client.Marshal(testData)
	if err != nil {
		t.Errorf("TestMarshal1: Unexpected Error: %s", err.Error())
	}
	if result == nil {
		t.Error("TestMarshal1: Result Unexpectedly Nil")
	}
}

//TestUnmarshal1 tests unloading YAML bytes into an object
func TestUnmarshal1(t *testing.T) {
	testData := []byte(`
   version: 1
   port: 6010
   queues:
      test_queue_1:
         read:
         - service1
         - service2
         write:
         - service3
      test_queue_2:
         read:
         - service3
         write:
         - service1
         - service2`)
	client := new(YAMLHandler)

	result, err := client.Unmarshal(testData)
	if err != nil {
		t.Errorf("TestUnmarshal1: Unexpected Error: %s", err.Error())
	}
	if len(result) != 3 {
		t.Error("TestUnmarshal1: Result Too Small")
	}
}
