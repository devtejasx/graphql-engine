package sources

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

func TestPlainTextConnectionString(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{
			"postgres url in plain text",
			`configuration: {connection_info: {database_url: "postgres://user:secret@db:5432/app"}}`,
			"database_url",
		},
		{
			"postgres url from env",
			`configuration: {connection_info: {database_url: {from_env: PG_DATABASE_URL}}}`,
			"",
		},
		{
			"postgres connection parameters",
			`configuration: {connection_info: {database_url: {connection_parameters: {host: db}}}}`,
			"",
		},
		{
			"mssql connection string in plain text",
			`configuration: {connection_info: {connection_string: "Server=db;Password=secret"}}`,
			"connection_string",
		},
		{
			"mssql connection string from env",
			`configuration: {connection_info: {connection_string: {from_env: MSSQL_URL}}}`,
			"",
		},
		{
			"no connection info",
			`configuration: {template: null}`,
			"",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var source struct {
				Configuration yaml.Node `yaml:"configuration"`
			}
			require.NoError(t, yaml.Unmarshal([]byte(tc.source), &source))
			assert.Equal(t, tc.want, plainTextConnectionString(source.Configuration))
		})
	}
}
