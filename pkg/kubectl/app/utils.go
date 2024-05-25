package kubectl

import (
	// "fmt"
	"gopkg.in/yaml.v3"
)

func parseApiObjKind(filecontent []byte) (string, error) {
	// fmt.Println("parseApiObj")
	var obj map[string]interface{}
	err := yaml.Unmarshal(filecontent, &obj)
	if err != nil {
		return "", err
	}
	return obj["kind"].(string), nil
}
