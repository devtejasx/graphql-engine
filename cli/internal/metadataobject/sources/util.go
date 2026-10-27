package sources

import (
	"fmt"
	"strings"

	"github.com/hasura/graphql-engine/cli/v2/internal/errors"
	"go.yaml.in/yaml/v3"
)

// mappingValue returns the value of key in a YAML mapping node, or nil.
func mappingValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}

	return nil
}

// plainTextConnectionString returns the key of a source's connection string
// (database_url for Postgres like sources, connection_string for MS SQL
// Server) when it is set to a literal value instead of {from_env: ...}, so it
// would be written to the project files. It returns "" otherwise.
func plainTextConnectionString(configuration yaml.Node) string {
	connectionInfo := mappingValue(&configuration, "connection_info")
	for _, key := range []string{"database_url", "connection_string"} {
		value := mappingValue(connectionInfo, key)
		if value != nil && value.Kind == yaml.ScalarNode && value.Value != "" {
			return key
		}
	}

	return ""
}

func getTableYamlFileName(node yaml.Node) (string, error) {
	var op errors.Op = "sources.getTableYamlFileName"

	if node.Kind == yaml.SequenceNode {
		filename, err := fileNameWhenTableIsArray(node)
		if err != nil {
			return "", errors.E(op, err)
		}

		return filename, nil
	}

	filename, err := fileNameWhenTableIsObject(node)
	if err != nil {
		return "", errors.E(op, err)
	}

	return filename, nil
}

func fileNameWhenTableIsObject(yamlNode yaml.Node) (string, error) {
	var (
		op   errors.Op = "sources.fileNameWhenTableIsObject"
		node struct {
			Name string `yaml:"name"`
			// Depending on the datasource the namespacing parameter can change.
			// for example in pg and mssql the namespacing is done with the concept of schemas
			// and in cases like bigquery it'll be called "dataset"
			Schema  string `yaml:"schema"`
			Dataset string `yaml:"dataset"`
		}
	)

	yamlString, _ := yaml.Marshal(yamlNode)

	err := yaml.Unmarshal(yamlString, &node)
	if err != nil {
		return "", errors.E(op, err)
	}

	if node.Schema == "" {
		return fmt.Sprintf("%s_%s.yaml", node.Dataset, node.Name), nil
	}

	return fmt.Sprintf("%s_%s.yaml", node.Schema, node.Name), nil
}

func fileNameWhenTableIsArray(yamlNode yaml.Node) (string, error) {
	var op errors.Op = "sources.fileNameWhenTableIsArray"

	var strArr []string // var that stores individual elememts as strings

	yamlString, _ := yaml.Marshal(yamlNode)

	err := yaml.Unmarshal(yamlString, &strArr)
	if err != nil {
		return "", errors.E(op, err)
	}

	fileName := strings.Join(
		strArr,
		"_",
	) + ".yaml" // Concat all elements in the array delimited with "_", and append .yaml at the end

	return fileName, nil
}
