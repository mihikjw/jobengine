package configload

import "gopkg.in/yaml.v2"

//YAMLHandler is an implementation of ConfigParser for the YAML file format
type YAMLHandler struct{}

//Unmarshal transforms a byte slice of YAML to a map
func (h *YAMLHandler) Unmarshal(in []byte) (map[interface{}]interface{}, error) {
	out := make(map[interface{}]interface{})
	err := yaml.Unmarshal(in, &out)
	return out, err
}

//Marshal transforms an input interface{} to a byte slice of YAML
func (h *YAMLHandler) Marshal(in interface{}) ([]byte, error) {
	return yaml.Marshal(in)
}
