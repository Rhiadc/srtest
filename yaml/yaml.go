package yaml

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/rhiadc/srtest/resources"
	"gopkg.in/yaml.v3"
)

func GenerateK8SYAMLAndValidate(jsonData []byte) (string, error) {
	var crd resources.CustomResourceDefinition

	if err := json.Unmarshal(jsonData, &crd); err != nil {
		return "", err
	}

	yamlData, err := yaml.Marshal(&crd)
	if err != nil {
		return "", err
	}

	//_, err = kubeval.Validate(yamlData)
	//if err != nil {
	//		return nil
	//}

	yamlFileName := fmt.Sprintf("crd-%s.yaml", crd.Metadata.Name)
	filename := filepath.Join("pp", yamlFileName)
	err = ioutil.WriteFile(filename, yamlData, 0644)
	if err != nil {
		return "", err
	}
	return yamlFileName, nil
}
